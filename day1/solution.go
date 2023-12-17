package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

const valueNotFound = -1

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

		for _, char := range line {
			if char >= '0' && char <= '9' {
				if first == valueNotFound {
					first = int(char - '0')
				}
				last = int(char - '0')
			}
		}

		total += (first * 10) + last
	}

	fmt.Println(total)
}
