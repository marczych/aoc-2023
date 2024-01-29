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

var OPERATIONAL = '.'
var DAMAGED = '#'

func getPossibilities(pattern string, groups []int) int {
	if len(groups) == 0 {
		return 0
	}

	requiredSize := len(groups) - 1
	for _, size := range groups {
		requiredSize += size
	}

	if len(pattern) < requiredSize {
		return 0
	}

	valid := true
	for _, char := range pattern[:groups[0]] {
		if char == OPERATIONAL {
			valid = false
			break
		}
	}

	// Only check if the end is damaged if we're not at the end of the string.
	if len(pattern) > groups[0] && pattern[groups[0]] == byte(DAMAGED) {
		valid = false
	}

	possibilities := 0

	if valid {
		if len(groups) == 1 {
			anyRemaining := false
			for i := groups[0] + 1; i < len(pattern); i += 1 {
				if pattern[i] == byte(DAMAGED) {
					anyRemaining = true
					break
				}
			}
			if !anyRemaining {
				// Points on the board if this is the last group and there
				// aren't any remaining damaged screws.
				possibilities += 1
			}
		} else {
			// This is a valid position for this group so now we have to check
			// the rest.
			possibilities += getPossibilities(pattern[groups[0]+1:], groups[1:])
		}
	}

	// Only check possibilities for the current group moved over by one if the
	// first position isn't definitely damaged.
	if pattern[0] != byte(DAMAGED) {
		possibilities += getPossibilities(pattern[1:], groups)
	}

	return possibilities
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	totalArrangements := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " ")
		pattern := splitLine[0]
		countStrings := strings.Split(splitLine[1], ",")
		counts := make([]int, 0, len(countStrings))
		for _, countString := range countStrings {
			count, err := strconv.Atoi(countString)
			if err != nil {
				log.Fatal("Failed to parse int: ", countString)
			}
			counts = append(counts, count)
		}
		totalArrangements += getPossibilities(pattern, counts)
	}

	fmt.Println(totalArrangements)
}
