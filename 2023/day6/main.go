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

	fmt.Println("Reading race info...")
	times := []int{}
	record_distances := []int{}
	part2_time_str := ""
	part2_record_distance_str := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ":")
		vals := strings.Fields(temp[1])

		int_slice := []int{}
		str := ""
		for _, val := range vals {
			str += val
			int_val, _ := strconv.Atoi(val)
			int_slice = append(int_slice, int_val)
		}

		if temp[0] == "Time" {
			part2_time_str = str
			times = append(times, int_slice...)
		} else if temp[0] == "Distance" {
			part2_record_distance_str = str
			record_distances = append(record_distances, int_slice...)
		}
	}

	fmt.Println("Calculating winning strategies...")
	num_winning_strats := []int{}
	for i := range len(times) {
		quadratic := getQuadratic(times[i])
		count := 0

		for x := range times[i] {
			distance := quadratic(x)
			if distance > record_distances[i] {
				count += 1
			}
		}

		num_winning_strats = append(num_winning_strats, count)
	}

	product := 1
	for _, val := range num_winning_strats {
		product *= val
	}

	fmt.Printf("[Part 1] %v \n", product)

	part2_time, _ := strconv.Atoi(part2_time_str)
	part2_record_distance, _ := strconv.Atoi(part2_record_distance_str)
	count_part2 := part2_time + 1 // Max number of strategies
	quadratic := getQuadratic(part2_time)
	for i := range part2_time {
		distance := quadratic(i)

		if distance > part2_record_distance {
			break
		} else if count_part2 <= 0 {
			count_part2 = 0
			break
		}

		count_part2 -= 2
	}

	fmt.Printf("[Part 2] %v \n", count_part2)
}

// Returns a function for a parabola with x as the non-zero root
func getQuadratic(x int) func(y int) int {
	return func(y int) int {
		return (x - y) * y
	}
}
