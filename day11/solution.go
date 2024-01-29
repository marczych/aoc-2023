package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
)

type Galaxy struct {
	x int
	y int
}

func abs(value int) int {
	if value < 0 {
		return -value
	}

	return value
}

func (g Galaxy) Distance(other Galaxy) int {
	return abs(g.x-other.x) + abs(g.y-other.y)
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	galaxies := make([]Galaxy, 0)
	xGalaxies := make(map[int]bool, 0)
	scanner := bufio.NewScanner(file)
	y := 0
	for scanner.Scan() {
		line := scanner.Text()
		foundGalaxy := false
		for x, char := range line {
			if char == '#' {
				foundGalaxy = true
				xGalaxies[x] = true
				galaxies = append(galaxies, Galaxy{x, y})
			}
		}
		if !foundGalaxy {
			y += 1
		}
		y += 1
	}

	for i, galaxy := range galaxies {
		xOffset := 0
		for x := 0; x < galaxy.x; x += 1 {
			if _, ok := xGalaxies[x]; !ok {
				xOffset += 1
			}
		}
		if xOffset != 0 {
			galaxies[i].x += xOffset
		}
	}

	totalDistance := 0
	for i, galaxy := range galaxies {
		for j := i + 1; j < len(galaxies); j += 1 {
			totalDistance += galaxy.Distance(galaxies[j])
		}
	}

	fmt.Println(totalDistance)
}
