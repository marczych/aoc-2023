package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func solveQuadratic(a float64, b float64, c float64) (float64, float64) {
	discriminant := b*b - 4*a*c
	positiveRoot := math.Sqrt(discriminant)
	positiveValue := (positiveRoot - b) / (2 * a)
	negativeValue := (-positiveRoot - b) / (2 * a)
	return negativeValue, positiveValue
}

func getNumbers(line string) []int {
	splitLine := strings.Split(line, " ")

	numbers := make([]int, 0)
	for i := range splitLine {
		token := splitLine[i]
		if i == 0 || token == "" {
			// Skip the prefix.
			continue
		}

		numbers = append(numbers, parseNumber(token))
	}

	return numbers
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

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	times := getNumbers(scanner.Text())
	scanner.Scan()
	distances := getNumbers(scanner.Text())

	if len(times) != len(distances) {
		log.Fatal(
			fmt.Sprintf("Mismatched length of times and distances: %d vs. %d",
				len(times),
				len(distances),
			),
		)
	}

	total := 1
	for i := range times {
		time := times[i]
		distance := distances[i]
		// Add one because we need to _beat_ the winning distance.
		start, end := solveQuadratic(1, float64(-time), float64(distance+1))
		// Add one because it's an inclusive range.
		winningTimeCount := int(math.Floor(end)) - int(math.Ceil(start)) + 1
		total = total * winningTimeCount
	}

	fmt.Println(total)
}
