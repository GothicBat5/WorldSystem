package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)


func trim(s string) string {
	return strings.TrimSpace(s)
}

func capitalize(s string) string {

	if s == "" {
		return ""
	}

	runes := []rune(strings.ToLower(s))

	first :=
		[]rune(
			strings.ToUpper(
				string(runes[0]),
			),
		)

	runes[0] = first[0]

	return string(runes)
}

func input(prompt string,
	reader *bufio.Reader) string {

	fmt.Print(prompt)

	text, _ :=
		reader.ReadString('\n')

	return trim(text)
}

// ---------- Core ----------

func addPerson(
	reader *bufio.Reader,
	roles map[string][]string,
) {

	name := input(
		"Enter name: ",
		reader,
	)

	if name == "" {

		fmt.Println(
			"Name cannot be empty.\n",
		)

		return
	}

	role := capitalize(
		input(
			"Enter role: ",
			reader,
		),
	)

	if role == "" {

		fmt.Println(
			"Role cannot be empty.\n",
		)

		return
	}

	// Duplicate prevention

	for _, person :=
		range roles[role] {

		if strings.EqualFold(
			person,
			name,
		) {

			fmt.Println(
				"Person already exists in this role.\n",
			)

			return
		}
	}

	roles[role] =
		append(
			roles[role],
			name,
		)

	fmt.Printf(
		"Added %s as %s\n\n",
		name,
		role,
	)
}

func showGroups(
	roles map[string][]string,
) {

	fmt.Println(
		"\n------ People Groups ------",
	)

	if len(roles) == 0 {

		fmt.Println(
			"No group found.\n",
		)

		return
	}

	var sortedRoles []string

	for role := range roles {

		sortedRoles =
			append(
				sortedRoles,
				role,
			)
	}

	sort.Strings(
		sortedRoles,
	)

	for _, role :=
		range sortedRoles {

		people :=
			append(
				[]string(nil),
				roles[role]...,
			)

		sort.Strings(
			people,
		)

		fmt.Printf(
			"%ss: %s\n",
			role,
			strings.Join(
				people,
				", ",
			),
		)
	}

	fmt.Println()
}

func deletePerson(
	reader *bufio.Reader,
	roles map[string][]string,
) {

	name :=
		input(
			"Enter name to delete: ",
			reader,
		)

	for role, people :=
		range roles {

		for i, p :=
			range people {

			if strings.EqualFold(
				p,
				name,
			) {

				roles[role] =
					append(
						people[:i],
						people[i+1:]...,
					)

				fmt.Printf(
					"Removed %s from %s\n\n",
					p,
					role,
				)

				return
			}
		}
	}

	fmt.Println(
		"Person not found.\n",
	)
}

func searchPerson(
	reader *bufio.Reader,
	roles map[string][]string,
) {

	name :=
		input(
			"Search name: ",
			reader,
		)

	target :=
		strings.ToLower(
			name,
		)

	found := false

	for role, people :=
		range roles {

		for _, person :=
			range people {

			if strings.Contains(
				strings.ToLower(
					person,
				),
				target,
			) {

				fmt.Printf(
					"%s is a %s\n",
					person,
					role,
				)

				found = true
			}
		}
	}

	if !found {

		fmt.Println(
			"No person found.\n",
		)

	} else {

		fmt.Println()
	}
}

// ---------- File ----------

func saveToFile(
	roles map[string][]string,
) {

	file, err :=
		os.Create(
			"roles.txt",
		)

	if err != nil {

		fmt.Println(
			"Error saving file.\n",
		)

		return
	}

	defer file.Close()

	writer :=
		bufio.NewWriter(
			file,
		)

	for role, people :=
		range roles {

		for _, person :=
			range people {

			fmt.Fprintf(
				writer,
				"%s|%s\n",
				role,
				person,
			)
		}
	}

	writer.Flush()

	fmt.Println(
		"Data saved!\n",
	)
}

func loadFromFile(roles map[string][]string) {

	file, err := os.Open("roles.txt",)

	if err != nil {

		return
	}

	defer file.Close()

	scanner :=
		bufio.NewScanner(file, )

	for scanner.Scan() {

		line :=	scanner.Text()

		parts := strings.Split(line, "|",)

		if len(parts) != 2 {

			continue
		}

		role :=	parts[0]

		name := parts[1]

		roles[role] = append(roles[role], name, )
	}
}


func showMenu() {

	fmt.Println("===== MENU =====",)

	fmt.Println("1. Add Person",)

	fmt.Println("2. Show Groups",)

	fmt.Println("3. Delete Person",)

	fmt.Println("4. Search Person",)

	fmt.Println("5. Save",)

	fmt.Println("6. Quit",)
}

func main() {

	reader :=
		bufio.NewReader(os.Stdin, )

	roles :=
		make(map[string][]string,)

	loadFromFile(roles, )

	for {

		showMenu()

		choice :=
			input("Choose: ", reader,)

		switch choice {

		case "1":

			addPerson(reader, roles, )

		case "2":

			showGroups(roles, )

		case "3":

			deletePerson(
				reader,
				roles,
			)

		case "4":

			searchPerson(reader, roles,)

		case "5":

			saveToFile(roles, )

		case "6":

			fmt.Println("\nProgram ended here.",)

			return

		default:

			fmt.Println("Invalid choice.\n",)
		}
	}
}
