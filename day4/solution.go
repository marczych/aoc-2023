package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var NO_NUMBER = -1

func getNumber(char byte) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	}

	return NO_NUMBER
}

func getWinningNumbers(line string) []int {
	winningNumbers := make([]int, 0, 10)
	currentNumber := 0

	maybeAppendNumber := func() {
		if currentNumber != 0 {
			winningNumbers = append(winningNumbers, currentNumber)
			currentNumber = 0
		}
	}

	for i := strings.Index(line, "|") + 2; i < len(line); i += 1 {
		char := line[i]
		number := getNumber(char)
		if number != NO_NUMBER {
			currentNumber = currentNumber*10 + int(char-'0')
		} else if char == ' ' {
			maybeAppendNumber()
		} else {
			panic("Invalid character: " + string(char))
		}
	}

	maybeAppendNumber()

	return winningNumbers
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	cardCopies := make(map[int]int)

	cardNumber := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cardNumber += 1
		line := scanner.Text()
		winningNumbers := getWinningNumbers(line)
		matchingNumberCount := 0

		currentNumber := 0
		for i := strings.Index(line, ":") + 2; line[i] != '|'; i += 1 {
			char := line[i]
			number := getNumber(char)
			if number != NO_NUMBER {
				currentNumber = currentNumber*10 + number
				continue
			}

			if char != ' ' {
				panic("Invalid character: " + string(char))
			}

			if slices.Contains(winningNumbers, currentNumber) {
				matchingNumberCount += 1
			}

			currentNumber = 0
		}

		for copyNumber := cardNumber + 1; copyNumber < cardNumber+1+matchingNumberCount; copyNumber += 1 {
			cardCopies[copyNumber] += cardCopies[cardNumber] + 1
		}
	}

	total := cardNumber
	for _, copyCount := range cardCopies {
		total += copyCount
	}

	fmt.Println(total)
}
