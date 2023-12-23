package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var NO_NUMBER = -1

func isSymbol(char byte) bool {
	return !((char >= '0' && char <= '9') || char == '.')
}

func getNumber(char byte) int {
	if char >= '0' && char <= '9' {
		return int(char - '0')
	}

	return NO_NUMBER
}

func getGearValue(input []string, x int, y int) int {
	numbers := make([]int, 0, 2)

	// Search through the 3x3 space surrounding the gear.
	for searchY := max(y-1, 0); searchY <= min(y+1, len(input)-1); searchY++ {
		line := input[searchY]
		for searchX := max(x-1, 0); searchX <= min(x+1, len(line)-1); searchX++ {
			number := getNumber(line[searchX])
			if number == NO_NUMBER {
				continue
			}

			// Already have two numbers - this isn't a gear.
			if len(numbers) == 2 {
				return 0
			}

			// Walk back until we find a non-number or the beginning of the line.
			for searchX -= 1; searchX >= 0 && getNumber(line[searchX]) != NO_NUMBER; searchX -= 1 {
			}

			currentNumber := 0
			// Move forward by one so we're on the first number.
			for searchX += 1; searchX < len(line); searchX += 1 {
				number = getNumber(line[searchX])
				if number == NO_NUMBER {
					break
				}

				currentNumber = currentNumber*10 + number
			}

			numbers = append(numbers, currentNumber)
		}
	}

	if len(numbers) == 2 {
		return numbers[0] * numbers[1]
	}

	return 0
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	input := make([]string, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		input = append(input, scanner.Text())
	}

	total := 0

	for y, line := range input {
		for x, char := range line {
			if char != '*' {
				continue
			}
			total += getGearValue(input, x, y)
		}
	}

	fmt.Println(total)
}
