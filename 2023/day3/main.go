package main

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type Point struct {
	Row int
	Col int
}

func (p Point) Add(other Point) Point {
	return Point{
		Row: p.Row + other.Row,
		Col: p.Col + other.Col,
	}
}

var Dir = map[string]Point{
	"UP":        {Row: -1, Col: 0},
	"UPRIGHT":   {Row: -1, Col: 1},
	"RIGHT":     {Row: 0, Col: 1},
	"DOWNRIGHT": {Row: 1, Col: 1},
	"DOWN":      {Row: 1, Col: 0},
	"DOWNLEFT":  {Row: 1, Col: -1},
	"LEFT":      {Row: 0, Col: -1},
	"UPLEFT":    {Row: -1, Col: -1},
}

type Number struct {
	Digits []Point
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Reading engine schematic...")
	schematic := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		new_row := strings.Split(line, "")
		schematic = append(schematic, new_row)
	}

	fmt.Println("Locating part numbers...")
	part_numbers := getPartNumbers(schematic)
	part_number_sum := 0
	for _, part_number := range part_numbers {
		part_number_sum += getPartNumberInt(schematic, part_number)
	}
	fmt.Printf("[Part 1] Part number sum: %d\n", part_number_sum)
}

// Numbers are built as the grid is read from left to right.
func getNumbers(grid [][]string) []Number {
	height := len(grid)
	width := len(grid[0])
	numbers := []Number{}
	validNum := regexp.MustCompile("[0-9]")

	for i := 0; i < height; i++ {
		new_number := Number{}
		readingNumber := false // True if a sequence of numbers has been found
		for j := 0; j < width; j++ {
			ch := grid[i][j]
			isDigit := validNum.MatchString(ch)
			if isDigit {
				new_number.Digits = append(new_number.Digits, Point{Row: i, Col: j})
				readingNumber = true
			} else if readingNumber && !isDigit { // If the end of the current number has been exceeded
				readingNumber = false
				numbers = append(numbers, new_number)
				new_number = Number{}
			}

			if isDigit && j == width-1 { // If the end of the row has been reached
				numbers = append(numbers, new_number)
			}
		}
	}

	return numbers
}

// A number is a part number if it is adjacent to to any non-numeric symbol, besides ".".
func getPartNumbers(grid [][]string) []Number {
	numbers := getNumbers(grid)
	part_numbers := []Number{}
	unwantedSymbol := regexp.MustCompile(`[0-9]|\.`)

	for _, number := range numbers {
		isPartNumber := false
	digitLoop:
		for _, digit := range number.Digits {
			for _, dir := range Dir {
				adj_point := digit.Add(dir)
				if isOOB(grid, adj_point) {
					continue
				}
				if !unwantedSymbol.MatchString(grid[adj_point.Row][adj_point.Col]) {
					isPartNumber = true
					break digitLoop
				}
			}
		}
		if isPartNumber {
			part_numbers = append(part_numbers, number)
		}
	}

	return part_numbers
}

// Checks if a point is Out of Bounds
func isOOB(grid [][]string, p Point) bool {
	height := len(grid)
	width := len(grid[0])
	if p.Row < 0 || p.Row >= height || p.Col < 0 || p.Col >= width {
		return true
	}
	return false
}

func getPartNumberInt(grid [][]string, n Number) int {
	int_str := ""

	for _, digit := range n.Digits {
		int_str += grid[digit.Row][digit.Col]
	}

	num, _ := strconv.Atoi(int_str) // Conversion error ignored
	return num
}
