package main

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

type CartItem struct {
	Item MenuItem
	Quantity int
}

func readInput(reader *bufio.Reader, prompt string) string {
	fmt.Print(prompt)

	text, _ := reader.ReadString('\n')

	return strings.TrimSpace(text)
}

func readInt(reader *bufio.Reader, prompt string) (int, bool) {

	input := readInput(reader, prompt)

	value, err := strconv.Atoi(input)

	if err != nil {
		return 0, false
	}

	return value, true
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

		"Desserts": {
			{"Ice Cream", 3.99},
			{"Cheesecake", 4.50},
			{"Brownies", 2.99},
		},

		"Drinks": {
			{"Cola", 1.99},
			{"Orange Juice", 2.50},
			{"Water", 1.00},
		},
	}

	categories := []string{
		"Breakfast",
		"Lunch",
		"Afternoon Snack",
		"Dinner",
		"Desserts",
		"Drinks",
	}

	reader := bufio.NewReader(os.Stdin)

	var cart []CartItem
	var total float64

CATEGORY_LOOP:

	for {

		fmt.Println("\n===== CATEGORIES =====")

		for i, cat := range categories {
			fmt.Printf("%d. %s\n", i+1, cat)
		}

		fmt.Println("7. View Cart")
		fmt.Println("0. Finish")

		categoryChoice, ok :=
			readInt(reader, "\nChoose category: ")

		if !ok {

			fmt.Println("Invalid input.")
			continue
		}

		switch categoryChoice {

		case 0:
			break CATEGORY_LOOP

		case 7:

			fmt.Println("\n===== CART =====")

			if len(cart) == 0 {

				fmt.Println("Cart empty.")

			} else {

				for _, c := range cart {

					sub := c.Item.Price * float64(c.Quantity)

					fmt.Printf("%dx %s = $%.2f\n",
						c.Quantity,
						c.Item.Name,
						sub,
					)
				}

				fmt.Printf("\nTotal: $%.2f\n",total,)
			}

			continue
		}

		if categoryChoice < 1 ||
			categoryChoice > len(categories) {

			fmt.Println("Invalid category.")
			continue
		}

		selectedCategory :=	categories[categoryChoice-1]

		items := menu[selectedCategory]


		for {

			fmt.Printf(
				"\n=== %s ===\n",
				selectedCategory,
			)

			for i, item := range items {

				fmt.Printf("%d. %s ($%.2f)\n",
					i+1,
					item.Name,
					item.Price,
				)
			}

			fmt.Println("0. Back")

			itemChoice, ok :=
				readInt(reader, "\nChoose item: ")

			if !ok {

				fmt.Println("Invalid input.")
				continue
			}

			if itemChoice == 0 {

				break
			}

			if itemChoice < 1 || itemChoice > len(items) {

				fmt.Println("Invalid item.")
				continue
			}

			qty, ok :=
				readInt(reader, "Quantity: ")

			if !ok || qty <= 0 {

				fmt.Println("Invalid quantity.")
				continue
			}

			selected :=	items[itemChoice-1]

			cart = append(
				cart,
				CartItem{
					Item: selected,
					Quantity: qty,
				},
			)

			subtotal :=	selected.Price * float64(qty)

			total += subtotal

			fmt.Printf("\nAdded %dx %s\n",	qty,selected.Name,)

			fmt.Printf("Subtotal: $%.2f\n",subtotal,)

			fmt.Printf("Running total: $%.2f\n",total,)
		}
	}

	fmt.Println("\n===== RECEIPT =====")

	if len(cart) == 0 {

		fmt.Println("No items ordered.")
		return
	}

	for _, c := range cart {

		sub :=
			c.Item.Price *
				float64(c.Quantity)

		fmt.Printf("%dx %s = $%.2f\n",
			c.Quantity,
			c.Item.Name,
			sub,
		)
	}

	fmt.Printf("\nTOTAL BILL: $%.2f\n",
		total,
	)

	fmt.Println("Gracias.")
}
