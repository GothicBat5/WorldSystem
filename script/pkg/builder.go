package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"github.com/spf13/cobra"
)

const volumeFlag = "volume"
const dockerArgsFlag = "docker-args"

func makeBuilderCmd(runE func(cmd *cobra.Command, args []string) error) *cobra.Command {
	builderCmd := &cobra.Command{
		Use: "builder",
		Short: "Run the Bazel builder image",
		Long: "Run the Bazel builder image.",
		Example: `dev builder`,
		Args: cobra.MinimumNArgs(0),
		RunE: runE,
	}
	builderCmd.Flags().String(volumeFlag, "bzlhome", "the Docker volume to use as the container home directory")
	return builderCmd
}

func (d *dev) builder(cmd *cobra.Command, extraArgs []string) error {
	ctx := cmd.Context()
	volume := mustGetFlagString(cmd, volumeFlag)
	var tty bool
	if len(extraArgs) == 0 {
		tty = true
	}
	args, err := d.getDockerRunArgs(ctx, volume, tty, nil)
	args = append(args, extraArgs...)
	if err != nil {
		return err
	}
	if tty {
		logCommand("docker", args...)
	}
	return d.exec.CommandContextInheritingStdStreams(ctx, "docker", args...)
}

func (d *dev) dockerIsPodman(ctx context.Context) (bool, error) {
	output, err := d.exec.CommandContextSilent(ctx, "docker", "help")
	if err != nil {
		return false, err
	}
	if bytes.Contains(output, []byte("podman")) {
		return true, nil
	}
	return false, nil
}

func (d *dev) getDockerRunArgs(
	ctx context.Context, volume string, tty bool, extraArgs []string,
) (args []string, err error) {
	err = d.ensureBinaryInPath("docker")
	if err != nil {
		return
	}
	uid, gid, err := d.os.CurrentUserAndGroup()
	if err != nil {
		return
	}
	uidInt, err := strconv.Atoi(uid)
	if err == nil && uidInt < 500 {
		uid = "501"
	}
	gidInt, err := strconv.Atoi(gid)
	if err == nil && gidInt < 500 {
		gid = uid
	}

	workspace, err := d.getWorkspace(ctx)
	if err != nil {
		return
	}
	bazelImageVersionFileContent, err := d.os.ReadFile(filepath.Join(workspace, "build", ".bazelbuilderversion"))
	if err != nil {
		return
	}
	bazelImage := strings.TrimSpace(bazelImageVersionFileContent)

	_, err = d.exec.CommandContextSilent(ctx, "docker", "volume", "inspect", volume)
	if err != nil {
		log.Printf("Creating volume %s with Docker...", volume)
		out, err := d.exec.CommandContextSilent(ctx, "docker", "volume", "create", volume)
		if err != nil {
			printStdoutAndErr(string(out), err)
			return nil, err
		}
		_, err = d.exec.CommandContextSilent(
			ctx, "docker", "run", "--rm", "-v", volume+":/vol",
			"--entrypoint=/usr/bin/chown", bazelImage,
			"-R", fmt.Sprintf("%s:%s", uid, gid), "/vol")
		if err != nil {
			return nil, err
		}
	}

	args = append(args, "run", "--rm")
	isPodman, err := d.dockerIsPodman(ctx)
	if err != nil {
		return nil, err
	}
	if isPodman {
		args = append(args, "--passwd=false")
	}
	if tty {
		args = append(args, "-it")
	} else {
		args = append(args, "-i")
	}
	gitDir, err := d.exec.CommandContextSilent(ctx, "git", "rev-parse", "--git-dir")
	if err != nil {
		return nil, err
	}
	gitCommonDir, err := d.exec.CommandContextSilent(ctx, "git", "rev-parse", "--git-common-dir")
	if err != nil {
		return nil, err
	}

	if string(gitDir) != string(gitCommonDir) {
		mountPath := strings.TrimSpace(string(gitCommonDir))
		args = append(args, "-v", mountPath+":"+mountPath)
	}
	args = append(args, "-v", workspace+":/cockroach")
	args = append(args, "--workdir=/cockroach")
	bazelRc := "empty.bazelrc"
	if !buildutil.CrdbTestBuild && (runtime.GOOS == "darwin" && runtime.GOARCH == "arm64") {
		bazelRc = "darwinarm64cross.bazelrc"
	}
	args = append(args, "-v", filepath.Join(workspace, "build", "bazelutil", bazelRc)+":/cockroach/.bazelrc.user")

	artifacts := filepath.Join(workspace, "artifacts")
	err = d.os.MkdirAll(artifacts)
	if err != nil {
		return
	}

	_ = d.os.Chmod(artifacts, 0777)
	args = append(args, "-v", artifacts+":/artifacts")
	args = append(args, "-v", volume+":/home/roach:delegated")
	args = append(args, extraArgs...)
	args = append(args, "-u", fmt.Sprintf("%s:%s", uid, gid))
	args = append(args, bazelImage)
	return
}
