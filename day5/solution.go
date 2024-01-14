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

var EMPTY_RANGE = Range{0, 0}

type Range struct {
	start  int
	length int
}

func getSeeds(line string) []Range {
	splitLine := strings.Split(line, " ")
	// Must be odd because of the "seeds:" prefix.
	if len(splitLine)%2 != 1 {
		log.Fatal(fmt.Sprintf("Invalid seed line: %s", line))
	}

	seeds := make([]Range, 0)
	for i := range splitLine {
		if i%2 == 0 {
			continue
		}

		seeds = append(seeds, Range{parseNumber(splitLine[i]), parseNumber(splitLine[i+1])})
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

func getMin(ranges []Range) int {
	minimum := ranges[0].start

	for _, rangeValue := range ranges {
		if rangeValue.start < minimum {
			minimum = rangeValue.start
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
	fmt.Println(values)
	newValues := make([]Range, 0, len(values))

	copyUnusedToNewValues := func() {
		for _, value := range values {
			if value.length != 0 {
				newValues = append(newValues, value)
			}
		}

		values = newValues
		newValues = make([]Range, 0, len(values))
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
		sourceRange := Range{parseNumber(splitLine[1]), parseNumber(splitLine[2])}

		for i := 0; i < len(values); i++ {
			valueRange := values[i]

			if valueRange.length == 0 {
				continue
			}

			if valueRange.start < sourceRange.start && valueRange.start+valueRange.length > sourceRange.start {
				// Split the range into two.
				values = append(values, Range{valueRange.start, sourceRange.start - valueRange.start})
				valueRange = Range{sourceRange.start, valueRange.start - sourceRange.start + valueRange.length}
			}

			if valueRange.start >= sourceRange.start && valueRange.start < sourceRange.start+sourceRange.length {
				rangeSplit := min(valueRange.start+valueRange.length, sourceRange.start+sourceRange.length)
				includedRange := Range{destinationStart + valueRange.start - sourceRange.start, rangeSplit - valueRange.start}
				newValues = append(newValues, includedRange)

				extraRangeLength := valueRange.start + valueRange.length - rangeSplit
				if extraRangeLength > 0 {
					values[i] = Range{rangeSplit, extraRangeLength}
				} else {
					values[i] = EMPTY_RANGE
				}
			}
		}
	}

	copyUnusedToNewValues()
	minValue := getMin(values)
	fmt.Println(minValue)
}
