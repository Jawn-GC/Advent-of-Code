package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s\n", filename)
		return
	}
	defer file.Close()

	pattern1 := `mul\([0-9]{1,3}\,[0-9]{1,3}\)`
	pattern2 := `mul\([0-9]{1,3}\,[0-9]{1,3}\)|do\(\)|don\'t\(\)`
	re1 := regexp.MustCompile(pattern1)
	re2 := regexp.MustCompile(pattern2)

	fmt.Printf("Scanning corrupted memory...\n")
	var matches1 []string
	var matches2 []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches1 = append(matches1, re1.FindAllString(line, -1)...)
		matches2 = append(matches2, re2.FindAllString(line, -1)...)
	}

	fmt.Printf("Executing instructions...\n")
	product_sum1 := getProductSum(matches1)
	product_sum2 := getProductSum(matches2)

	fmt.Printf("[Part 1] The sum of products is: %d\n", product_sum1)
	fmt.Printf("[Part 2] The sum of products is: %d\n", product_sum2)
}

func getProductSum(matches []string) int {
	product_sum := 0
	instructionsEnabled := true
	for _, match := range matches {
		if strings.Contains(match, "don't") {
			instructionsEnabled = false
			continue
		} else if strings.Contains(match, "do") {
			instructionsEnabled = true
			continue
		}

		if !instructionsEnabled {
			continue
		}

		trimmed_match := strings.TrimLeft(match, "mul(")
		trimmed_match = strings.TrimRight(trimmed_match, ")")
		parsed_ints := strings.Split(trimmed_match, ",")

		num1, err := strconv.Atoi(parsed_ints[0])
		if err != nil {
			fmt.Printf("Could not convert %s to integer: %v\n", parsed_ints[0], err)
			return -1
		}
		num2, err := strconv.Atoi(parsed_ints[1])
		if err != nil {
			fmt.Printf("Could not convert %s to integer: %v\n", parsed_ints[1], err)
			return -1
		}

		product_sum += num1 * num2
	}
	return product_sum
}
