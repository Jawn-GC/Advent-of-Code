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

	fmt.Printf("Parsing initial stone values...\n")
	stones := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		stone_strings := strings.Fields(line)
		for _, str := range stone_strings {
			num, err := strconv.Atoi(str)
			if err != nil {
				fmt.Printf("Error converting string: %v\n", err)
				return
			}
			stones = append(stones, num)
		}
	}

	fmt.Printf("Blinking...\n")
	num_blinks := 75
	stone_freq := map[int]int{}

	for _, stone := range stones {
		if _, ok := stone_freq[stone]; !ok {
			stone_freq[stone] = 0
		}
		stone_freq[stone] += 1
	}

	for i := 0; i < num_blinks; i++ {
		stone_freq = blink(stone_freq)
	}

	num_stones := 0
	for _, freq := range stone_freq {
		num_stones += freq
	}
	fmt.Printf("After %d blinks there are %d stones\n", num_blinks, num_stones)
}

func blink(stone_freq map[int]int) map[int]int {
	new_stone_freq := map[int]int{}

	for stone, freq := range stone_freq {
		if stone == 0 {
			new_stone_freq[1] += freq
		} else if numDigits(stone)%2 == 0 {
			stone_left, stone_right := splitNumber(stone)
			new_stone_freq[stone_left] += freq
			new_stone_freq[stone_right] += freq
		} else {
			new_stone_freq[stone*2024] += freq
		}
	}

	return new_stone_freq
}

func splitNumber(num int) (int, int) {
	digits := numDigits(num)
	magnitude := 1
	for i := 0; i < digits/2; i++ {
		magnitude *= 10
	}
	left := num / magnitude
	right := num % magnitude
	return left, right
}

func numDigits(num int) int {
	temp := num
	digits := 0
	for temp > 0 {
		digits++
		temp /= 10
	}
	return digits
}
