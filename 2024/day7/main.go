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

	// By storing the equations in a map I am assuming that a "result"
	// value does not appear more than once in "input.txt"
	equations := make(map[string][]string)
	scanner := bufio.NewScanner(file)
	fmt.Printf("Parsing equations...\n")
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ":")
		result := strings.TrimSpace(temp[0])
		terms := strings.Fields(temp[1])
		equations[result] = terms
	}

	fmt.Printf("Calibrating...\n")
	calibration_result1, err := getCalibrationResult(equations, 2)
	if err != nil {
		fmt.Printf("Could not complete calibration: %v", err)
	}
	calibration_result2, err := getCalibrationResult(equations, 3)
	if err != nil {
		fmt.Printf("Could not complete calibration: %v", err)
	}
	fmt.Printf("[Part 1] Calibration Result: %d\n", calibration_result1)
	fmt.Printf("[Part 2] Calibration Result: %d\n", calibration_result2)
}

func getCalibrationResult(equations map[string][]string, base int) (int, error) {
	total := 0

	for key, value := range equations {
		result, err := strconv.Atoi(key)
		if err != nil {
			return 0, err
		}

		terms := []int{}
		for _, term_str := range value {
			term, err := strconv.Atoi(term_str)
			if err != nil {
				return 0, err
			}
			terms = append(terms, term)
		}

		is_valid_permutation, err := validPermutationFound(result, terms, base)
		if err != nil {
			return 0, err
		}
		if is_valid_permutation {
			total += result
		}
	}

	return total, nil
}

// Determine if there exists some order of operations between the "terms"
// that gives "result". Operations are resolved from left to right.
// In the operations string, 0 is +, 1 is *, and 2 is concantenation.
func validPermutationFound(result int, terms []int, base int) (bool, error) {
	n := len(terms) - 1
	num_permutations := intPow(base, n)
	operations := 0
	for i := 0; i < num_permutations; i++ {
		total := terms[0]
		base_str := strconv.FormatInt(int64(operations), base)
		operations_str := fmt.Sprintf("%0*s", n, base_str) // Pad the number with leading zeros
		operations_runes := []rune(operations_str)
		for j := 0; j < n; j++ {
			if operations_runes[j] == '0' {
				total += terms[j+1]
			} else if operations_runes[j] == '1' {
				total *= terms[j+1]
			} else if operations_runes[j] == '2' {
				temp_str := strconv.FormatInt(int64(total), 10) + strconv.FormatInt(int64(terms[j+1]), 10)
				temp_int, err := strconv.Atoi(temp_str)
				if err != nil {
					return false, err
				}
				total = temp_int
			}
		}
		if result == total {
			return true, nil
		}
		operations++
	}

	return false, nil
}

// Helper function for calculating powers
func intPow(base int, exp int) int {
	result := 1
	for i := 0; i < exp; i++ {
		result *= base
	}
	return result
}
