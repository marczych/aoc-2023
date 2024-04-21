package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type Junction struct {
	x        int
	y        int
	segments []Segment
	isEnd    bool
}

type Segment struct {
	length       int
	nextJunction *Junction
}

func junctionKey(x int, y int) int {
	return x<<32 | y
}

func getTile(mapInput []string, x int, y int) byte {
	if x < 0 || x >= len(mapInput) || y < 0 || y >= len(mapInput) {
		return '#'
	}

	return mapInput[y][x]
}

func isPath(mapInput []string, x int, y int) bool {
	return getTile(mapInput, x, y) != '#'
}

func buildJunctionGraph(mapInput []string) *Junction {
	mapSize := len(mapInput)
	var startJunction *Junction
	junctionMap := make(map[int]*Junction)

	pathCount := func(x int, y int) int {
		if isPath(mapInput, x, y) {
			return 1
		}
		return 0
	}

	for y, line := range mapInput {
		for x := range line {
			if !isPath(mapInput, x, y) {
				continue
			}

			if y == 0 {
				startJunction = &Junction{x, y, nil, false}
				junctionMap[junctionKey(x, y)] = startJunction
			} else if y == mapSize-1 {
				endJunction := &Junction{x, y, nil, true}
				junctionMap[junctionKey(x, y)] = endJunction
			} else {
				connectionCount := pathCount(x, y-1) + pathCount(x, y+1) + pathCount(x-1, y) + pathCount(x+1, y)

				if connectionCount < 2 {
					log.Fatal("Invalid connection count: ", x, y, connectionCount)
				} else if connectionCount > 2 {
					junctionMap[junctionKey(x, y)] = &Junction{x, y, nil, false}
				}
			}
		}
	}

	for _, junction := range junctionMap {
		connectJunctions(mapInput, junctionMap, junction, junction.x, junction.y, -1, -1, 0)
	}

	return startJunction
}

func connectJunctions(mapInput []string, junctionMap map[int]*Junction, junction *Junction, x int, y int, previousX int, previousY int, length int) {
	if length > 0 {
		newJunction, found := junctionMap[junctionKey(x, y)]
		if found {
			junction.segments = append(junction.segments, Segment{length, newJunction})
			return
		}
	}
	if !isPath(mapInput, x, y) {
		return
	}
	next := func(newX int, newY int) {
		if !(newX == previousX && newY == previousY) {
			connectJunctions(mapInput, junctionMap, junction, newX, newY, x, y, length+1)
		}
	}
	next(x, y-1)
	next(x, y+1)
	next(x-1, y)
	next(x+1, y)
}

func findLongestPath(junction *Junction, path map[*Junction]bool, length int) int {
	if junction.isEnd {
		return length
	}

	longestPath := -1
	path[junction] = true

	for _, segment := range junction.segments {
		_, found := path[segment.nextJunction]
		if !found {
			longestPath = max(longestPath, findLongestPath(segment.nextJunction, path, length+segment.length))
		}
	}

	delete(path, junction)

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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		mapInput = append(mapInput, scanner.Text())
	}

	startJunction := buildJunctionGraph(mapInput)
	fmt.Println(findLongestPath(startJunction, make(map[*Junction]bool), 0))
}
