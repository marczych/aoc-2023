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
		// Now a Joker which is least valuable on its own.
		'J': 0,
		'2': 1,
		'3': 2,
		'4': 3,
		'5': 4,
		'6': 5,
		'7': 6,
		'8': 7,
		'9': 8,
		'T': 9,
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

	// Count the jokers and remove them from the map because we're going to handle them specially.
	jokerCount := cardsMap['J']
	delete(cardsMap, 'J')
	cardCounts := make([]int, 0, 5)
	for _, count := range cardsMap {
		// Filter out single cards because they aren't useful.
		if count > 1 {
			cardCounts = append(cardCounts, count)
		}
	}

	sort.Slice(cardCounts, func(i, j int) bool {
		return cardCounts[i] > cardCounts[j]
	})

	// Special case for all jokers.
	if jokerCount == 5 {
		return FIVE_KIND_RANK
	}

	if jokerCount > 0 {
		// If we have no duplicate cards but we have jokers, we can make some duplicates.
		if len(cardCounts) == 0 {
			cardCounts = append(cardCounts, jokerCount+1)
		} else {
			// Otherwise, each joker adds another of the card we have the most.
			cardCounts[0] += jokerCount
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
		if cardCounts[0] == 3 {
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
