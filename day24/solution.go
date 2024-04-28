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

const MIN_X_Y = 200000000000000
const MAX_X_Y = 400000000000000

// const MIN_X_Y = 7
// const MAX_X_Y = 27

type Hailstone struct {
	x int
	y int
	z int

	xVel int
	yVel int
	zVel int

	m float64
	b float64
}

func (hailstone Hailstone) getTime(x float64) float64 {
	return (x - float64(hailstone.x)) / float64(hailstone.xVel)
}

type Point struct {
	x float64
	y float64
}

func parseNumber(input string) int {
	number, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(fmt.Sprintf("Invalid number: %s", err))
	}

	return number
}

func intersectLines(left Hailstone, right Hailstone) Point {
	x := (right.b - left.b) / (left.m - right.m)
	return Point{x, left.m*x + left.b}
}

func main() {
	filePath := flag.String("file", "", "input file")
	flag.Parse()

	file, err := os.Open(*filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	hailstones := make([]Hailstone, 0, 1024)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		splitLine := strings.Split(line, " @ ")
		splitPosition := strings.Split(splitLine[0], ", ")
		splitVelocity := strings.Split(splitLine[1], ", ")

		xPos := parseNumber(splitPosition[0])
		yPos := parseNumber(splitPosition[1])
		zPos := parseNumber(splitPosition[2])
		xVel := parseNumber(splitVelocity[0])
		yVel := parseNumber(splitVelocity[1])
		zVel := parseNumber(splitVelocity[2])

		m := float64(yVel) / float64(xVel)
		b := float64(yPos) - m*float64(xPos)

		hailstones = append(hailstones, Hailstone{xPos, yPos, zPos, xVel, yVel, zVel, m, b})
	}

	intersectionCount := 0
	for leftIndex := 0; leftIndex < len(hailstones)-1; leftIndex += 1 {
		for rightIndex := leftIndex + 1; rightIndex < len(hailstones); rightIndex += 1 {
			left := hailstones[leftIndex]
			right := hailstones[rightIndex]
			intersection := intersectLines(left, right)
			if intersection.x >= MIN_X_Y && intersection.x <= MAX_X_Y && intersection.y >= MIN_X_Y && intersection.y <= MAX_X_Y {
				if left.getTime(intersection.x) >= 0 && right.getTime(intersection.x) >= 0 {
					intersectionCount += 1
				}
			}
		}
	}

	fmt.Println(intersectionCount)
}
