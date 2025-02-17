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

	var reports [][]int
	num_safe_reports := 0

	fmt.Printf("Parsing reports...\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		report_strings := strings.Fields(line)
		report_ints, err := convertReportStringsToInts(report_strings)
		if err != nil {
			fmt.Printf("Error converting strings to ints: %v\n", err)
		}

		reports = append(reports, report_ints)
	}

	for _, report := range reports {
		if isReportSafe(report, true) {
			num_safe_reports += 1
		}
	}

	fmt.Printf("The number of safe reports is: %d\n", num_safe_reports)
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

func isReportSafe(report []int, dampen ...bool) bool {
	isSafe := true
	isDecreasing, isIncreasing := true, true
	isDamp := false
	var previous_num int

	if len(report) == 0 {
		return true
	}

	if len(dampen) > 0 {
		isDamp = dampen[0]
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

	if isDamp && !isSafe {
		for i := 0; i < len(report); i++ {
			sliced_report := []int{}
			for j := 0; j < len(report); j++ {
				if i != j {
					sliced_report = append(sliced_report, report[j])
				}
			}
			isSafe = isReportSafe(sliced_report)
			if isSafe {
				break
			}
		}
	}

	return isSafe
}
