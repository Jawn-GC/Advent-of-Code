package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Card struct {
	ID             string
	WinningNumbers []string
	OwnedNumbers   []string
	NumMatches     int
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Parsing scratchcards...")
	cards := []Card{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		temp := strings.Split(line, ":")     // {"Card #", Numbers}
		temp2 := strings.Fields(temp[0])     // {"Card", "#"}
		temp3 := strings.Split(temp[1], "|") // {WinningNumbers, OwnedNumbers}

		id := temp2[1]
		winning_numbers := strings.Fields(temp3[0])
		owned_numbers := strings.Fields(temp3[1])

		cards = append(cards, Card{ID: id, WinningNumbers: winning_numbers, OwnedNumbers: owned_numbers})
	}

	fmt.Println("Calculating scratchcard points...")
	points_sum := 0
	for i := range cards {
		points_sum += getScratchcardPoints(&cards[i])
	}
	total_scratchcards := getTotalScratchcards(cards)

	fmt.Printf("[Part 1] Total points: %d\n", points_sum)
	fmt.Printf("[Part 2] Total scratchcards: %d\n", total_scratchcards)
}

// Assumes that a winning number is not repeated.
func getScratchcardPoints(card *Card) int {
	points := 0
	num_matches := 0

	for _, winning_num := range card.WinningNumbers {
		if isInSlice(winning_num, card.OwnedNumbers) {
			num_matches++
		}
	}

	card.NumMatches = num_matches
	if num_matches > 0 {
		points = getPower(2, num_matches-1)
	}

	return points
}

// There is intially 1 of each card. Suppose that a card is worth n points. Then, copies
// of the next n cards are created.
func getTotalScratchcards(cards []Card) int {
	total_count := 0
	card_counts := make([]int, len(cards))

	for i, card := range cards {
		card_counts[i] += 1 // There is at least one of each card
		for j := 0; j < card.NumMatches; j++ {
			index_ahead := i + j + 1
			if index_ahead < len(cards) { // OOB check
				card_counts[index_ahead] += 1 * card_counts[i]
			}
		}
	}

	for _, count := range card_counts {
		total_count += count
	}

	return total_count
}

// For n >= 0
func getPower(base, n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= base
	}
	return result
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
