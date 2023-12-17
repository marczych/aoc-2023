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

func getGamePower(gameLine string) int {
	splitLine := strings.Split(gameLine, ": ")
	if len(splitLine) != 2 {
		log.Fatal(fmt.Sprintf("Invalid game: %s", gameLine))
	}

	handfuls := strings.Split(splitLine[1], "; ")
	colorCounts := map[string]int{
		"red":   0,
		"green": 0,
		"blue":  0,
	}

	for _, handful := range handfuls {
		cubes := strings.Split(handful, ", ")
		for _, cubeString := range cubes {
			cubeParts := strings.Split(cubeString, " ")
			if len(cubeParts) != 2 {
				log.Fatal(fmt.Sprintf("Invalid cubes: %s", handful))
			}

			count, err := strconv.Atoi(cubeParts[0])
			if err != nil {
				log.Fatal(fmt.Sprintf("Invalid count: %s", err))
			}

			color := cubeParts[1]
			maxCubeCount, ok := colorCounts[color]

			if !ok {
				log.Fatal(fmt.Sprintf("Invalid color: %s", color))
			}

			colorCounts[color] = max(maxCubeCount, count)
		}
	}

	power := 1
	for _, count := range colorCounts {
		power *= count
	}
	return power
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
		line := scanner.Text()
		total += getGamePower(line)
	}

	fmt.Println(total)
}
