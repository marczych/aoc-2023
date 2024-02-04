package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var CUBE_ROCK = '#'
var ROUND_ROCK = 'O'

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	totalLoad := 0
	heights := make([]int, 1000)
	// Assume the input is a square.
	size := 0
	// We want the first loop to start at 0.
	currentHeight := -1

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		size = len(line)
		currentHeight += 1

		for position, char := range line {
			switch char {
			case CUBE_ROCK:
				heights[position] = currentHeight + 1
			case ROUND_ROCK:
				totalLoad += size - heights[position]
				heights[position] += 1
			}
		}
	}

	fmt.Println(totalLoad)
}
