package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type ATM struct {
	balance float64
}


func (a *ATM) checkBalance() {
	fmt.Printf("\nYour current balance is: %.2f\n", a.balance)
}

func (a *ATM) deposit(amount float64) {
	if amount <= 0 {
		fmt.Println("\nDeposit amount must be greater than zero.")
		return
	}
	a.balance += amount
	fmt.Printf("\nYou deposited %.2f. New balance: %.2f\n", amount, a.balance)
}

func (a *ATM) withdraw(amount float64) {
	if amount <= 0 {
		fmt.Println("\nWithdrawal amount must be greater than zero.")
		return
	}

	if amount > a.balance {
		fmt.Printf("\nInsufficient funds. Your balance is %.2f\n", a.balance)
		return
	}

	a.balance -= amount
	fmt.Printf("\nYou withdrew %.2f. New balance: %.2f\n", amount, a.balance)
}

//utils

func input(prompt string, reader *bufio.Reader) string {
	fmt.Print(prompt)
	text, _ := reader.ReadString('\n')
	return strings.TrimSpace(text)
}

func getFloat(prompt string, reader *bufio.Reader) (float64, bool) {
	text := input(prompt, reader)
	val, err := strconv.ParseFloat(text, 64)
	if err != nil {
		fmt.Println("Invalid input.")
		return 0, false
	}
	return val, true
}

//showcase menu

func showMenu() {
	fmt.Println(`
=========================
Welcome to Go ATM
1. Check Balance
2. Deposit
3. Withdraw
4. Exit
=========================`)
}


func main() {
	reader := bufio.NewReader(os.Stdin)
	atm := ATM{balance: 1000.0}

	for {
		showMenu()

		choice := input("Enter your choice: ", reader)

		switch choice {

		case "1":
			atm.checkBalance()

		case "2":
			if amount, ok := getFloat("Enter deposit amount: ", reader); ok {
				atm.deposit(amount)
			}

		case "3":
			if amount, ok := getFloat("Enter withdrawal amount: ", reader); ok {
				atm.withdraw(amount)
			}

		case "4":
			fmt.Println("\nProgram Done!")
			return

		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}
