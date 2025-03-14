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
	num_dir_robots := 2
	dir_memo := make(map[string]string) // Key: "start_symbol"+"end_symbol"
	for _, code := range codes {
		main_robot := num_ref["A"]
		dir_robots := []Point{}
		instructions := numToDir(&main_robot, code)

		for i := 0; i < num_dir_robots; i++ {
			dir_robots = append(dir_robots, dir_ref["A"])
		}

		for _, robot := range dir_robots {
			instructions = dirToDir(&robot, instructions, dir_memo)
		}

		num, _ := strconv.Atoi(strings.TrimRight(code, "A"))
		complexity_sum += len(instructions) * num
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

		// Add the higher-cost instructions first except in the above cases
		// where the empty space would be crossed. The "<" key is the highest
		// cost move because it takes the most steps to reach from the "A"
		// key, for example.
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

func dirToDir(robot *Point, code string, dir_memo map[string]string) string {
	var builder strings.Builder
	chars := strings.Split(code, "")
	current_symbol := "A"

	for _, ch := range chars {
		delta := dir_ref[ch].Minus(*robot) // Desired position minus current position
		next_symbol := ch
		sub_instructions := "" // Instructions to move from current to next

		// Look up step sequence for previously encountered scenarios
		if val, ok := dir_memo[current_symbol+next_symbol]; ok {
			builder.WriteString(val)
			current_symbol = next_symbol
			*robot = dir_ref[next_symbol]
			continue
		}

		// Avoid crossing empty space in direction pad
		if robot.Col == 0 && dir_ref[ch].Row == 0 {
			for i := 0; i < delta.Col; i++ {
				sub_instructions += ">"
				*robot = robot.Plus(Deltas[">"])
			}
			for i := 0; i < delta.Row*(-1); i++ {
				sub_instructions += "^"
				*robot = robot.Plus(Deltas["^"])
			}

			sub_instructions += "A"
			if _, ok := dir_memo[current_symbol+next_symbol]; !ok {
				dir_memo[current_symbol+next_symbol] = sub_instructions
			}
			builder.WriteString(sub_instructions)
			current_symbol = next_symbol
			continue
		}

		if robot.Row == 0 && dir_ref[ch].Col == 0 {
			for i := 0; i < delta.Row; i++ {
				sub_instructions += "v"
				*robot = robot.Plus(Deltas["v"])
			}
			for i := 0; i < delta.Col*(-1); i++ {
				sub_instructions += "<"
				*robot = robot.Plus(Deltas["<"])
			}

			sub_instructions += "A"
			if _, ok := dir_memo[current_symbol+next_symbol]; !ok {
				dir_memo[current_symbol+next_symbol] = sub_instructions
			}
			builder.WriteString(sub_instructions)
			current_symbol = next_symbol
			continue
		}

		// Add the higher-cost instructions first except in the above cases
		// where the empty space would be crossed. The "<" key is the highest
		// cost move because it takes the most steps to reach from the "A"
		// key, for example.
		if delta.Col < 0 {
			for i := 0; i < delta.Col*(-1); i++ {
				sub_instructions += "<"
				*robot = robot.Plus(Deltas["<"])
			}
		}

		if delta.Row > 0 {
			for i := 0; i < delta.Row; i++ {
				sub_instructions += "v"
				*robot = robot.Plus(Deltas["v"])
			}
		}

		if delta.Row < 0 {
			for i := 0; i < delta.Row*(-1); i++ {
				sub_instructions += "^"
				*robot = robot.Plus(Deltas["^"])
			}
		}

		if delta.Col > 0 {
			for i := 0; i < delta.Col; i++ {
				sub_instructions += ">"
				*robot = robot.Plus(Deltas[">"])
			}
		}

		sub_instructions += "A"
		if _, ok := dir_memo[current_symbol+next_symbol]; !ok {
			dir_memo[current_symbol+next_symbol] = sub_instructions
		}
		builder.WriteString(sub_instructions)
		current_symbol = next_symbol
	}

	return builder.String()
}
