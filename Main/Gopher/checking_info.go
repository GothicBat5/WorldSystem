package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

var roles = make(map[string][]string)

//Utility 

func trim(s string) string {
	return strings.TrimSpace(s)
}

func capitalize(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func input(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return trim(text)
}

//core

func addPerson(reader *bufio.Reader) {
	name := input("Enter name: ", reader)

	if name == "" {
		fmt.Println("Name cannot be empty.\n")
		return
	}

	role := capitalize(input("Enter role: ", reader))

	if role == "" {
		fmt.Println("Role cannot be empty.\n")
		return
	}

	roles[role] = append(roles[role], name)

	fmt.Printf("Added %s as %s\n\n", name, role)
}

func showGroups() {
	fmt.Println("\n------ People Groups ------")

	if len(roles) == 0 {
		fmt.Println("No group found.\n")
		return
	}

	var sortedRoles []string
	for role := range roles {
		sortedRoles = append(sortedRoles, role)
	}
	sort.Strings(sortedRoles)

	for _, role := range sortedRoles {
		people := roles[role]
		sort.Strings(people)
		fmt.Printf("%ss: %s\n", role, strings.Join(people, ", "))
	}

	fmt.Println()
}

func deletePerson(reader *bufio.Reader) {
	name := input("Enter name to delete: ", reader)

	for role, people := range roles {
		for i, p := range people {
			if strings.EqualFold(p, name) {
				roles[role] = append(people[:i], people[i+1:]...)
				fmt.Printf("Removed %s from %s\n\n", p, role)
				return
			}
		}
	}

	fmt.Println("Person not found.\n")
}

func searchPerson(reader *bufio.Reader) {
	name := input("Search name: ", reader)
	found := false

	for role, people := range roles {
		for _, person := range people {
			if strings.Contains(strings.ToLower(person), strings.ToLower(name)) {
				fmt.Printf("%s is a %s\n", person, role)
				found = true
			}
		}
	}

	if !found {
		fmt.Println("No person found.\n")
	} else {
		fmt.Println()
	}
}

//file handling

func saveToFile() {
	file, err := os.Create("roles.txt")
	if err != nil {
		fmt.Println("Error saving file.\n")
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for role, people := range roles {
		for _, person := range people {
			fmt.Fprintf(writer, "%s|%s\n", role, person)
		}
	}

	writer.Flush()
	fmt.Println("Data saved!\n")
}

func loadFromFile() {
	file, err := os.Open("roles.txt")
	if err != nil {
		return //file doesn't exist = ignore
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "|")

		if len(parts) == 2 {
			role := parts[0]
			name := parts[1]

			roles[role] = append(roles[role], name)
		}
	}
}


func showMenu() {
	fmt.Println("===== MENU =====")
	fmt.Println("1. Add Person")
	fmt.Println("2. Show Groups")
	fmt.Println("3. Delete Person")
	fmt.Println("4. Search Person")
	fmt.Println("5. Save")
	fmt.Println("6. Quit")
}

//m loop

func main() {
	reader := bufio.NewReader(os.Stdin)

	loadFromFile()

	for {
		showMenu()
		choice := input("Choose: ", reader)

		switch choice {
		case "1":
			addPerson(reader)
		case "2":
			showGroups()
		case "3":
			deletePerson(reader)
		case "4":
			searchPerson(reader)
		case "5":
			saveToFile()
		case "6":
			fmt.Println("\nProgram ended here.")
			return
		default:
			fmt.Println("Invalid choice.\n")
		}
	}
}
