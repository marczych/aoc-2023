package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const valueNotFound = -1

var numberWords = [...]string{
	"zero",
	"one",
	"two",
	"three",
	"four",
	"five",
	"six",
	"seven",
	"eight",
	"nine",
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	total := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		first := valueNotFound
		last := valueNotFound
		line := scanner.Text()
		lineLength := len(line)

		for linePosition, char := range line {
			currentNumber := valueNotFound

			if char >= '0' && char <= '9' {
				currentNumber = int(char - '0')
			} else {
				charactersRemaining := lineLength - linePosition

			WORD_LOOP:
				for number, word := range numberWords {
					if len(word) > charactersRemaining {
						continue
					}

					for charPosition, numberChar := range word {
						if numberChar != rune(line[linePosition+charPosition]) {
							continue WORD_LOOP
						}
					}

					// We found the number!
					currentNumber = number
					break
				}
			}

			if currentNumber != valueNotFound {
				if first == valueNotFound {
					first = currentNumber
				}
				last = currentNumber
			}
		}

		if first != valueNotFound && last != valueNotFound {
			total += (first * 10) + last
		}
	}

	fmt.Println(total)
}
