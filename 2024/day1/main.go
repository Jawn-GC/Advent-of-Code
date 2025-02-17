package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

	var col1, col2 []int
	counts := make(map[int]int)
	scanner := bufio.NewScanner(file)

	fmt.Printf("Parsing integers...\n")
	for scanner.Scan() {
		line := scanner.Text()
		num_strs := strings.Fields(line)

		num1, err := strconv.Atoi(num_strs[0])
		if err != nil {
			fmt.Printf("Error converting %s to an integer: %v\n", num_strs[0], err)
		}
		col1 = append(col1, num1)

		num2, err := strconv.Atoi(num_strs[1])
		if err != nil {
			fmt.Printf("Error converting %s to an integer: %v\n", num_strs[1], err)
		}
		col2 = append(col2, num2)

		if _, ok := counts[num2]; !ok {
			counts[num2] = 0
		}
		counts[num2] += 1
	}

	if len(col1) != len(col2) {
		fmt.Printf("Error: unequal list lengths\n")
	}

	sort.Ints(col1)
	sort.Ints(col2)

	total_difference := 0
	for i := range col1 {
		difference := col1[i] - col2[i]
		if difference < 0 {
			difference = -difference
		}
		total_difference += difference
	}

	similarity := 0
	for _, val := range col1 {
		count, ok := counts[val]
		if ok {
			similarity += val * count
		}
	}

	fmt.Printf("The total difference between the lists is: %d\n", total_difference)
	fmt.Printf("The similarity between the lists is: %d\n", similarity)
}
