package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"sort"
	"strconv"
	"strings"
)

type Hand struct {
	Cards     []string
	Bid       int
	Hand_Type string
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Reading Camel Cards hands...")
	hands := []Hand{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Fields(line) // {Cards, Bid}

		cards := strings.Split(temp[0], "")
		bid, _ := strconv.Atoi(temp[1])
		hand_type := identifyHandType(cards)
		hands = append(hands, Hand{Cards: cards, Bid: bid, Hand_Type: hand_type})
	}

	fmt.Println("Sorting hands...")
	sortedHands := sortHands(hands)

	fmt.Println("Calculating winnings...")
	total_winnings := 0
	for i, hand := range sortedHands {
		total_winnings += (i + 1) * hand.Bid
	}

	applyJokers(&hands)
	sortedHands2 := sortHands2(hands)
	total_winnings2 := 0
	for i, hand := range sortedHands2 {
		total_winnings2 += (i + 1) * hand.Bid
	}

	fmt.Printf("[Part 1] Total winnings: %d\n", total_winnings)
	fmt.Printf("[Part 2] Total winnings: %d\n", total_winnings2)
}

func identifyHandType(cards []string) string {
	if len(cards) == 0 {
		return ""
	}

	card_count := make(map[string]int)
	for _, val := range cards {
		if _, ok := card_count[val]; !ok {
			card_count[val] = 1
		} else {
			card_count[val] += 1
		}
	}

	set_found := false // Detect Three of a kind
	num_pairs := 0     // Number of pairs found
	for _, count := range card_count {
		if count == 5 {
			return "Five of a kind"
		} else if count == 4 {
			return "Four of a kind"
		} else if count == 3 {
			set_found = true
		} else if count == 2 {
			num_pairs += 1
		}
	}

	if set_found && num_pairs == 1 {
		return "Full house"
	} else if set_found {
		return "Three of a kind"
	} else if num_pairs == 2 {
		return "Two pair"
	} else if num_pairs == 1 {
		return "One pair"
	}

	return "High card"
}

func sortHands(hands []Hand) []Hand {
	fiveOfAKind := getHandsOfType(hands, "Five of a kind")
	fourOfAKind := getHandsOfType(hands, "Four of a kind")
	fullHouse := getHandsOfType(hands, "Full house")
	threeOfAKind := getHandsOfType(hands, "Three of a kind")
	twoPair := getHandsOfType(hands, "Two pair")
	onePair := getHandsOfType(hands, "One pair")
	highCard := getHandsOfType(hands, "High card")
	sortedHands := []Hand{}

	// Secondary sorting order within each hand type group
	secondarySortHands(&fiveOfAKind)
	secondarySortHands(&fourOfAKind)
	secondarySortHands(&fullHouse)
	secondarySortHands(&threeOfAKind)
	secondarySortHands(&twoPair)
	secondarySortHands(&onePair)
	secondarySortHands(&highCard)

	// Primary sorting order of the hand type groups
	sortedHands = append(sortedHands, highCard...)
	sortedHands = append(sortedHands, onePair...)
	sortedHands = append(sortedHands, twoPair...)
	sortedHands = append(sortedHands, threeOfAKind...)
	sortedHands = append(sortedHands, fullHouse...)
	sortedHands = append(sortedHands, fourOfAKind...)
	sortedHands = append(sortedHands, fiveOfAKind...)

	return sortedHands
}

func sortHands2(hands []Hand) []Hand {
	fiveOfAKind := getHandsOfType(hands, "Five of a kind")
	fourOfAKind := getHandsOfType(hands, "Four of a kind")
	fullHouse := getHandsOfType(hands, "Full house")
	threeOfAKind := getHandsOfType(hands, "Three of a kind")
	twoPair := getHandsOfType(hands, "Two pair")
	onePair := getHandsOfType(hands, "One pair")
	highCard := getHandsOfType(hands, "High card")
	sortedHands := []Hand{}

	// Secondary sorting order within each hand type group
	secondarySortHands2(&fiveOfAKind)
	secondarySortHands2(&fourOfAKind)
	secondarySortHands2(&fullHouse)
	secondarySortHands2(&threeOfAKind)
	secondarySortHands2(&twoPair)
	secondarySortHands2(&onePair)
	secondarySortHands2(&highCard)

	// Primary sorting order of the hand type groups
	sortedHands = append(sortedHands, highCard...)
	sortedHands = append(sortedHands, onePair...)
	sortedHands = append(sortedHands, twoPair...)
	sortedHands = append(sortedHands, threeOfAKind...)
	sortedHands = append(sortedHands, fullHouse...)
	sortedHands = append(sortedHands, fourOfAKind...)
	sortedHands = append(sortedHands, fiveOfAKind...)

	return sortedHands
}

func getHandsOfType(hands []Hand, hand_type string) []Hand {
	subset := []Hand{}

	for _, hand := range hands {
		if hand.Hand_Type == hand_type {
			subset = append(subset, hand)
		}
	}

	return subset
}

func secondarySortHands(hands *[]Hand) {
	labels := []string{"2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K", "A"}
	hand_size := 5

	sort.Slice(*hands, func(i, j int) bool {
		for n := 0; n < hand_size; n++ {
			i_index := slices.Index(labels, (*hands)[i].Cards[n])
			j_index := slices.Index(labels, (*hands)[j].Cards[n])

			if i_index != j_index {
				return i_index < j_index
			}
		}

		return false
	})
}

func secondarySortHands2(hands *[]Hand) {
	labels := []string{"J", "2", "3", "4", "5", "6", "7", "8", "9", "T", "Q", "K", "A"}
	hand_size := 5

	sort.Slice(*hands, func(i, j int) bool {
		for n := 0; n < hand_size; n++ {
			i_index := slices.Index(labels, (*hands)[i].Cards[n])
			j_index := slices.Index(labels, (*hands)[j].Cards[n])

			if i_index != j_index {
				return i_index < j_index
			}
		}

		return false
	})
}

// Part 2:
// Jacks are replaced with Jokers. A Joker will become whatever label makes the hand the most powerful.
func applyJokers(hands *[]Hand) {
	for i := 0; i < len(*hands); i++ {
		num_jokers := 0
		sub_hand := []string{} // Non-Joker cards in hand
		for _, label := range (*hands)[i].Cards {
			if label == "J" {
				num_jokers += 1
			} else {
				sub_hand = append(sub_hand, label)
			}
		}

		sub_hand_type := identifyHandType(sub_hand)

		// If the number of Jokers is 5, then the hand_type stays the same ("Five of a kind")

		if num_jokers == 4 { // Sub hand type: High Card
			(*hands)[i].Hand_Type = "Five of a kind"
		}

		if num_jokers == 3 && sub_hand_type == "One pair" {
			(*hands)[i].Hand_Type = "Five of a kind"
		} else if num_jokers == 3 { // Sub hand type: High Card
			(*hands)[i].Hand_Type = "Four of a kind"
		}

		if num_jokers == 2 && sub_hand_type == "Three of a kind" {
			(*hands)[i].Hand_Type = "Five of a kind"
		} else if num_jokers == 2 && sub_hand_type == "One pair" {
			(*hands)[i].Hand_Type = "Four of a kind"
		} else if num_jokers == 2 { // Sub hand type: High Card
			(*hands)[i].Hand_Type = "Three of a kind"
		}

		if num_jokers == 1 && sub_hand_type == "Four of a kind" {
			(*hands)[i].Hand_Type = "Five of a kind"
		} else if num_jokers == 1 && sub_hand_type == "Three of a kind" {
			(*hands)[i].Hand_Type = "Four of a kind"
		} else if num_jokers == 1 && sub_hand_type == "Two pair" {
			(*hands)[i].Hand_Type = "Full house"
		} else if num_jokers == 1 && sub_hand_type == "One pair" {
			(*hands)[i].Hand_Type = "Three of a kind"
		} else if num_jokers == 1 { // Sub hand type: High Card
			(*hands)[i].Hand_Type = "One pair"
		}
	}
}
