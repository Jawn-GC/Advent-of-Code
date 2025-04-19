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

	fmt.Println("Reading histories...")
	histories := [][]int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		int_row := []int{}
		str_row := strings.Fields(line)
		for _, str := range str_row {
			val, _ := strconv.Atoi(str)
			int_row = append(int_row, val)
		}

		histories = append(histories, int_row)
	}

	fmt.Println("Extrapolating values...")
	predicted_value_sum := 0
	predicted_value_sum2 := 0
	for _, history := range histories {
		predicted_value_sum += predictNextValue(history)
	}
	for _, history := range histories {
		predicted_value_sum2 += predictPreviousValue(history)
	}
	fmt.Printf("[Part 1] The sum of next predicted values is %d\n", predicted_value_sum)
	fmt.Printf("[Part 1] The sum of previous predicted values is %d\n", predicted_value_sum2)
}

// It is assumed that a sequence of only zeroes will eventually be generated.
// It is assmued that the starting history is not all zeroes.
func predictNextValue(history []int) int {
	sequences := getHistorySequences(history)

	extrapolated_val := 0
	for i := len(sequences) - 2; i >= 0; i-- {
		row := sequences[i]
		n := len(row)
		extrapolated_val += row[n-1]
	}

	return extrapolated_val
}

func predictPreviousValue(history []int) int {
	sequences := getHistorySequences(history)

	extrapolated_val := 0
	for i := len(sequences) - 2; i >= 0; i-- {
		row := sequences[i]
		extrapolated_val = row[0] - extrapolated_val
	}

	return extrapolated_val
}

func isAllZeroes(sequence []int) bool {
	for _, val := range sequence {
		if val != 0 {
			return false
		}
	}

	return true
}

func getHistorySequences(history []int) [][]int {
	sequences := [][]int{}
	sequences = append(sequences, history)
	current_sequence := history

	for {
		new_sequence := []int{}

		for i := 0; i < len(current_sequence)-1; i++ {
			new_val := current_sequence[i+1] - current_sequence[i]
			new_sequence = append(new_sequence, new_val)
		}
		sequences = append(sequences, new_sequence)
		current_sequence = new_sequence

		if isAllZeroes(new_sequence) {
			break
		}
	}

	return sequences
}
