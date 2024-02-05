package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func hash(input string) int {
	currentValue := 0

	for _, char := range input {
		currentValue += int(char)
		currentValue *= 17
		currentValue %= 256
	}

	return currentValue
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
	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		steps := strings.Split(line, ",")
		for _, step := range steps {
			total += hash(step)
		}
	}

	fmt.Println(total)
}
