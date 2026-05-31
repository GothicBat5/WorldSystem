package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const ROOM_COUNT = 5

type Room struct {
	Customer string
	Occupied bool
}


func input(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getValidatedInt(prompt string, reader *bufio.Reader) int {

	for {

		text := input(prompt, reader)

		value, err := strconv.Atoi(text)

		if err != nil || value <= 0 {

			fmt.Println("Invalid input.")
			continue
		}

		return value
	}
}

func getValidatedFloat(prompt string, reader *bufio.Reader) float64 {

	for {

		text := input(prompt, reader)

		value, err := strconv.ParseFloat(text, 64)

		if err != nil || value <= 0 {

			fmt.Println("Invalid input.")
			continue
		}

		return value
	}
}


func (r Room) isVacant() bool {
	return !r.Occupied
}


func logIn(reader *bufio.Reader, users map[string]string) bool {

	const MAX_ATTEMPTS = 3

	for attempt := 1; attempt <= MAX_ATTEMPTS; attempt++ {

		fmt.Printf("\nLOG IN (%d/%d)\n", attempt, MAX_ATTEMPTS)

		username := input("Enter username: ", reader)
		password := input("Pin pass: ", reader)

		if pass, ok := users[username]; ok && pass == password {

			fmt.Printf("Welcome %s!\n", username)
			return true
		}

		fmt.Println("Invalid credentials.")
	}

	fmt.Println("\nToo many failed attempts.")
	return false
}


func handleReservation(rooms []Room, reader *bufio.Reader) {

	roomNum := getValidatedInt("Enter room number: ", reader)

	if roomNum < 1 || roomNum > len(rooms) {

		fmt.Println("Invalid room number.")
		return
	}

	selected := &rooms[roomNum-1]

	if !selected.isVacant() {

		fmt.Printf("Room already occupied by %s\n", selected.Customer,)

		return
	}

	name := input("Enter customer name: ", reader)

	if name == "" {

		fmt.Println("Name required.")
		return
	}

	hours := getValidatedInt("Enter hours: ", reader)
	rate := getValidatedFloat("Rate per hour: ", reader)

	total := float64(hours) * rate

	fmt.Printf("Total Cost: P%.2f\n", total)

	payment := getValidatedFloat("Payment: ", reader)

	if payment < total {

		fmt.Println("Insufficient payment.")
		return
	}

	selected.Customer = name
	selected.Occupied = true

	fmt.Printf("Change: P%.2f\n", payment-total)
	fmt.Println("Reservation Successful.")
}


func checkOut(rooms []Room, reader *bufio.Reader) {

	roomNum := getValidatedInt("Checkout room number: ", reader)

	if roomNum < 1 || roomNum > len(rooms) {

		fmt.Println("Invalid room.")
		return
	}

	selected := &rooms[roomNum-1]

	if selected.isVacant() {

		fmt.Println("Room already vacant.")
		return
	}

	fmt.Printf(
		"%s checked out from room %d\n",
		selected.Customer,
		roomNum,
	)

	selected.Customer = ""
	selected.Occupied = false
}


func showRooms(rooms []Room) {

	fmt.Println("\n===== ROOM STATUS =====")

	for i, room := range rooms {

		status := "Vacant"

		if room.Occupied {
			status = room.Customer
		}

		fmt.Printf("Room %d : %s\n", i+1, status,)
	}

	fmt.Println()
}


func main() {

	reader := bufio.NewReader(os.Stdin)

	users := map[string]string{
		"Emma": "5655",
		"Justine": "8902",
		"Claude": "1236",
		"Miller": "3279",
		"Laura": "4350",
	}

	rooms := make([]Room, ROOM_COUNT)

	if !logIn(reader, users) {
		return
	}

	for {

		showRooms(rooms)

		fmt.Println("[1] Check In")
		fmt.Println("[2] Check Out")
		fmt.Println("[3] Exit")

		choice := input("Choice: ", reader)

		switch choice {

		case "1":
			handleReservation(rooms, reader)

		case "2":
			checkOut(rooms, reader)

		case "3":
			fmt.Println("\nProgram Ended.")
			return

		default:
			fmt.Println("Invalid choice.")
		}
	}
}
