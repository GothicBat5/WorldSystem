package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MenuItem struct {
	Name string
	Price float64
}

func main() {

	menu := []MenuItem{
		{"Nuggets with Soup", 6.99},
		{"Ala King", 10.99},
		{"Mushroom Chip", 5.50},
		{"Caesar Salad", 7.25},
	}

	fmt.Println("Welcome\n")
	fmt.Println("Here's our menu:")

	for i, item := range menu {
		fmt.Printf("%d. %s - $%.2f\n", i+1, item.Name, item.Price)
	}

	reader := bufio.NewReader(os.Stdin)

	var total float64
	var order []MenuItem

	for {
		fmt.Print("\nEnter item number (or 'done'): ")

		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading input.")
			continue
		}

		input = strings.TrimSpace(input)

		if strings.ToLower(input) == "done" {
			break
		}

		choice, err := strconv.Atoi(input)
		if err != nil || choice < 1 || choice > len(menu) {
			fmt.Println("Invalid choice.")
			continue
		}

		item := menu[choice-1]

		order = append(order, item)
		total += item.Price

		fmt.Printf("Added %s. Current total: $%.2f\n", item.Name, total)
	}


	fmt.Println("\nYour Order:")

	if len(order) == 0 {
		fmt.Println("No items ordered.")
		return
	}

	for _, item := range order {
		fmt.Printf("- %s ($%.2f)\n", item.Name, item.Price)
	}

	fmt.Printf("\nTotal bill: $%.2f\n", total)
	fmt.Println("Gracias.")
}
