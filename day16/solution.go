package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

var TILE_MIRROR_FORWARD = '/'
var TILE_MIRROR_BACKWARD = '\\'
var TILE_SPLITTER_VERTICAL = '|'
var TILE_SPLITTER_HORIZONTAL = '-'

type Contraption struct {
	rows  []string
	beams map[int]map[int]bool
}

func createContraption() Contraption {
	contraption := Contraption{}
	contraption.beams = make(map[int]map[int]bool)
	return contraption
}

func (contraption Contraption) Get(x int, y int) rune {
	return rune(contraption.rows[y][x])
}

func (contraption *Contraption) AddRow(row string) {
	contraption.rows = append(contraption.rows, row)
}

func getTileKey(x int, y int) int {
	return x*1_000_000 + y
}

func (contraption Contraption) addBeam(x int, y int, xVel int, yVel int) {
	tileKey := getTileKey(x, y)
	_, beamExists := contraption.beams[tileKey]
	if !beamExists {
		contraption.beams[tileKey] = make(map[int]bool)
	}

	contraption.beams[tileKey][getTileKey(xVel, yVel)] = true
}

func (contraption Contraption) doesBeamExist(x int, y int, xVel int, yVel int) bool {
	tile, tileOk := contraption.beams[getTileKey(x, y)]
	if !tileOk {
		return false
	}

	_, velocityOk := tile[getTileKey(xVel, yVel)]
	return velocityOk
}

func (contraption *Contraption) GetEnergizedTileCount() int {
	contraption.fireBeam(0, 0, 1, 0)
	return len(contraption.beams)
}

func (contraption Contraption) size() int {
	// Assume it's a square.
	return len(contraption.rows)
}

func (contraption *Contraption) fireBeam(xPos int, yPos int, xVel int, yVel int) {
	for {
		if contraption.doesBeamExist(xPos, yPos, xVel, yVel) ||
			xPos < 0 ||
			yPos < 0 ||
			xPos >= contraption.size() ||
			yPos >= contraption.size() {
			break
		}
		contraption.addBeam(xPos, yPos, xVel, yVel)
		tile := contraption.Get(xPos, yPos)

		if tile == TILE_MIRROR_FORWARD {
			xVel, yVel = -yVel, -xVel
		} else if tile == TILE_MIRROR_BACKWARD {
			xVel, yVel = yVel, xVel
		} else if tile == TILE_SPLITTER_VERTICAL {
			if xVel != 0 {
				// Fire a beam up.
				contraption.fireBeam(xPos, yPos-1, 0, -1)
				// Change the course of this beam.
				xVel = 0
				yVel = 1
			}
		} else if tile == TILE_SPLITTER_HORIZONTAL {
			if yVel != 0 {
				// Fire a beam left.
				contraption.fireBeam(xPos-1, yPos, -1, 0)
				// Change the course of this beam.
				xVel = 1
				yVel = 0
			}
		}

		// Advance the current beam.
		xPos += xVel
		yPos += yVel
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

	contraption := createContraption()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		contraption.AddRow(line)
	}

	fmt.Println(contraption.GetEnergizedTileCount())
}
