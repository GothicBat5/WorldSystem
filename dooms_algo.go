package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)



func input(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getInt(prompt string, reader *bufio.Reader) (int, bool) {
	text := input(prompt, reader)

	value, err := strconv.Atoi(text) 
	if err != nil {
		fmt.Println("Invalid number input.")
		return 0, false
	}

	return value, true
}


func main() {
	reader := bufio.NewReader(os.Stdin)

	day, ok := getInt("Day: ", reader)
	if !ok {
		return
	}

	month, ok := getInt("Month: ", reader)
	if !ok {
		return
	}

	year, ok := getInt("Year: ", reader)
	if !ok {
		return
	}


	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)


	if date.Day() != day ||
		int(date.Month()) != month ||
		date.Year() != year {

		fmt.Println("\nInvalid date entered!\n")
		return
	}

	formattedDate := date.Format("January 02, 2006")
	weekday := date.Weekday()

	fmt.Printf("\n%s was a %s.\n", formattedDate, weekday)
}
