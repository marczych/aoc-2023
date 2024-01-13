package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

var CONSUMED_NUMBER = -1

func getSeeds(line string) []int {
	splitLine := strings.Split(line, " ")
	seeds := make([]int, 0, len(splitLine)-1)
	for i, seedString := range splitLine {
		if i == 0 {
			continue
		}
		seeds = append(seeds, parseNumber(seedString))
	}

	return seeds
}

func parseNumber(input string) int {
	number, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid number: %s", err))
	}

	return number
}

func getMin(values []int) int {
	minimum := values[0]

	for _, value := range values {
		if value < minimum {
			minimum = value
		}
	}

	return minimum
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	scanner.Scan()
	values := getSeeds(scanner.Text())
	newValues := make([]int, 0, len(values))

	copyUnusedToNewValues := func() {
		for _, value := range values {
			if value != CONSUMED_NUMBER {
				newValues = append(newValues, value)
			}
		}

		values = newValues
		newValues = make([]int, 0, len(values))
	}

	// Start the loop to convert the existing values to the next "phase" e.g.
	// seed => soil => fertilizer => etc.
	for scanner.Scan() {
		line := scanner.Text()

		if line == "" {
			// Throwaway the next header line which is useless.
			scanner.Scan()

			copyUnusedToNewValues()
			continue
		}

		splitLine := strings.Split(line, " ")
		if len(splitLine) != 3 {
			log.Fatal(fmt.Sprintf("Invalid line: %s", line))
		}
		destinationStart := parseNumber(splitLine[0])
		sourceStart := parseNumber(splitLine[1])
		sourceEnd := sourceStart + parseNumber(splitLine[2])

		for i, value := range values {
			if value >= sourceStart && value <= sourceEnd {
				newValues = append(newValues, destinationStart+(value-sourceStart))
				values[i] = CONSUMED_NUMBER
			}
		}
	}

	copyUnusedToNewValues()
	minValue := getMin(values)
	fmt.Println(minValue)
}
