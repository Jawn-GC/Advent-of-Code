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

	pattern := `mul\([0-9]{1,3}\,[0-9]{1,3}\)`
	re := regexp.MustCompile(pattern)
	var matches []string

	fmt.Printf("Scanning corrupted memory...\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		matches = append(matches, re.FindAllString(line, -1)...)
	}

	product_sum := 0
	fmt.Printf("Parsing integer pairs...\n")
	for _, match := range matches {
		trimmed_match := strings.TrimLeft(match, "mul(")
		trimmed_match = strings.TrimRight(trimmed_match, ")")
		parsed_ints := strings.Split(trimmed_match, ",")

		num1, err := strconv.Atoi(parsed_ints[0])
		if err != nil {
			fmt.Printf("Could not convert %s to integer: %v", parsed_ints[0], err)
		}
		num2, err := strconv.Atoi(parsed_ints[1])
		if err != nil {
			fmt.Printf("Could not convert %s to integer: %v", parsed_ints[1], err)
		}

		product_sum += num1 * num2
	}

	fmt.Printf("The sum of products is: %d\n", product_sum)
}
