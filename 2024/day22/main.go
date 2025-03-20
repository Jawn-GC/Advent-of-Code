package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

var step1_memo = make(map[int]int)
var step2_memo = make(map[int]int)
var step3_memo = make(map[int]int)

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Reading initial secret numbers...")
	init_secret_numbers := []int{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		num, _ := strconv.Atoi(line)
		init_secret_numbers = append(init_secret_numbers, num)
	}

	fmt.Println("Calculating new secret numbers...")
	secret_num_sum := 0
	n := 2000
	for _, num := range init_secret_numbers {
		secret_num_sum += getNthSecretNumber(num, n)
	}
	fmt.Printf("[Part 1] Sum of 2000th secret numbers for each buyer: %d\n", secret_num_sum)
}

func getNthSecretNumber(init_num, n int) int {
	new_secret_num := init_num

	for i := 0; i < n; i++ {
		new_secret_num = applyStep1(new_secret_num)
		new_secret_num = applyStep2(new_secret_num)
		new_secret_num = applyStep3(new_secret_num)
	}

	return new_secret_num
}

// Bitwise XOR
func mix(num, val int) int {
	return num ^ val
}

// Modulo 16777216
func prune(num int) int {
	return num % 16777216
}

// Multiply the secret number by 64 and mix the result into the
// secret number. Then, prune and return the secret number.
func applyStep1(num int) int {
	if val, ok := step1_memo[num]; ok {
		return val
	}

	init_num := num
	num = mix(num, num*64)
	num = prune(num)

	step1_memo[init_num] = num
	return num
}

// Divide the secret number by 32 and mix the result (rounded down
// to the nearest integer) into the secret number. Then, prune
// and return the secret number.
func applyStep2(num int) int {
	if val, ok := step2_memo[num]; ok {
		return val
	}

	init_num := num
	num = mix(num, num/32)
	num = prune(num)

	step2_memo[init_num] = num
	return num
}

// Multiply the secret number by 2048 and mix the result into the
// secret number. Then, prune and return the secret number.
func applyStep3(num int) int {
	if val, ok := step3_memo[num]; ok {
		return val
	}

	init_num := num
	num = mix(num, num*2048)
	num = prune(num)

	step3_memo[init_num] = num
	return num
}
