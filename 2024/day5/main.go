package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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

	fmt.Print("Parsing ordering rules...\n")
	ordering_rules := map[string][]string{}
	orderingRulesRead := false
	updates := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			orderingRulesRead = true
			fmt.Print("Parsing updates...\n")
			continue
		}
		if !orderingRulesRead {
			pair := strings.Split(line, "|")
			key := pair[0]
			if _, ok := ordering_rules[key]; !ok {
				ordering_rules[key] = []string{}
			}
			ordering_rules[key] = append(ordering_rules[key], pair[1])
		} else {
			update := strings.Split(line, ",")
			updates = append(updates, update)
		}
	}

	fmt.Print("Validating updates...\n")
	var isUpdateValid bool
	middle_sum := 0
	fixed_middle_sum := 0
	for _, update := range updates {
		size := len(update)
		isUpdateValid = true
	updateValidationLoop:
		for i := 1; i < size; i++ { // Loop through each page included in the potential update
			current_page := update[i]
			for j := i - 1; j >= 0; j-- { // Loop through all pages before the current page
				previous_page := update[j]
				for _, val := range ordering_rules[current_page] { // Loop through all ordering rules for the current page
					if previous_page == val {
						isUpdateValid = false
						break updateValidationLoop
					}
				}
			}
		}
		if isUpdateValid {
			num, err := strconv.Atoi(update[size/2])
			if err != nil {
				fmt.Printf("String to int conversion error: %v\n", err)
			}
			middle_sum += num
		} else {
			fixed_update := fixUpdate(update, ordering_rules)
			num, err := strconv.Atoi(fixed_update[size/2])
			if err != nil {
				fmt.Printf("String to int conversion error: %v\n", err)
			}
			fixed_middle_sum += num
		}
	}

	fmt.Printf("[Part 1] The middle sum of all valid updates is: %d\n", middle_sum)
	fmt.Printf("[Part 2] The middle sum of all fixed updates is: %d\n", fixed_middle_sum)
}

// For each page in invalid_update, count how many of the other pages exist in its ordering_rules value (a slice).
// Order the pages by this new count value.
// The new index for each page should be equal to their corresponding count value.
// The solution is expected to be unique.
func fixUpdate(invalid_update []string, ordering_rules map[string][]string) []string {
	fixed_update := []string{}
	new_indexes := map[string]int{}
	for i := 0; i < len(invalid_update); i++ {
		current_page := invalid_update[i]
		current_rules := ordering_rules[current_page]
		new_indexes[current_page] = 0
		for j := 0; j < len(invalid_update); j++ {
			if i == j {
				continue
			}
			other_page := invalid_update[j]
			if isInSlice(other_page, current_rules) {
				new_indexes[current_page] += 1
			}
		}
	}
	for i := 0; i < len(invalid_update); i++ {
		for page, index := range new_indexes {
			if index == i {
				fixed_update = append(fixed_update, page)
			}
		}
	}
	return fixed_update
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
