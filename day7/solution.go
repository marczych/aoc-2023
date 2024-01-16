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

const NO_HAND_RANK = 0
const TWO_KIND_RANK = 1
const TWO_PAIR_RANK = 2
const THREE_KIND_RANK = 3
const FULL_HOUSE_RANK = 4
const FOUR_KIND_RANK = 5
const FIVE_KIND_RANK = 6

type Hand struct {
	// We can define a single integer for the value of the entire hand which is
	// a series of bytes: hand type, value of card 1, value of card 2, etc.
	HandValue int
	Bid       int
}

func parseHand(line string) Hand {
	splitLine := strings.Split(line, " ")
	return Hand{calculateHandValue(splitLine[0]), parseNumber(splitLine[1])}
}

func calculateHandValue(hand string) int {
	cardValues := map[rune]int{
		'2': 0,
		'3': 1,
		'4': 2,
		'5': 3,
		'6': 4,
		'7': 5,
		'8': 6,
		'9': 7,
		'T': 8,
		'J': 9,
		'Q': 10,
		'K': 11,
		'A': 12,
	}

	handValue := calculateHandType(hand)
	for _, card := range hand {
		handValue = handValue<<8 + cardValues[card]
	}

	return handValue
}

func calculateHandType(hand string) int {
	cardsMap := make(map[rune]int)
	for _, card := range hand {
		cardsMap[card] += 1
	}

	cardCounts := make([]int, 0, 5)
	for _, count := range cardsMap {
		// Filter out single cards because they aren't useful.
		if count > 1 {
			cardCounts = append(cardCounts, count)
		}
	}

	if len(cardCounts) == 1 {
		if cardCounts[0] == 5 {
			return FIVE_KIND_RANK
		}
		if cardCounts[0] == 4 {
			return FOUR_KIND_RANK
		}
		if cardCounts[0] == 3 {
			return THREE_KIND_RANK
		}
		if cardCounts[0] == 2 {
			return TWO_KIND_RANK
		}
	} else if len(cardCounts) == 2 {
		if cardCounts[0] == 3 || cardCounts[1] == 3 {
			return FULL_HOUSE_RANK
		} else {
			return TWO_PAIR_RANK
		}
	}

	return NO_HAND_RANK
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

	hands := make([]Hand, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		hands = append(hands, parseHand(line))
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].HandValue < hands[j].HandValue
	})

	total := 0
	for rank, hand := range hands {
		total += hand.Bid * (rank + 1)
	}

	fmt.Println(total)
}
