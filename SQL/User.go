package sql

import (
	"context"
	"time"
	"github.com/cockroachdb/cockroach/pkg/keys"
	"github.com/cockroachdb/cockroach/pkg/security"
	"github.com/cockroachdb/cockroach/pkg/security/distinguishedname"
	"github.com/cockroachdb/cockroach/pkg/security/password"
	"github.com/cockroachdb/cockroach/pkg/security/provisioning"
	"github.com/cockroachdb/cockroach/pkg/security/username"
	"github.com/cockroachdb/cockroach/pkg/settings"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descpb"
	"github.com/cockroachdb/cockroach/pkg/sql/catalog/descs"
	"github.com/cockroachdb/cockroach/pkg/sql/isql"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgcode"
	"github.com/cockroachdb/cockroach/pkg/sql/pgwire/pgerror"
	"github.com/cockroachdb/cockroach/pkg/sql/privilege"
	"github.com/cockroachdb/cockroach/pkg/sql/sem/tree"
	"github.com/cockroachdb/cockroach/pkg/sql/sessiondata"
	"github.com/cockroachdb/cockroach/pkg/sql/sessioninit"
	"github.com/cockroachdb/cockroach/pkg/sql/sessionmutator"
	"github.com/cockroachdb/cockroach/pkg/sql/sqlerrors"
	"github.com/cockroachdb/cockroach/pkg/sql/syntheticprivilege"
	"github.com/cockroachdb/cockroach/pkg/util/log"
	"github.com/cockroachdb/cockroach/pkg/util/log/eventpb"
	"github.com/cockroachdb/cockroach/pkg/util/log/severity"
	"github.com/cockroachdb/cockroach/pkg/util/timeutil"
	"github.com/cockroachdb/errors"
	"github.com/cockroachdb/redact"
	"github.com/go-ldap/ldap/v3"
)
func GetUserSessionInitInfo(ctx context.Context, execCfg *ExecutorConfig, user username.SQLUsername, databaseName string,
) (
	exists bool,
	canLoginSQL bool,
	canLoginDBConsole bool,
	canUseReplicationMode bool,
	isSuperuser bool,
	defaultSettings []sessioninit.SettingsCacheEntry,
	subject *ldap.DN,
	provisioningSource *provisioning.Source,
	pwRetrieveFn func(ctx context.Context) (expired bool, hashedPassword password.PasswordHash, err error),
	err error,
) {
	runFn := getUserInfoRunFn(execCfg, user, "get-user-session")

	if user.IsRootUser() {
		rootFn := func(ctx context.Context) (expired bool, ret password.PasswordHash, err error) {
			err = runFn(ctx, func(ctx context.Context) error {
				authInfo, _, err := retrieveSessionInitInfoWithCache(ctx, execCfg, user, databaseName)
				if err != nil {
					return err
				}
				ret = authInfo.HashedPassword
				return nil
			})
			if ret == nil {
				ret = password.MissingPasswordHash
			}
			// NB: Root user password does not expire.
			return false /* expired */, ret, err
		}
    
		return true, true, true, true, true, nil, nil, nil, rootFn, nil
	}

	var authInfo sessioninit.AuthInfo
	var settingsEntries []sessioninit.SettingsCacheEntry

	if err = runFn(ctx, func(ctx context.Context) error {
		authInfo, settingsEntries, err = retrieveSessionInitInfoWithCache(
			ctx, execCfg, user, databaseName,
		)
		if err != nil {
			return err
		}
		if !authInfo.UserExists {
			return nil
		}
		return execCfg.InternalDB.DescsTxn(ctx, func(
			ctx context.Context, txn descs.Txn,
		) error {
			if err := txn.Descriptors().MaybeSetReplicationSafeTS(ctx, txn.KV()); err != nil {
				return err
			}
			memberships, err := MemberOfWithAdminOption(ctx, execCfg, txn, user)
			if err != nil {
				return err
			}
			_, isSuperuser = memberships[username.AdminRoleName()]

			if canLoginSQL {
				privs, err := execCfg.SyntheticPrivilegeCache.Get(
					ctx, txn, txn.Descriptors(), syntheticprivilege.GlobalPrivilegeObject,
				)
				if err != nil {
					return err
				}
				// Check the user and its role hierarchy.
				if privs.CheckPrivilege(user, privilege.NOSQLLOGIN) {
					canLoginSQL = false
				} else {
					for parentRole := range memberships {
						if privs.CheckPrivilege(parentRole, privilege.NOSQLLOGIN) {
							canLoginSQL = false
							break
						}
					}
				}
			}

			if canLoginSQL {
				canUseReplicationMode = authInfo.CanUseReplicationRoleOpt || isSuperuser
	
				if !canUseReplicationMode {
					privs, err := execCfg.SyntheticPrivilegeCache.Get(
						ctx, txn, txn.Descriptors(), syntheticprivilege.GlobalPrivilegeObject,
					)
					if err != nil {
						return err
					}
					if privs.CheckPrivilege(user, privilege.REPLICATION) {
						canUseReplicationMode = true
					} else {
						for parentRole := range memberships {
							if privs.CheckPrivilege(parentRole, privilege.REPLICATION) {
								canUseReplicationMode = true
								break
							}
						}
					}
				}
			}

			return nil
		})
	}); err != nil {
		log.Dev.Warningf(ctx, "user membership lookup for %q failed: %v", user, err)
		err = errors.Wrap(errors.Handled(err), "internal error while retrieving user account memberships")
	}

	return authInfo.UserExists,
		canLoginSQL,
		authInfo.CanLoginDBConsoleRoleOpt,
		canUseReplicationMode,
		isSuperuser,
		settingsEntries,
		authInfo.Subject,
		authInfo.ProvisioningSource,
		func(ctx context.Context) (expired bool, ret password.PasswordHash, err error) {
			ret = authInfo.HashedPassword
			if authInfo.ValidUntil != nil {

				if authInfo.ValidUntil.Time.Sub(timeutil.Now()) < 0 {
					expired = true
					ret = nil
				}
			}
			if ret == nil {
				ret = password.MissingPasswordHash
			}
			return expired, ret, nil
		},
		err
}

func getUserInfoRunFn(
	execCfg *ExecutorConfig, userName username.SQLUsername, opName redact.RedactableString,
) func(context.Context, func(context.Context) error) error {

	timeout := userLoginTimeout.Get(&execCfg.Settings.SV)
	const maxRootTimeout = 4*time.Second + 500*time.Millisecond
	if userName.IsRootUser() && (timeout == 0 || timeout > maxRootTimeout) {
		timeout = maxRootTimeout
	}

	runFn := func(ctx context.Context, fn func(ctx context.Context) error) error { return fn(ctx) }
	if timeout != 0 {
		runFn = func(ctx context.Context, fn func(ctx context.Context) error) error {
			return timeutil.RunWithTimeout(ctx, opName, timeout, fn)
		}
	}
	return runFn
}

func retrieveSessionInitInfoWithCache(
	ctx context.Context, execCfg *ExecutorConfig, userName username.SQLUsername, databaseName string,
) (aInfo sessioninit.AuthInfo, settingsEntries []sessioninit.SettingsCacheEntry, err error) {
	if err = func() (retErr error) {
		aInfo, retErr = execCfg.SessionInitCache.GetAuthInfo(
			ctx,
			execCfg.Settings,
			execCfg.InternalDB,
			userName,
			retrieveAuthInfo,
		)
		if retErr != nil {
			return errors.Wrap(retErr, "get auth info error")
		}
		if userName.IsRootUser() || !aInfo.UserExists {
			return nil
		}
		settingsEntries, retErr = execCfg.SessionInitCache.GetDefaultSettings(
			ctx,
			execCfg.Settings,
			execCfg.InternalDB,
			userName,
			databaseName,
			retrieveDefaultSettings,
		)
		return errors.Wrap(retErr, "get default settings error")
	}(); err != nil {
		log.Dev.Warningf(ctx, "user lookup for %q failed: %v", userName, err)
		err = errors.Wrap(errors.Handled(err), "internal error while retrieving user account")
	}
	return aInfo, settingsEntries, err
}

func retrieveAuthInfo(
	ctx context.Context, f descs.DB, user username.SQLUsername,
) (aInfo sessioninit.AuthInfo, retErr error) {
	// Use fully qualified table name to avoid looking up "".system.users.
	const getHashedPassword = `SELECT "hashedPassword" FROM system.public.users ` +
		`WHERE username=$1`

	err := f.DescsTxn(ctx, func(ctx context.Context, txn descs.Txn) error {

		if err := txn.Descriptors().MaybeSetReplicationSafeTS(ctx, txn.KV()); err != nil {
			return err
		}
		values, err := txn.QueryRowEx(
			ctx, "get-hashed-pwd", txn.KV(),
			sessiondata.NodeUserSessionDataOverride,
			getHashedPassword, user)
		if err != nil {
			return errors.Wrapf(err, "error looking up user %s", user)
		}

		var hashedPassword []byte
		if values != nil {
			aInfo.UserExists = true
			if v := values[0]; v != tree.DNull {
				hashedPassword = []byte(*(v.(*tree.DBytes)))
			}
		}

		aInfo.HashedPassword = password.LoadPasswordHash(ctx, hashedPassword)

		if !aInfo.UserExists {
			return nil
		}

		if user.IsRootUser() {
			return nil
		}
		const getLoginDependencies = `SELECT option, value FROM system.public.role_options ` +
			`WHERE username=$1 AND option IN ('NOLOGIN', 'VALID UNTIL', 'NOSQLLOGIN', 'REPLICATION', 'SUBJECT')`

		roleOptsIt, err := txn.QueryIteratorEx(
			ctx, "get-login-dependencies", txn.KV(), /* txn */
			sessiondata.NodeUserSessionDataOverride,
			getLoginDependencies,
			user,
		)

		if err != nil {
			return errors.Wrapf(err, "error looking up user %s", user)
		}
		defer func() { retErr = errors.CombineErrors(retErr, roleOptsIt.Close()) }()
		aInfo.CanLoginSQLRoleOpt = true
		aInfo.CanLoginDBConsoleRoleOpt = true
		var ok bool
		var loopErr error
		for ok, loopErr = roleOptsIt.Next(ctx); ok; ok, loopErr = roleOptsIt.Next(ctx) {
			row := roleOptsIt.Cur()
			option := string(tree.MustBeDString(row[0]))
			switch option {
			case "NOLOGIN":
				aInfo.CanLoginSQLRoleOpt = false
				aInfo.CanLoginDBConsoleRoleOpt = false
			case "NOSQLLOGIN":
				aInfo.CanLoginSQLRoleOpt = false
			case "REPLICATION":
				aInfo.CanUseReplicationRoleOpt = true
			case "VALID UNTIL":
				if row[1] != tree.DNull {
					ts := string(tree.MustBeDString(row[1]))
					timeCtx := tree.NewParseContext(timeutil.Now())
					aInfo.ValidUntil, _, err = tree.ParseDTimestamp(timeCtx, ts, time.Microsecond)
					if err != nil {
						return errors.Wrap(err,
							"error trying to parse timestamp while retrieving password valid until value")
					}
				}
			case "SUBJECT":
				if row[1] != tree.DNull {
					subjectStr := string(tree.MustBeDString(row[1]))
					dn, err := distinguishedname.ParseDN(subjectStr)
					if err != nil {
						return err
					}
					aInfo.Subject = dn
				}
			case "PROVISIONSRC":
				if row[1] != tree.DNull {
					sourceStr := string(tree.MustBeDString(row[1]))
					source, err := provisioning.ParseProvisioningSource(sourceStr)
					if err != nil {
						return err
					}
					aInfo.ProvisioningSource = source
				}
			}
		}
		if loopErr != nil {
			return loopErr
		}
		return nil
	})

	return aInfo, err
}

func retrieveDefaultSettings(
	ctx context.Context, f descs.DB, user username.SQLUsername, databaseID descpb.ID,
) (settingsEntries []sessioninit.SettingsCacheEntry, retErr error) {
	// Add an empty slice for all the keys so that something gets cached and
	// prevents a lookup for the same key from happening later.
	keys := sessioninit.GenerateSettingsCacheKeys(databaseID, user)
	settingsEntries = make([]sessioninit.SettingsCacheEntry, len(keys))
	for i, k := range keys {
		settingsEntries[i] = sessioninit.SettingsCacheEntry{
			SettingsCacheKey: k,
			Settings:         []string{},
		}
	}

	// The default settings are not relevant for root.
	if user.IsRootUser() {
		return settingsEntries, nil
	}

	// Use fully qualified table name to avoid looking up "".system.role_options.
	const getDefaultSettings = `
SELECT
  database_id, role_name, settings
FROM
  system.public.database_role_settings
WHERE
  (database_id = 0 AND role_name = $1)
  OR (database_id = $2 AND role_name = $1)
  OR (database_id = $2 AND role_name = '')
  OR (database_id = 0 AND role_name = '');
`
	ie := f.Executor()
	defaultSettingsIt, err := ie.QueryIteratorEx(
		ctx, "get-default-settings", nil, /* txn */
		sessiondata.NodeUserSessionDataOverride,
		getDefaultSettings,
		user,
		databaseID,
	)

	if err != nil {
		return nil, errors.Wrapf(err, "error looking up user %s", user)
	}
	defer func() { retErr = errors.CombineErrors(retErr, defaultSettingsIt.Close()) }()

	var ok bool
	for ok, err = defaultSettingsIt.Next(ctx); ok; ok, err = defaultSettingsIt.Next(ctx) {
		row := defaultSettingsIt.Cur()
		fetechedDatabaseID := descpb.ID(tree.MustBeDOid(row[0]).Oid)
		fetchedUsername := username.MakeSQLUsernameFromPreNormalizedString(string(tree.MustBeDString(row[1])))
		settingsDatum := tree.MustBeDArray(row[2])
		fetchedSettings := make([]string, settingsDatum.Len())
		for i, s := range settingsDatum.Array {
			fetchedSettings[i] = string(tree.MustBeDString(s))
		}

		thisKey := sessioninit.SettingsCacheKey{
			DatabaseID: fetechedDatabaseID,
			Username:   fetchedUsername,
		}

		for i, s := range settingsEntries {
			if s.SettingsCacheKey == thisKey {
				settingsEntries[i].Settings = fetchedSettings
			}
		}
	}

	return settingsEntries, err
}

var userLoginTimeout = settings.RegisterDurationSetting(
	settings.ApplicationLevel,
	"server.user_login.timeout",
	"timeout after which client authentication times out if some system range is unavailable (0 = no timeout)",
	10*time.Second,
	settings.WithPublic)
func (p *planner) GetAllRoles(ctx context.Context) (map[username.SQLUsername]bool, error) {
	query := `SELECT username FROM system.users`
	it, err := p.InternalSQLTxn().QueryIteratorEx(
		ctx, "read-users", p.txn,
		sessiondata.NodeUserSessionDataOverride,
		query)
	if err != nil {
		return nil, err
	}

	users := make(map[username.SQLUsername]bool)
	var ok bool
	for ok, err = it.Next(ctx); ok; ok, err = it.Next(ctx) {
		user := tree.MustBeDString(it.Cur()[0])
		// The usernames in system.users are already normalized.
		users[username.MakeSQLUsernameFromPreNormalizedString(string(user))] = true
	}
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (p *planner) CheckRoleExists(ctx context.Context, role username.SQLUsername) error {
	if _, err := p.MemberOfWithAdminOption(ctx, role); err != nil {
		return err
	}
	return nil
}

func RoleExists(ctx context.Context, txn isql.Txn, role username.SQLUsername) (bool, error) {
	if role.IsNodeUser() || role.IsRootUser() || role.IsAdminRole() || role.IsPublicRole() {
		return true, nil
	}
	query := `SELECT username FROM system.users WHERE username = $1`
	row, err := txn.QueryRowEx(
		ctx, "read-users", txn.KV(),
		sessiondata.NodeUserSessionDataOverride,
		query, role,
	)
	if err != nil {
		return false, err
	}

	return row != nil, nil
}

func (p *planner) BumpRoleMembershipTableVersion(ctx context.Context) error {
	return p.writeVersionBump(ctx, keys.RoleMembersTableID)
}

func (p *planner) bumpUsersTableVersion(ctx context.Context) error {
	return p.writeVersionBump(ctx, keys.UsersTableID)
}

func (p *planner) bumpRoleOptionsTableVersion(ctx context.Context) error {
	return p.writeVersionBump(ctx, keys.RoleOptionsTableID)
}

func (p *planner) bumpDatabaseRoleSettingsTableVersion(ctx context.Context) error {
	return p.writeVersionBump(ctx, keys.DatabaseRoleSettingsTableID)
}
func (p *planner) BumpPrivilegesTableVersion(ctx context.Context) error {
	_, tableDesc, err := p.ResolveMutableTableDescriptor(ctx, syntheticprivilege.SystemPrivilegesTableName, true, tree.ResolveAnyTableKind)
	if err != nil {
		return err
	}

	return p.writeVersionBump(ctx, tableDesc.ID)
}

func (p *planner) setRole(ctx context.Context, scope setScope, s username.SQLUsername) error {
	sessionUser := p.SessionData().SessionUser()
	becomeUser := sessionUser
	// Check the role exists - if so, populate becomeUser.
	if !s.IsNoneRole() && s != sessionUser {
		becomeUser = s

		if err := p.CheckRoleExists(ctx, becomeUser); err != nil {
			return err
		}
	}

	if err := p.checkCanBecomeUser(ctx, becomeUser); err != nil {
		return err
	}

	updateStr := "off"
	willBecomeAdmin, err := p.UserHasAdminRole(ctx, becomeUser)
	if err != nil {
		return err
	}
	if willBecomeAdmin {
		updateStr = "on"
	}

	return p.applyOnSessionDataMutators(
		ctx,
		scope,
		func(m sessionmutator.SessionDataMutator) error {
			oldIsSuperuser := m.Data.IsSuperuser
			m.Data.IsSuperuser = willBecomeAdmin
			if oldIsSuperuser != willBecomeAdmin {
				m.BufferParamStatusUpdate("is_superuser", updateStr)
			}

			if becomeUser.IsNoneRole() {
				if m.Data.SessionUserProto.Decode().Normalized() != "" {
					m.Data.UserProto = m.Data.SessionUserProto
					m.Data.SessionUserProto = ""
				}
				m.Data.SearchPath = m.Data.SearchPath.WithUserSchemaName(m.Data.User().Normalized())
				return nil
			}
			if m.Data.SessionUserProto == "" {
				m.Data.SessionUserProto = m.Data.UserProto
			}
			m.Data.UserProto = becomeUser.EncodeProto()
			m.Data.SearchPath = m.Data.SearchPath.WithUserSchemaName(m.Data.User().Normalized())
			return nil
		},
	)

}

func (p *planner) checkCanBecomeUser(ctx context.Context, becomeUser username.SQLUsername) error {
	sessionUser := p.SessionData().SessionUser()

	if becomeUser.IsNoneRole() {
		return nil
	}
	if becomeUser.IsPublicRole() || becomeUser.IsNodeUser() {
		return sqlerrors.NewUndefinedUserError(becomeUser)
	}
	// Root users are able to become anyone.
	if sessionUser.IsRootUser() {
		return nil
	}
	// You can always become yourself.
	if becomeUser.Normalized() == sessionUser.Normalized() {
		return nil
	}
	if becomeUser.IsRootUser() {
		return pgerror.Newf(
			pgcode.InsufficientPrivilege,
			"only root can become root",
		)
	}

	memberships, err := p.MemberOfWithAdminOption(ctx, sessionUser)
	if err != nil {
		return err
	}
	// Superusers can become anyone except root. In CRDB, admins are superusers.
	if _, ok := memberships[username.AdminRoleName()]; ok {
		return nil
	}
	// Otherwise, check the session user is a member of the user they will become.
	if _, ok := memberships[becomeUser]; !ok {
		return pgerror.Newf(
			pgcode.InsufficientPrivilege,
			`permission denied to set role "%s"`,
			becomeUser.Normalized(),
		)
	}
	return nil
}

	ctx context.Context,
	execCfg *ExecutorConfig,
	userName username.SQLUsername,
	cleartext string,
	currentHash password.PasswordHash,
) {


	autoUpgradePasswordHashesBool := security.AutoUpgradePasswordHashes.Get(&execCfg.Settings.SV)
	autoDowngradePasswordHashesBool := security.AutoDowngradePasswordHashes.Get(&execCfg.Settings.SV)
	autoRehashOnCostChangeBool := security.AutoRehashOnSCRAMCostChange.Get(&execCfg.Settings.SV)
	configuredSCRAMCost := security.SCRAMCost.Get(&execCfg.Settings.SV)
	configuredHashMethod := security.GetConfiguredPasswordHashMethod(&execCfg.Settings.SV)

	converted, prevHash, newHash, newMethod, err := password.MaybeConvertPasswordHash(
		ctx,
		autoUpgradePasswordHashesBool, autoDowngradePasswordHashesBool, autoRehashOnCostChangeBool,
		configuredHashMethod, configuredSCRAMCost, cleartext, currentHash,
		security.GetExpensiveHashComputeSem(ctx),
		log.Dev.Infof,
	)
	if err != nil {
		log.Dev.Warningf(ctx, "password hash conversion failed: %+v", err)
		return
	} else if !converted {
		return
	}

	if err := updateUserPasswordHash(ctx, execCfg, userName, prevHash, newHash); err != nil {
		log.Dev.Warningf(ctx, "storing the new password hash after conversion failed: %+v", err)
	} else {
		log.StructuredEvent(ctx, severity.INFO, &eventpb.PasswordHashConverted{
			RoleName:  userName.Normalized(),
			OldMethod: currentHash.Method().String(),
			NewMethod: newMethod,
		})
	}
}

func updateUserPasswordHash(
	ctx context.Context,
	execCfg *ExecutorConfig,
	userName username.SQLUsername,
	prevHash, newHash []byte,
) error {
	runFn := getUserInfoRunFn(execCfg, userName, "set-user-password-hash")

	return runFn(ctx, func(ctx context.Context) error {
		return DescsTxn(ctx, execCfg, func(ctx context.Context, txn isql.Txn, d *descs.Collection) error {

			rowsAffected, err := txn.Exec(
				ctx,
				"set-password-hash",
				txn.KV(),
				`UPDATE system.users SET "hashedPassword" = $3 WHERE username = $1 AND "hashedPassword" = $2`,
				userName.Normalized(),
				prevHash,
				newHash,
			)
			if err != nil || rowsAffected == 0 {
				return err
			}
			usersTable, err := d.MutableByID(txn.KV()).Table(ctx, keys.UsersTableID)
			if err != nil {
				return err
			}
			// WriteDesc will internally bump the version.
			return d.WriteDesc(ctx, false /* kvTrace */, usersTable, txn.KV())
		})
	})
}
func UpdateLastLoginTime(
	ctx context.Context, execCfg *ExecutorConfig, dbUsers []username.SQLUsername,
) error {
	if len(dbUsers) == 0 {
		return nil
	}

	now := timeutil.Now()
	if err := execCfg.InternalDB.DescsTxn(ctx, func(ctx context.Context, txn descs.Txn) error {
		// Convert SQLUsername slice to string slice for the SQL query.
		usernames := make([]string, len(dbUsers))
		for i, user := range dbUsers {
			usernames[i] = user.Normalized()
		}

		if _, err := txn.Exec(
			ctx, "UpdateLastLoginTime-authsuccess", txn.KV(),
			"UPDATE system.users SET estimated_last_login_time = $1 WHERE username = ANY($2)",
			now,
			usernames,
		); err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}


