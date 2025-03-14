package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Obtaining towel types...\n")
	input_type := "towels"
	towels := []string{}
	designs := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			input_type = "designs"
			continue
		}

		if input_type == "towels" {
			towels = append(towels, strings.Split(line, ", ")...)
		} else if input_type == "designs" {
			designs = append(designs, line)
		}
	}

	fmt.Printf("Determining possible designs...\n")
	num_possible_designs := 0
	num_valid_arrangements := 0
	for _, design := range designs {
		if isDesignPossible(towels, design) {
			num_possible_designs++
			num_valid_arrangements += numValidArrangements(towels, design)
		}
	}

	fmt.Printf("Number of possible designs: %d\n", num_possible_designs)
	fmt.Printf("Number of possible arrangements: %d\n", num_valid_arrangements)
}

func isDesignPossible(towels []string, design string) bool {
	memo := make(map[string]bool)

	var helper func(string) bool
	helper = func(remaining string) bool {
		if remaining == "" {
			return true
		}

		if val, ok := memo[remaining]; ok {
			return val
		}

		for _, towel := range towels {
			if len(remaining) >= len(towel) && towel == remaining[:len(towel)] {
				if helper(remaining[len(towel):]) {
					memo[remaining] = true
					return true
				}
			}
		}

		memo[remaining] = false
		return false
	}

	return helper(design)
}

func numValidArrangements(towels []string, design string) int {
	memo := make(map[string]int)

	var helper func(string) int
	helper = func(remaining string) int {
		if remaining == "" {
			return 1
		}

		if val, ok := memo[remaining]; ok {
			return val
		}

		count := 0
		for _, towel := range towels {
			if len(remaining) >= len(towel) && towel == remaining[:len(towel)] {
				count += helper(remaining[len(towel):])
			}
		}

		memo[remaining] = count
		return count
	}

	return helper(design)
}
