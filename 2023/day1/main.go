package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"unicode"
)

var digit_ref = map[string]string{
	"one":   "1",
	"two":   "2",
	"three": "3",
	"four":  "4",
	"five":  "5",
	"six":   "6",
	"seven": "7",
	"eight": "8",
	"nine":  "9",
}

type Overlap struct {
	Name   string
	Symbol string
	Offset int
}

var possible_overlaps = []Overlap{
	{Name: "oneight", Symbol: "e", Offset: 2},
	{Name: "twone", Symbol: "o", Offset: 2},
	{Name: "threeight", Symbol: "e", Offset: 4},
	{Name: "fiveight", Symbol: "e", Offset: 3},
	{Name: "sevenine", Symbol: "n", Offset: 4},
	{Name: "eightwo", Symbol: "t", Offset: 4},
	{Name: "eighthree", Symbol: "t", Offset: 4},
	{Name: "nineight", Symbol: "e", Offset: 3},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
	}
	defer file.Close()

	fmt.Println("[Part 1] Reading calibration document...")
	document := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		document = append(document, line)
	}

	// Part 1
	fmt.Println("[Part 1] Parsing calibration values...")
	calibration_sum1 := 0
	for _, line := range document {
		calibration_sum1 += getCalibrationValue(line)
	}

	// Part 2
	fmt.Println("[Part 2] Detecting overlaps...")
	modified_document := splitOverlaps(document)
	fmt.Println("[Part 2] Reformatting document...")
	reformatted_document := reformatDocument(modified_document)
	fmt.Println("[Part 2] Parsing calibration values...")
	calibration_sum2 := 0
	for _, line := range reformatted_document {
		calibration_sum2 += getCalibrationValue(line)
	}

	fmt.Printf("[Part 1] Sum of calibration values: %d\n", calibration_sum1)
	fmt.Printf("[Part 2] Sum of calibration values: %d\n", calibration_sum2)
}

// The calibration value of a line is the first numeric value
// concantenated with the last numeric value.
func getCalibrationValue(line string) int {
	digits := []rune{}
	calibration_value := ""
	for _, ch := range line {
		if unicode.IsDigit(ch) {
			digits = append(digits, ch)
		}
	}

	// The digits slice is assumed to be non-empty.
	calibration_value += string(digits[0])
	calibration_value += string(digits[len(digits)-1])

	// Conversion errors are ignored.
	calibration_value_num, _ := strconv.Atoi(calibration_value)
	return calibration_value_num
}

// Some numeric values of the original document have been replaced
// by strings of their names (ex. 1 -> one). This function reverts this
// change (ex. one -> 1). Lines are assumed to be lowercase.
func reformatDocument(doc []string) []string {
	reformatted_doc := []string{}

	for _, line := range doc {
		chars := strings.Split(line, "")
		temp_str := ""
		reformatted_line := ""

		// Build temp strings from left to right
		for _, ch := range chars {
			temp_str += ch
			// Check to see if the accumulated string matches any number's
			// name (up to the length of temp_str).
			partial_match := false
			complete_match := false
			for key, _ := range digit_ref {
				if temp_str == key {
					partial_match = true
					complete_match = true
					break
				}
				if len(temp_str) > len(key) {
					continue
				}
				for i := 0; i < len(temp_str); i++ {
					if temp_str[i] == key[i] && i == len(temp_str)-1 {
						partial_match = true
					} else if temp_str[i] != key[i] {
						break
					}
				}
			}

			// If no substring match, add temp str to the reformatted
			// line. If a complete match was found, append the digit
			// instead of the temp_str.
			if complete_match {
				reformatted_line += digit_ref[temp_str]
				temp_str = ""
				partial_match = false
				complete_match = false
			} else if !partial_match {
				reformatted_line += temp_str[:len(temp_str)-1]
				temp_str = string(temp_str[len(temp_str)-1])
			}
		}

		// Append any remaining chars
		if len(temp_str) != 0 {
			reformatted_line += temp_str
		}
		reformatted_doc = append(reformatted_doc, reformatted_line)
	}

	return reformatted_doc
}

// Some lines in the document may have two or more number names
// overlapping by one letter (ex. twone). This function will
// duplicate each overlapped letter.
func splitOverlaps(doc []string) []string {
	modified_doc := []string{}
	for _, line := range doc {
		for _, overlap := range possible_overlaps {
			for {
				match_found := false

				index := strings.Index(line, overlap.Name)
				if index != -1 {
					match_found = true
					dupe_index := index + overlap.Offset
					line = line[:dupe_index] + overlap.Symbol + line[dupe_index:]
				}

				if !match_found {
					break
				}
			}
		}
		modified_doc = append(modified_doc, line)
	}
	return modified_doc
}
