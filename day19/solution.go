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

func (part Part) GetNewPart(category byte, predicate byte, predicateValue int, offsetMagnitude int) Part {
	applyPredicate := func(value int) int {
		if predicate == '>' {
			return max(predicateValue+offsetMagnitude, value)
		}
		return min(predicateValue-offsetMagnitude, value)
	}

	newPart := part
	switch category {
	case 'x':
		newPart.x = applyPredicate(newPart.x)
	case 'm':
		newPart.m = applyPredicate(newPart.m)
	case 'a':
		newPart.a = applyPredicate(newPart.a)
	case 's':
		newPart.s = applyPredicate(newPart.s)
	default:
		log.Fatal(fmt.Sprintf("Invalid category: %s", string(category)))
	}
	return newPart
}

type Rule struct {
	category  byte
	predicate byte
	value     int
	result    string
}

func getValidCombinations(workflowsMap map[string][]Rule, workflowName string, minPart Part, maxPart Part) int {
	if !(minPart.x <= maxPart.x && minPart.m <= maxPart.m && minPart.a <= maxPart.a && minPart.s <= maxPart.s) {
		return 0
	}
	if workflowName == "R" {
		return 0
	}
	if workflowName == "A" {
		return (maxPart.x - minPart.x + 1) * (maxPart.m - minPart.m + 1) * (maxPart.a - minPart.a + 1) * (maxPart.s - minPart.s + 1)
	}

	workflow, found := workflowsMap[workflowName]
	if !found {
		log.Fatal(fmt.Sprintf("Invalid workflow: %s", workflowName))
	}

	total := 0
	for _, rule := range workflow {
		switch rule.predicate {
		case '*':
			total += getValidCombinations(workflowsMap, rule.result, minPart, maxPart)
		case '>':
			total += getValidCombinations(workflowsMap, rule.result, minPart.GetNewPart(rule.category, rule.predicate, rule.value, 1), maxPart)
			maxPart = maxPart.GetNewPart(rule.category, '<', rule.value, 0)
		case '<':
			total += getValidCombinations(workflowsMap, rule.result, minPart, maxPart.GetNewPart(rule.category, rule.predicate, rule.value, 1))
			minPart = minPart.GetNewPart(rule.category, '>', rule.value, 0)
		}
	}
	return total
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

	validCombinations := getValidCombinations(workflowsMap, "in", Part{1, 1, 1, 1}, Part{4000, 4000, 4000, 4000})
	fmt.Println(validCombinations)
}
