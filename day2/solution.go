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

var colorCounts = map[string]int{
	"red":   12,
	"green": 13,
	"blue":  14,
}

func getValidGameValue(gameLine string) int {
	splitLine := strings.Split(gameLine, ": ")
	if len(splitLine) != 2 {
		log.Fatal(fmt.Sprintf("Invalid game: %s", gameLine))
	}

	gameNumber := getGameNumber(splitLine[0])

	if isValidGame(splitLine[1]) {
		return gameNumber
	}

	return 0
}

func getGameNumber(gameString string) int {
	splitGameString := strings.Split(gameString, " ")
	if len(splitGameString) != 2 {
		log.Fatal(fmt.Sprintf("Invalid game prefix: %s", gameString))
	}

	gameNumber, err := strconv.Atoi(splitGameString[1])
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid game number: %s", err))
	}

	return gameNumber
}

func isValidGame(game string) bool {
	handfuls := strings.Split(game, "; ")
	for _, handful := range handfuls {
		if !isValidHandful(handful) {
			return false
		}
	}

	return true
}

func isValidHandful(handful string) bool {
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

		if count > maxCubeCount {
			return false
		}
	}

	return true
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
		total += getValidGameValue(line)
	}

	fmt.Println(total)
}
