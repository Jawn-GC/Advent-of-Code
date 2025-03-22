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

type Digit struct {
	Pos    Point
	Symbol string  // The string stored at Point
	Number *Number // Reference to the Number this Digit is part of
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
	Digits []Digit
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
	var part_numbers []Number
	getPartNumbers(schematic, &part_numbers)
	part_number_sum := 0
	for _, part_number := range part_numbers {
		part_number_sum += getPartNumberInt(part_number)
	}

	fmt.Println("Locating gears...")
	gears := getGears(schematic)

	fmt.Println("Calculating gear ratios...")
	gear_ratio_sum := 0
	for _, gear := range gears {
		gear_ratio_sum += getGearRatio(gear, &part_numbers)
	}

	fmt.Printf("[Part 1] Part number sum: %d\n", part_number_sum)
	fmt.Printf("[Part 2] Gear ratio sum: %d\n", gear_ratio_sum)
}

// Numbers are built as the grid is read from left to right.
func getNumbers(grid [][]string, numbers *[]Number) {
	height := len(grid)
	width := len(grid[0])
	validNum := regexp.MustCompile("[0-9]")

	for i := 0; i < height; i++ {
		new_number := Number{}
		readingNumber := false // True if a sequence of numbers has been found
		for j := 0; j < width; j++ {
			ch := grid[i][j]
			isDigit := validNum.MatchString(ch)
			if isDigit {
				new_number.Digits = append(new_number.Digits, Digit{Pos: Point{Row: i, Col: j}, Symbol: ch})
				readingNumber = true
			} else if readingNumber && !isDigit { // If the end of the current number has been exceeded
				readingNumber = false
				*numbers = append(*numbers, new_number)
				new_number = Number{}
			}

			if isDigit && j == width-1 { // If the end of the row has been reached
				*numbers = append(*numbers, new_number)
			}
		}
	}

	for i := 0; i < len(*numbers); i++ {
		for j := 0; j < len((*numbers)[i].Digits); j++ {
			(*numbers)[i].Digits[j].Number = &(*numbers)[i]
		}
	}
}

// A number is a part number if it is adjacent to to any non-numeric symbol, besides ".".
func getPartNumbers(grid [][]string, part_numbers *[]Number) {
	var numbers []Number
	getNumbers(grid, &numbers)
	unwantedSymbol := regexp.MustCompile(`[0-9]|\.`)

	for _, number := range numbers {
		isPartNumber := false
	digitLoop:
		for _, digit := range number.Digits {
			for _, dir := range Dir {
				adj_point := digit.Pos.Add(dir)
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
			*part_numbers = append(*part_numbers, number)
		}
	}
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

func getPartNumberInt(n Number) int {
	int_str := ""

	for _, digit := range n.Digits {
		int_str += digit.Symbol
	}

	num, _ := strconv.Atoi(int_str) // Conversion error ignored
	return num
}

func getGears(grid [][]string) []Point {
	height := len(grid)
	width := len(grid[0])
	gears := []Point{}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			if grid[i][j] == "*" {
				gears = append(gears, Point{Row: i, Col: j})
			}
		}
	}

	return gears
}

// If a gear is adjacent to exactly two part numbers, the ratio is
// the product of the numbers.
func getGearRatio(gear Point, part_numbers *[]Number) int {
	adj_numbers := []*Number{}
	digits := []Digit{}

	for _, number := range *part_numbers {
		digits = append(digits, number.Digits...)
	}

	for _, delta := range Dir {
		adj_point := gear.Add(delta)
		for _, digit := range digits {
			if adj_point == digit.Pos {
				if !isInSlice(digit.Number, adj_numbers) {
					adj_numbers = append(adj_numbers, digit.Number)
				}
			}
		}
	}
	if len(adj_numbers) == 2 {
		num1 := getPartNumberInt(*adj_numbers[0])
		num2 := getPartNumberInt(*adj_numbers[1])
		return num1 * num2
	}

	return 0 // Not a valid gear
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
