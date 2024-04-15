package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

type Brick struct {
	startX int
	startY int
	startZ int
	endX   int
	endY   int
	endZ   int
	deps   map[*Brick]bool
	rdeps  map[*Brick]bool
}

func NewBrick(startX int, startY int, startZ int, endX int, endY int, endZ int) *Brick {
	return &Brick{
		startX,
		startY,
		startZ,
		endX,
		endY,
		endZ,
		make(map[*Brick]bool),
		make(map[*Brick]bool),
	}
}

func (brick *Brick) GetTotalBricksThatFall(removedBricks map[*Brick]bool) int {
	allDepsRemoved := true
	for depBrick := range brick.deps {
		_, found := removedBricks[depBrick]
		if !found {
			allDepsRemoved = false
			break
		}
	}

	if !allDepsRemoved {
		return 0
	}

	removedBricks[brick] = true
	totalBricksThatFall := 1

	for rdepBrick := range brick.rdeps {
		totalBricksThatFall += rdepBrick.GetTotalBricksThatFall(removedBricks)
	}

	return totalBricksThatFall
}

func parseNumber(input string) int {
	number, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid number: %s", err))
	}

	return number
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	bricks := make([]*Brick, 0, 1024)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		brickPieces := strings.Split(line, "~")
		startPieces := strings.Split(brickPieces[0], ",")
		endPieces := strings.Split(brickPieces[1], ",")

		bricks = append(bricks, NewBrick(
			parseNumber(startPieces[0]),
			parseNumber(startPieces[1]),
			parseNumber(startPieces[2]),
			parseNumber(endPieces[0]),
			parseNumber(endPieces[1]),
			parseNumber(endPieces[2]),
		))
	}

	sort.Slice(bricks, func(i int, j int) bool {
		return bricks[i].startZ < bricks[j].startZ
	})
	space := make(map[int]*Brick)

	coordKey := func(x int, y int, z int) int {
		return x<<40 | y<<20 | z
	}

	for _, brick := range bricks {
		iterateBrickSpace := func(callable func(int, int, int)) {
			for x := brick.startX; x <= brick.endX; x += 1 {
				for y := brick.startY; y <= brick.endY; y += 1 {
					for z := brick.startZ; z <= brick.endZ; z += 1 {
						callable(x, y, z)
					}
				}
			}
		}
		for {
			brick.startZ -= 1
			brick.endZ -= 1

			if brick.startZ == 0 {
				// This one is on the ground.
				break
			}

			iterateBrickSpace(func(x int, y int, z int) {
				collidingBrick, found := space[coordKey(x, y, z)]
				if found {
					brick.deps[collidingBrick] = true
					collidingBrick.rdeps[brick] = true
				}
			})

			if len(brick.deps) > 0 {
				break
			}
		}

		brick.startZ += 1
		brick.endZ += 1
		iterateBrickSpace(func(x int, y int, z int) {
			space[coordKey(x, y, z)] = brick
		})
	}

	totalFallingBrickCount := 0
	for _, brick := range bricks {
		removedBricks := make(map[*Brick]bool)
		removedBricks[brick] = true
		for rdepBrick := range brick.rdeps {
			totalFallingBrickCount += rdepBrick.GetTotalBricksThatFall(removedBricks)
		}
	}

	fmt.Println(totalFallingBrickCount)
}
