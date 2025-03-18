package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Used to traverse a 2D slice representing a word search grid
// Dy is change in row index
// Dx is change in column index
type Delta struct {
	Dy int
	Dx int
}

var Directions = map[string]Delta{
	"UP":         {-1, 0},
	"UP_RIGHT":   {-1, 1},
	"RIGHT":      {0, 1},
	"DOWN_RIGHT": {1, 1},
	"DOWN":       {1, 0},
	"DOWN_LEFT":  {1, -1},
	"LEFT":       {0, -1},
	"UP_LEFT":    {-1, -1},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s\n", filename)
		return
	}
	defer file.Close()

	fmt.Printf("Reading word search...\n")
	word_search := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []string{}
		row = append(row, strings.Split(line, "")...)
		word_search = append(word_search, row)
	}

	// The word search grid is assumed to be nonempty and rectangular
	// The word is assumed to be at least 2 characters
	word_count := 0
	pattern_count := 0
	height := len(word_search)
	width := len(word_search[0])
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// Part 1
			word_count += numXMASWord(word_search, i, j)
			// Part 2
			if hasXMASPattern(word_search, i, j) {
				pattern_count += 1
			}
		}
	}

	fmt.Printf("[Part 1] The word 'XMAS' was found %d times.\n", word_count)
	fmt.Printf("[Part 2] The X-MAS pattern was found %d times.\n", pattern_count)
}

func numXMASWord(word_search [][]string, row int, col int) int {
	height := len(word_search)
	width := len(word_search[0])
	word := "XMAS"
	word_chars := strings.Split(word, "")
	word_count := 0
	if word_search[row][col] != word_chars[0] {
		return 0
	}
	// If the intial letter is found, then search all eight directions for the remaining sequence of letters
	for _, delta := range Directions {
		new_row := row + delta.Dy
		new_col := col + delta.Dx
		for k := 1; k < len(word_chars); k++ {
			// Check if search goes out of bounds
			if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
				break
			}
			if word_search[new_row][new_col] != word_chars[k] {
				break
			}
			// Increment count if the entire word has been found
			if k == len(word_chars)-1 {
				word_count += 1
			}
			new_row += delta.Dy
			new_col += delta.Dx
		}
	}
	return word_count
}

func hasXMASPattern(word_search [][]string, row int, col int) bool {
	if word_search[row][col] != "A" {
		return false
	}

	height := len(word_search)
	width := len(word_search[0])

	var new_row, new_col int

	// Set upper left corner
	new_row, new_col = row+Directions["UP_LEFT"].Dy, col+Directions["UP_LEFT"].Dx
	if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
		return false
	}
	upper_left := word_search[new_row][new_col]

	// Set upper right corner
	new_row, new_col = row+Directions["UP_RIGHT"].Dy, col+Directions["UP_RIGHT"].Dx
	if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
		return false
	}
	upper_right := word_search[new_row][new_col]

	// Set lower left corner
	new_row, new_col = row+Directions["DOWN_LEFT"].Dy, col+Directions["DOWN_LEFT"].Dx
	if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
		return false
	}
	lower_left := word_search[new_row][new_col]

	// Set lower right corner
	new_row, new_col = row+Directions["DOWN_RIGHT"].Dy, col+Directions["DOWN_RIGHT"].Dx
	if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
		return false
	}
	lower_right := word_search[new_row][new_col]

	// Check if opposite corners are a pair of "M" and "S"
	if !(upper_left == "M" && lower_right == "S") && !(upper_left == "S" && lower_right == "M") {
		return false
	}
	if !(upper_right == "M" && lower_left == "S") && !(upper_right == "S" && lower_left == "M") {
		return false
	}

	return true
}
