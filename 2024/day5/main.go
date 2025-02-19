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

	ordering_rules := map[string][]string{}
	orderingRulesRead := false
	updates := [][]string{}
	fmt.Print("Parsing ordering rules...\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			orderingRulesRead = true
			fmt.Print("Parsing updates...\n")
			continue
		}
		if orderingRulesRead == false {
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
	valid_updates := [][]string{}
	var isUpdateValid bool
	middle_sum := 0
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
			valid_updates = append(valid_updates, update)
			num, err := strconv.Atoi(update[size/2])
			if err != nil {
				fmt.Printf("String to int conversion error: %v\n", err)
			}
			middle_sum += num
		}
	}

	fmt.Printf("The middle sum of all valid updates is: %d\n", middle_sum)
}
