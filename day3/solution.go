package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

func isSymbol(char byte) bool {
	return !((char >= '0' && char <= '9') || char == '.')
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
		currentNumber := 0
		places := 0

		resetNumbers := func() {
			currentNumber = 0
			places = 0
		}

		handleNumber := func(x int, y int) {
			if currentNumber == 0 {
				return
			}

			defer resetNumbers()

			for searchY := max(y-1, 0); searchY <= min(y+1, len(input)-1); searchY++ {
				for searchX := max(x-places, 0); searchX <= min(x+1, len(line)-1); searchX++ {
					if isSymbol(input[searchY][searchX]) {
						total += currentNumber
						return
					}
				}
			}
		}

		for x, char := range line {
			if char >= '0' && char <= '9' {
				currentNumber = currentNumber*10 + int(char-'0')
				places += 1
				continue
			}

			handleNumber(x-1, y)
		}
		handleNumber(len(line)-1, y)
	}

	fmt.Println(total)
}
