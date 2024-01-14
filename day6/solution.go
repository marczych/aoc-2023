package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
)

func solveQuadratic(a float64, b float64, c float64) (float64, float64) {
	discriminant := b*b - 4*a*c
	positiveRoot := math.Sqrt(discriminant)
	positiveValue := (positiveRoot - b) / (2 * a)
	negativeValue := (-positiveRoot - b) / (2 * a)
	return negativeValue, positiveValue
}

func getNumber(line string) int {
	number := 0
	for _, char := range line {
		if char >= '0' && char <= '9' {
			number = number*10 + int(char-'0')
		}
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

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	time := getNumber(scanner.Text())
	scanner.Scan()
	distance := getNumber(scanner.Text())

	// Add one because we need to _beat_ the winning distance.
	start, end := solveQuadratic(1, float64(-time), float64(distance+1))
	// Add one because it's an inclusive range.
	winningTimeCount := int(math.Floor(end)) - int(math.Ceil(start)) + 1
	fmt.Println(winningTimeCount)
}
