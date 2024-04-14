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

const START = 'S'
const PLOT = '.'
const ROCK = '#'

func getGardenPositions(garden []string, state map[int]bool, endingPositions map[int]bool, xPos int, yPos int, steps int) {
	checkPosition := func(x int, y int) {
		if x < 0 || y < 0 || x >= len(garden) || y >= len(garden) {
			return
		}

		stateKey := (x << 48) | (y << 32) | steps
		_, found := state[stateKey]
		if found {
			return
		}
		state[stateKey] = true

		if garden[y][x] != ROCK {
			if steps == 1 {
				endingPositions[(x<<48)|(y<<32)] = true
			} else {
				getGardenPositions(garden, state, endingPositions, x, y, steps-1)
			}
		}
	}

	checkPosition(xPos-1, yPos)
	checkPosition(xPos+1, yPos)
	checkPosition(xPos, yPos-1)
	checkPosition(xPos, yPos+1)
}

func main() {
	filePath := flag.String("file", "", "input file")
	stepsFlag := flag.String("steps", "", "number of steps")
	flag.Parse()

	steps, err := strconv.Atoi(*stepsFlag)
	if err != nil {
		log.Fatal(err)
	}

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	var garden []string
	xPos := -1
	yPos := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		garden = append(garden, line)
		if xPos == -1 {
			startingPosition := strings.Index(line, string(START))
			if startingPosition != -1 {
				xPos = startingPosition
				yPos = len(garden) - 1
			}
		}
	}

	endingPositions := make(map[int]bool)
	getGardenPositions(garden, make(map[int]bool), endingPositions, xPos, yPos, steps)
	fmt.Println(len(endingPositions))
}
