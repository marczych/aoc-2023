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

type Part struct {
	x int
	m int
	a int
	s int
}

func (part Part) GetSum() int {
	return part.x + part.m + part.a + part.s
}

type Rule struct {
	category  byte
	predicate byte
	value     int
	result    string
}

func (rule Rule) GetResult(part Part) string {
	if rule.category == '*' {
		return rule.result
	}

	partValue := 0
	switch rule.category {
	case 'x':
		partValue = part.x
	case 'm':
		partValue = part.m
	case 'a':
		partValue = part.a
	case 's':
		partValue = part.s
	default:
		log.Fatal(fmt.Sprintf("Invalid category: %s", string(rule.category)))
	}

	switch rule.predicate {
	case '>':
		if partValue > rule.value {
			return rule.result
		}
	case '<':
		if partValue < rule.value {
			return rule.result
		}
	}

	return ""
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

	workflowsMap := make(map[string][]Rule)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			break
		}

		openBracePosition := strings.Index(line, "{")
		workflowName := line[0:openBracePosition]
		rules := make([]Rule, 0, 10)
		rulesStrings := strings.Split(line[openBracePosition+1:], ",")

		for _, rulesString := range rulesStrings {
			endOfValuePosition := strings.Index(rulesString, ":")

			if endOfValuePosition == -1 {
				result := rulesString[0 : len(rulesString)-1]
				rules = append(rules, Rule{'*', '*', 0, result})
				break
			}

			category := rulesString[0]
			predicate := rulesString[1]
			value := parseNumber(rulesString[2:endOfValuePosition])
			result := rulesString[endOfValuePosition+1:]
			rules = append(rules, Rule{category, predicate, value, result})
		}

		workflowsMap[workflowName] = rules
	}

	total := 0
	for scanner.Scan() {
		line := scanner.Text()
		categoryStrings := strings.Split(line[1:len(line)-1], ",")
		part := Part{
			parseNumber(categoryStrings[0][2:]),
			parseNumber(categoryStrings[1][2:]),
			parseNumber(categoryStrings[2][2:]),
			parseNumber(categoryStrings[3][2:]),
		}

		result := "in"
		for {
			rules, found := workflowsMap[result]
			if !found {
				log.Fatal(fmt.Sprintf("Workflow not found: %s", result))
			}
			for _, rule := range rules {
				result = rule.GetResult(part)
				if result != "" {
					break
				}
			}

			if result == "R" {
				break
			}
			if result == "A" {
				total += part.GetSum()
				break
			}
		}
	}

	fmt.Println(total)
}
