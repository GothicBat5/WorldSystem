
package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"unicode"
)


func isInteger(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsDigit(r) {
			return false
		}
	}
	return true
}

func isDecimal(s string) bool {
	if s == "" {
		return false
	}

	dotSeen := false

	for _, r := range s {
		if r == '.' {
			if dotSeen {
				return false
			}
			dotSeen = true
		} else if !unicode.IsDigit(r) {
			return false
		}
	}

	return dotSeen //must contain exactly one dot
}

func isSingleChar(s string) bool {
	runes := []rune(s)
	return len(runes) == 1 && unicode.IsLetter(runes[0])
}

func isText(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}


func main() {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Input: ")
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	tokens := strings.Fields(line)

	hasInt := false
	hasDec := false
	hasText := false
	hasChar := false
	hasMixed := false

	for _, token := range tokens {

		switch {
		case isInteger(token):
			fmt.Println(token, "-> Integer")
			hasInt = true

		case isDecimal(token):
			fmt.Println(token, "-> Decimal")
			hasDec = true

		case isSingleChar(token):
			fmt.Println(token, "-> Single Character")
			hasChar = true

		case isText(token):
			fmt.Println(token, "-> Text")
			hasText = true

		default:
			fmt.Println(token, "-> Mixed")
			hasMixed = true
		}
	}

	fmt.Println("\nSummary:")

	if hasInt {
		fmt.Println("- Contains integers")
	}
	if hasDec {
		fmt.Println("- Contains decimals")
	}
	if hasText {
		fmt.Println("- Contains text")
	}
	if hasChar {
		fmt.Println("- Contains single characters")
	}
	if hasMixed {
		fmt.Println("- Contains mixed values")
	}
}
