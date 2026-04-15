package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Room struct {
	Customer string
}

func (r Room) isVacant() bool {
	return r.Customer == "" || r.Customer == "Vacant"
}

//Utility 

func input(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getValidatedInt(prompt string, reader *bufio.Reader) int {
	for {
		text := input(prompt, reader)
		val, err := strconv.Atoi(text)

		if err != nil || val <= 0 {
			fmt.Println("\nInvalid input. Try again.")
			continue
		}
		return val
	}
}

//login

func logIn(reader *bufio.Reader) bool {
	users := map[string]string{
		"Emma":    "5655",
		"Justine": "8902",
		"Claude":  "1236",
		"Miller":  "3279",
		"Laura":   "4350",
	}

	fmt.Println("LOG IN")

	username := input("Enter username: ", reader)
	password := input("Pin pass: ", reader)

	if pass, ok := users[username]; ok && pass == password {
		fmt.Printf("Welcome %s!\n", username)
		return true
	}

	fmt.Println("Invalid Input.\nMake sure it was correct.")
	return false
}

//Reservation 

func handleReservation(rooms []Room, reader *bufio.Reader) {
	roomNum := getValidatedInt("Enter room number: ", reader)

	if roomNum < 1 || roomNum > len(rooms) {
		fmt.Println("\nInvalid room number.")
		return
	}

	selected := &rooms[roomNum-1]

	if !selected.isVacant() {
		fmt.Printf("Room is already occupied by %s\n", selected.Customer)
		return
	}

	name := input("Enter name: ", reader)
	selected.Customer = name

	hrs := getValidatedInt("Enter number of hours: ", reader)
	rate := getValidatedInt("Enter rate per hour: ", reader)

	total := hrs * rate
	fmt.Printf("Total Cost: P%d\n", total)

	payment := getValidatedInt("Enter payment: ", reader)

	if payment < total {
		fmt.Println("\nInsufficient payment.")
		selected.Customer = "Vacant"
		return
	}

	fmt.Printf("Change: P%d\n", payment-total)
	fmt.Println("Reservation Successful.")
}

// showw

func showRooms(rooms []Room) {
	fmt.Println("\nRoom Status:")

	for i, r := range rooms {
		name := r.Customer
		if name == "" {
			name = "Vacant"
		}
		fmt.Printf("Room %d: %s\n", i+1, name)
	}

	fmt.Println()
}


func main() {
	reader := bufio.NewReader(os.Stdin)

	//rooms
	rooms := make([]Room, 5)
	for i := range rooms {
		rooms[i].Customer = "Vacant"
	}

	for {
		if !logIn(reader) {
			continue
		}

		for {
			showRooms(rooms)

			fmt.Println("[1] Check In")
			fmt.Println("[2] Exit")

			choice := input("Choice: ", reader)

			switch choice {
			case "1":
				handleReservation(rooms, reader)
			case "2":
				fmt.Println("\nGoodbye. Program Ended.")
				return
			default:
				fmt.Println("\nInvalid Choice.")
			}
		}
	}
}
