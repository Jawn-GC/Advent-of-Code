package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Point struct {
	Row int
	Col int
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Locating antennas...\n")
	antennas := make(map[rune][]Point)
	scanner := bufio.NewScanner(file)
	row_count := 0 // Will count up to the height of the grid
	col_count := 0 // Will count up to the width of the grid
	for scanner.Scan() {
		line := scanner.Text()
		line_runes := []rune(line)
		for col_count = 0; col_count < len(line_runes); col_count++ {
			r := line_runes[col_count]
			if unicode.IsLetter(r) || unicode.IsDigit(r) {
				if _, ok := antennas[r]; !ok {
					antennas[r] = []Point{}
				}
				antennas[r] = append(antennas[r], Point{Row: row_count, Col: col_count})
			}
		}
		row_count++
	}

	fmt.Printf("Locating antinodes..\n")
	antinodes1 := findAntinodes1(antennas, row_count, col_count)
	antinodes2 := findAntinodes2(antennas, row_count, col_count)
	fmt.Printf("[Part 1] There are %d unique antinodes.\n", len(antinodes1))
	fmt.Printf("[Part 2] There are %d unique antinodes.\n", len(antinodes2))
}

// A point is an antinode of a pair of antennas of the same type if it is in line
// with both of them, but only when one of the antennas is twice as far away as
// the other. The antinode must be within the bounds of the grid.
func findAntinodes1(antennas map[rune][]Point, height int, width int) []Point {
	antinodes := []Point{}

	for _, points := range antennas {
		n := len(points)

		// Antinodes are created if there are at least two antennas of a type.
		// Skip any group of 1 or fewer antennas.
		if n < 2 {
			continue
		}

		// Iterate over all possible pairs of antennas.
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				dx := points[j].Col - points[i].Col
				dy := points[j].Row - points[i].Row

				new_row := points[j].Row + dy
				new_col := points[j].Col + dx
				new_point := Point{Row: new_row, Col: new_col}
				if new_row >= 0 && new_row < height && new_col >= 0 && new_col < width {
					if !isInSlice(new_point, antinodes) {
						antinodes = append(antinodes, new_point)
					}
				}

				new_row = points[j].Row - 2*dy
				new_col = points[j].Col - 2*dx
				new_point = Point{Row: new_row, Col: new_col}
				if new_row >= 0 && new_row < height && new_col >= 0 && new_col < width {
					if !isInSlice(new_point, antinodes) {
						antinodes = append(antinodes, new_point)
					}
				}
			}
		}
	}

	return antinodes
}

// A point is an antinode of a pair of antennas of the same type if it is in line
// with both of them. The antinode must be within the bounds of the grid.
func findAntinodes2(antennas map[rune][]Point, height int, width int) []Point {
	antinodes := []Point{}

	for _, points := range antennas {
		n := len(points)

		// Antinodes are created if there are at least two antennas of a type.
		// Skip any group of 1 or fewer antennas.
		if n < 2 {
			continue
		}

		// Iterate over all possible pairs of antennas.
		for i := 0; i < n-1; i++ {
			for j := i + 1; j < n; j++ {
				dx := points[j].Col - points[i].Col
				dy := points[j].Row - points[i].Row
				if !isInSlice(points[i], antinodes) {
					antinodes = append(antinodes, points[i])
				}
				if !isInSlice(points[j], antinodes) {
					antinodes = append(antinodes, points[j])
				}

				// In the direction from Point i to Point j.
				new_row := points[j].Row + dy
				new_col := points[j].Col + dx
				for {
					if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
						break
					}
					new_point := Point{Row: new_row, Col: new_col}
					if !isInSlice(new_point, antinodes) {
						antinodes = append(antinodes, new_point)
					}

					new_row += dy
					new_col += dx
				}

				// In the direction from Point j to Point i.
				new_row = points[i].Row - dy
				new_col = points[i].Col - dx
				for {
					if new_row < 0 || new_row >= height || new_col < 0 || new_col >= width {
						break
					}
					new_point := Point{Row: new_row, Col: new_col}
					if !isInSlice(new_point, antinodes) {
						antinodes = append(antinodes, new_point)
					}

					new_row -= dy
					new_col -= dx
				}
			}
		}
	}

	return antinodes
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
