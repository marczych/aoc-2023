package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"maps"
	"os"
	"strings"
)

func pathKey(x int, y int) int {
	return x<<32 | y
}

func printMap(mapInput []string, path map[int]bool) {
	fmt.Println("MAP:")
	for y, line := range mapInput {
		for x, tile := range line {
			if _, found := path[pathKey(x, y)]; found {
				fmt.Print("O")
			} else {
				fmt.Print(string(tile))
			}
		}
		fmt.Println()
	}
}

func getLongestPath(mapInput []string, path map[int]bool, x int, y int) int {
	key := pathKey(x, y)
	_, found := path[key]
	if x < 0 || x >= len(mapInput) || y < 0 || y >= len(mapInput) || found || mapInput[y][x] == '#' {
		return -1
	}

	currentTile := mapInput[y][x]

	if y == len(mapInput) - 1 && currentTile == '.' {
		return len(path)
	}

	// Copy on write.
	path = maps.Clone(path)
	path[key] = true

	switch currentTile {
	case '^':
		return getLongestPath(mapInput, path, x, y-1)
	case 'v':
		return getLongestPath(mapInput, path, x, y+1)
	case '<':
		return getLongestPath(mapInput, path, x-1, y)
	case '>':
		return getLongestPath(mapInput, path, x+1, y)
	}

	longestPath := -1
	longestPath = max(longestPath, getLongestPath(mapInput, path, x, y-1))
	longestPath = max(longestPath, getLongestPath(mapInput, path, x, y+1))
	longestPath = max(longestPath, getLongestPath(mapInput, path, x-1, y))
	longestPath = max(longestPath, getLongestPath(mapInput, path, x+1, y))

	return longestPath
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mapInput := make([]string, 0, 1024)
	startingXpos := -1
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mapInput = append(mapInput, scanner.Text())
		if startingXpos == -1 {
			startingXpos = strings.Index(mapInput[0], ".")
		}
	}

	longestPath := getLongestPath(mapInput, make(map[int]bool), startingXpos, 0)
	fmt.Println(longestPath)
}
