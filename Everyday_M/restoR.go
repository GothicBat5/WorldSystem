package main

//Soon fix: make it stay in breakfast or any other choices after there original choice
//So if they choose breakfast, stay in breakfast so it's consistent.

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type MenuItem struct {
	Name  string
	Price float64
}

func main() {

	menu := map[string][]MenuItem{
		"Breakfast": {
			{"Pancakes", 5.99},
			{"Omelette", 4.50},
			{"Coffee", 2.00},
		},

		"Lunch": {
			{"Nuggets with Soup", 6.99},
			{"Ala King", 10.99},
			{"Caesar Salad", 7.25},
		},

		"Afternoon Snack": {
			{"Mushroom Chips", 5.50},
			{"Fries", 3.25},
			{"Milk Tea", 4.75},
		},

		"Dinner": {
			{"Steak", 15.99},
			{"Grilled Fish", 12.50},
			{"Pasta", 8.99},
		},
	}

	categories := []string{
		"Breakfast",
		"Lunch",
		"Afternoon Snack",
		"Dinner",
	}

	reader := bufio.NewReader(os.Stdin)

	var order []MenuItem
	var total float64

	for {

		fmt.Println("\n===== MENU =====")

		for i, cat := range categories {
			fmt.Printf("%d. %s\n", i+1, cat)
		}

		fmt.Println("0. Finish Order")

		fmt.Print("\nChoose category: ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			break
		}

		categoryChoice, err := strconv.Atoi(input)

		if err != nil ||
			categoryChoice < 1 ||
			categoryChoice > len(categories) {

			fmt.Println("Invalid category.")
			continue
		}

		selectedCategory :=
			categories[categoryChoice-1]

		items := menu[selectedCategory]

		fmt.Printf(
			"\n--- %s ---\n",
			selectedCategory,
		)

		for i, item := range items {
			fmt.Printf(
				"%d. %s - $%.2f\n",
				i+1,
				item.Name,
				item.Price,
			)
		}

		fmt.Println("0. Back")

		fmt.Print("Choose item: ")

		input, _ = reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "0" {
			continue
		}

		itemChoice, err := strconv.Atoi(input)

		if err != nil ||
			itemChoice < 1 ||
			itemChoice > len(items) {

			fmt.Println("Invalid item.")
			continue
		}

		selected := items[itemChoice-1]

		order = append(order, selected)

		total += selected.Price

		fmt.Printf(
			"Added %s ($%.2f)\n",
			selected.Name,
			selected.Price,
		)

		fmt.Printf(
			"Current total: $%.2f\n",
			total,
		)
	}

	fmt.Println("\n===== ORDER =====")

	if len(order) == 0 {
		fmt.Println("No items ordered.")
		return
	}

	for _, item := range order {
		fmt.Printf(
			"- %s ($%.2f)\n",
			item.Name,
			item.Price,
		)
	}

	fmt.Printf(
		"\nTotal bill: $%.2f\n",
		total,
	)

	fmt.Println("Gracias.")
}
