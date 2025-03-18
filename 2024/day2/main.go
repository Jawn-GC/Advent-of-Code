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
		fmt.Printf("Could not open %s\n", filename)
		return
	}

	defer file.Close()

	fmt.Printf("Parsing reports...\n")
	var reports [][]int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		report_strings := strings.Fields(line)
		report_ints, err := convertReportStringsToInts(report_strings)
		if err != nil {
			fmt.Printf("Error converting strings to ints: %v\n", err)
			return
		}

		reports = append(reports, report_ints)
	}

	// Part 1
	fmt.Println("[Part 1] Detecting safe reports...")
	num_safe_reports1 := 0
	for _, report := range reports {
		if isReportSafe(report, false) {
			num_safe_reports1 += 1
		}
	}
	// Part 2
	fmt.Println("[Part 2] Detecting safe reports with dampening...")
	num_safe_reports2 := 0
	for _, report := range reports {
		if isReportSafe(report, true) {
			num_safe_reports2 += 1
		}
	}

	fmt.Printf("[Part 1] The number of safe reports is: %d\n", num_safe_reports1)
	fmt.Printf("[Part 2] The number of safe reports is: %d\n", num_safe_reports2)
}

func convertReportStringsToInts(report_strings []string) ([]int, error) {
	var report_ints []int
	for _, str := range report_strings {
		num, err := strconv.Atoi(str)
		if err != nil {
			return []int{}, err
		}
		report_ints = append(report_ints, num)
	}
	return report_ints, nil
}

// A report is labeled as safe if the ints are strictly increasing/decreasing
// and differences between consecutive ints are no more than 3.
// If damping is applied, a report is safe if a single int can be removed to
// make the report safe.
func isReportSafe(report []int, dampen bool) bool {
	isSafe := true
	isDecreasing, isIncreasing := true, true
	var previous_num int

	if len(report) == 0 {
		return true
	}

	for i, num := range report {
		if i == 0 {
			previous_num = num
			continue
		}

		difference := num - previous_num
		if difference == 0 {
			isSafe = false
			break
		} else if difference > 0 {
			isDecreasing = false
		} else if difference < 0 {
			difference = -difference
			isIncreasing = false
		}

		if !isDecreasing && !isIncreasing {
			isSafe = false
			break
		}

		if difference > 3 {
			isSafe = false
			break
		}

		previous_num = num
	}

	if dampen && !isSafe {
		for i := 0; i < len(report); i++ {
			sliced_report := []int{}
			for j := 0; j < len(report); j++ {
				if i != j {
					sliced_report = append(sliced_report, report[j])
				}
			}
			isSafe = isReportSafe(sliced_report, false)
			if isSafe {
				break
			}
		}
	}

	return isSafe
}
