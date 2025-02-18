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
	count := 0
	height := len(word_search)
	width := len(word_search[0])
	word := "XMAS"
	word_chars := strings.Split(word, "")
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if word_search[i][j] != word_chars[0] {
				continue
			}
			// If the intial letter is found, then search all eight directions for the remaining sequence of letters
			for _, delta := range Directions {
				new_row := i + delta.Dy
				new_col := j + delta.Dx
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
						count += 1
					}
					new_row += delta.Dy
					new_col += delta.Dx
				}
			}
		}
	}

	fmt.Printf("The word '%s' was found %d times in the word search\n", word, count)
}
