package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Point struct {
	Row int
	Col int
}

func (p Point) Plus(other Point) Point {
	return Point{
		Row: p.Row + other.Row,
		Col: p.Col + other.Col,
	}
}

func (p Point) Minus(other Point) Point {
	return Point{
		Row: p.Row - other.Row,
		Col: p.Col - other.Col,
	}
}

var Deltas = map[string]Point{
	"^": {Row: -1, Col: 0},
	">": {Row: 0, Col: 1},
	"v": {Row: 1, Col: 0},
	"<": {Row: 0, Col: -1},
}

// Numerical Keypad Reference
// +---+---+---+
// | 7 | 8 | 9 |
// +---+---+---+
// | 4 | 5 | 6 |
// +---+---+---+
// | 1 | 2 | 3 |
// +---+---+---+
//     | 0 | A |
//     +---+---+

var num_ref = map[string]Point{
	"0": {Row: 3, Col: 1},
	"1": {Row: 2, Col: 0},
	"2": {Row: 2, Col: 1},
	"3": {Row: 2, Col: 2},
	"4": {Row: 1, Col: 0},
	"5": {Row: 1, Col: 1},
	"6": {Row: 1, Col: 2},
	"7": {Row: 0, Col: 0},
	"8": {Row: 0, Col: 1},
	"9": {Row: 0, Col: 2},
	"A": {Row: 3, Col: 2},
}

// Directional Keypad Reference
//     +---+---+
//     | ^ | A |
// +---+---+---+
// | < | v | > |
// +---+---+---+

var dir_ref = map[string]Point{
	"^": {Row: 0, Col: 1},
	">": {Row: 1, Col: 2},
	"v": {Row: 1, Col: 1},
	"<": {Row: 1, Col: 0},
	"A": {Row: 0, Col: 2},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading codes...\n")
	codes := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		codes = append(codes, line)
	}

	// Robot 1 is at the numpad
	// Robots 2 & 3 are at direction pads
	// The user is at a direction pad
	fmt.Printf("Writing instructions...\n")
	complexity_sum := 0
	for _, code := range codes {
		robot1 := num_ref["A"]
		robot2 := dir_ref["A"]
		robot3 := dir_ref["A"]

		r2_code := numToDir(&robot1, code)
		r3_code := dirToDir(&robot2, r2_code)
		user_code := dirToDir(&robot3, r3_code)

		num, _ := strconv.Atoi(strings.TrimRight(code, "A"))
		complexity_sum += len(user_code) * num
	}

	fmt.Printf("Complexity sum: %d\n", complexity_sum)
}

func numToDir(robot *Point, code string) string {
	instructions := ""
	chars := strings.Split(code, "")

	for _, ch := range chars {
		delta := num_ref[ch].Minus(*robot) // Desired position minus current position

		// Avoid crossing empty space in numpad
		if robot.Row == 3 && num_ref[ch].Col == 0 {
			if delta.Row < 0 {
				for i := 0; i < delta.Row*(-1); i++ {
					instructions += "^"
					*robot = robot.Plus(Deltas["^"])
				}
			}
			if delta.Col < 0 {
				for i := 0; i < delta.Col*(-1); i++ {
					instructions += "<"
					*robot = robot.Plus(Deltas["<"])
				}
			}
			instructions += "A"
			continue
		}

		if robot.Col == 0 && num_ref[ch].Row == 3 {
			if delta.Col > 0 {
				for i := 0; i < delta.Col; i++ {
					instructions += ">"
					*robot = robot.Plus(Deltas[">"])
				}
			}
			if delta.Row > 0 {
				for i := 0; i < delta.Row; i++ {
					instructions += "v"
					*robot = robot.Plus(Deltas["v"])
				}
			}
			instructions += "A"
			continue
		}

		if delta.Col < 0 {
			for i := 0; i < delta.Col*(-1); i++ {
				instructions += "<"
				*robot = robot.Plus(Deltas["<"])
			}
		}

		if delta.Row > 0 {
			for i := 0; i < delta.Row; i++ {
				instructions += "v"
				*robot = robot.Plus(Deltas["v"])
			}
		}

		if delta.Row < 0 {
			for i := 0; i < delta.Row*(-1); i++ {
				instructions += "^"
				*robot = robot.Plus(Deltas["^"])
			}
		}

		if delta.Col > 0 {
			for i := 0; i < delta.Col; i++ {
				instructions += ">"
				*robot = robot.Plus(Deltas[">"])
			}
		}

		instructions += "A"
	}

	return instructions
}

func dirToDir(robot *Point, code string) string {
	instructions := ""
	chars := strings.Split(code, "")

	for _, ch := range chars {
		delta := dir_ref[ch].Minus(*robot) // Desired position minus current position

		// Avoid crossing empty space in direction pad
		if robot.Col == 0 && dir_ref[ch].Row == 0 {
			for i := 0; i < delta.Col; i++ {
				instructions += ">"
				*robot = robot.Plus(Deltas[">"])
			}
			for i := 0; i < delta.Row*(-1); i++ {
				instructions += "^"
				*robot = robot.Plus(Deltas["^"])
			}
			instructions += "A"
			continue
		}

		if robot.Row == 0 && dir_ref[ch].Col == 0 {
			for i := 0; i < delta.Row; i++ {
				instructions += "v"
				*robot = robot.Plus(Deltas["v"])
			}
			for i := 0; i < delta.Col*(-1); i++ {
				instructions += "<"
				*robot = robot.Plus(Deltas["<"])
			}
			instructions += "A"
			continue
		}

		if delta.Col < 0 {
			for i := 0; i < delta.Col*(-1); i++ {
				instructions += "<"
				*robot = robot.Plus(Deltas["<"])
			}
		}

		if delta.Row > 0 {
			for i := 0; i < delta.Row; i++ {
				instructions += "v"
				*robot = robot.Plus(Deltas["v"])
			}
		}

		if delta.Row < 0 {
			for i := 0; i < delta.Row*(-1); i++ {
				instructions += "^"
				*robot = robot.Plus(Deltas["^"])
			}
		}

		if delta.Col > 0 {
			for i := 0; i < delta.Col; i++ {
				instructions += ">"
				*robot = robot.Plus(Deltas[">"])
			}
		}

		instructions += "A"
	}

	return instructions
}
