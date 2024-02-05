package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"log"
	"os"
)

var EMPTY_SPACE = '.'
var CUBE_ROCK = '#'
var ROUND_ROCK = 'O'
var ITERATION_COUNT = 1000000000

type MapData struct {
	data []rune
	size int
}

func (mapData MapData) Serialize() string {
	return string(mapData.data)
}

func (mapData *MapData) AddRow(row string) {
	for _, char := range row {
		mapData.data = append(mapData.data, char)
	}
}

func (mapData *MapData) Rotate() {
	for layer := 0; layer < mapData.size/2; layer += 1 {
		x := layer
		for y := layer; y < mapData.size-layer-1; y += 1 {
			original := mapData.Get(x, y)
			mapData.Set(x, y, mapData.Get(y, mapData.size-x-1))
			mapData.Set(y, mapData.size-x-1, mapData.Get(mapData.size-x-1, mapData.size-y-1))
			mapData.Set(mapData.size-x-1, mapData.size-y-1, mapData.Get(mapData.size-y-1, x))
			mapData.Set(mapData.size-y-1, x, original)
		}
	}
}

func (mapData MapData) Get(x int, y int) rune {
	return mapData.data[x+y*mapData.size]
}

func (mapData MapData) Set(x int, y int, value rune) {
	mapData.data[x+y*mapData.size] = value
}

func (mapData MapData) CalculateLoad() int {
	totalLoad := 0
	for x := 0; x < mapData.size; x += 1 {
		for y := 0; y < mapData.size; y += 1 {
			if mapData.Get(x, y) == ROUND_ROCK {
				totalLoad += mapData.size - y
			}
		}
	}
	return totalLoad
}

func (mapData MapData) ToString() string {
	buffer := bytes.Buffer{}
	for i := 0; i < mapData.size*mapData.size; i += 1 {
		buffer.WriteRune(mapData.data[i])
		if i%mapData.size == mapData.size-1 {
			buffer.WriteRune('\n')
		}
	}

	return buffer.String()
}

func (mapData *MapData) RollNorth() {
	for x := 0; x < mapData.size; x += 1 {
		nextEmptySpace := 0
		for y := 0; y < mapData.size; y += 1 {
			switch mapData.Get(x, y) {
			case CUBE_ROCK:
				nextEmptySpace = y + 1
			case ROUND_ROCK:
				if y != nextEmptySpace {
					mapData.Set(x, nextEmptySpace, ROUND_ROCK)
					mapData.Set(x, y, EMPTY_SPACE)
				}
				nextEmptySpace += 1
			}
		}
	}
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	mapData := MapData{}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		mapData.size = len(line)
		mapData.AddRow(line)
	}

	cache := make(map[string]int)
	for i := 1; i < ITERATION_COUNT; i += 1 {
		for r := 0; r < 4; r += 1 {
			mapData.RollNorth()
			mapData.Rotate()
		}

		mapKey := mapData.Serialize()
		previousIteration, ok := cache[mapKey]

		if ok {
			period := i - previousIteration
			if (ITERATION_COUNT-i)%period == 0 {
				break
			}
		} else {
			cache[mapKey] = i
		}
	}
	fmt.Println(mapData.CalculateLoad())
}
