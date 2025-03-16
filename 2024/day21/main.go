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

// Symbol pair to path reference.
// Obtained by printing the map after use in Part 1.
// map[<<:A <A:>>^A <^:>^A <v:>A >>:A >A:^A >^:<^A >v:<A A<:v<<A A>:vA AA:A A^:<A Av:<vA ^<:v<A ^>:v>A ^A:>A ^^:A v<:<A v>:>A vA:^>A vv:A]
var pair_sequences = make(map[string]map[string]int)

func init_pair_sequences() {
	pair_sequences["AA"] = map[string]int{"AA": 1}
	pair_sequences["A^"] = map[string]int{"A<": 1, "<A": 1}
	pair_sequences["A>"] = map[string]int{"Av": 1, "vA": 1}
	pair_sequences["Av"] = map[string]int{"A<": 1, "<v": 1, "vA": 1}
	pair_sequences["A<"] = map[string]int{"Av": 1, "v<": 1, "<<": 1, "<A": 1}

	pair_sequences["^A"] = map[string]int{"A>": 1, ">A": 1}
	pair_sequences["^^"] = map[string]int{"AA": 1}
	pair_sequences["^>"] = map[string]int{"Av": 1, "v>": 1, ">A": 1}
	// No entry required for "^v"
	pair_sequences["^<"] = map[string]int{"Av": 1, "v<": 1, "<A": 1}

	pair_sequences[">A"] = map[string]int{"A^": 1, "^A": 1}
	pair_sequences[">^"] = map[string]int{"A<": 1, "<^": 1, "^A": 1}
	pair_sequences[">>"] = map[string]int{"AA": 1}
	pair_sequences[">v"] = map[string]int{"A<": 1, "<A": 1}
	// No entry required for "><"

	pair_sequences["vA"] = map[string]int{"A^": 1, "^>": 1, ">A": 1}
	// No entry required for "v^"
	pair_sequences["v>"] = map[string]int{"A>": 1, ">A": 1}
	pair_sequences["vv"] = map[string]int{"AA": 1}
	pair_sequences["v<"] = map[string]int{"A<": 1, "<A": 1}

	pair_sequences["<A"] = map[string]int{"A>": 1, ">>": 1, ">^": 1, "^A": 1}
	pair_sequences["<^"] = map[string]int{"A>": 1, ">^": 1, "^A": 1}
	// No entry required for "<>"
	pair_sequences["<v"] = map[string]int{"A>": 1, ">A": 1}
	pair_sequences["<<"] = map[string]int{"AA": 1}
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
	complexity_sum1 := 0
	complexity_sum2 := 0
	num_dir_robots1 := 2
	num_dir_robots2 := 25
	dir_memo := make(map[string]string) // Key: "start_symbol"+"end_symbol"
	init_pair_sequences()

	// Part 1
	for _, code := range codes {
		main_robot := num_ref["A"]
		dir_robots := []Point{}
		instructions := numToDir(&main_robot, code)

		for i := 0; i < num_dir_robots1; i++ {
			dir_robots = append(dir_robots, dir_ref["A"])
		}

		for _, robot := range dir_robots {
			instructions = dirToDir(&robot, instructions, dir_memo)
		}

		num, _ := strconv.Atoi(strings.TrimRight(code, "A"))
		complexity_sum1 += len(instructions) * num
	}

	// Part 2
	// Uses a different method because strings eventually become too large.
	for _, code := range codes {
		main_robot := num_ref["A"]
		instructions := numToDir(&main_robot, code)
		pair_counts := get_initial_pair_counts(instructions)

		for i := 0; i < num_dir_robots2; i++ {
			pair_counts = get_next_pair_counts(pair_counts)
		}

		num_instructions := 0
		num, _ := strconv.Atoi(strings.TrimRight(code, "A"))
		for _, val := range pair_counts {
			num_instructions += val
		}
		complexity_sum2 += num_instructions * num
	}

	fmt.Printf("Complexity sum [Part 1]: %d\n", complexity_sum1)
	fmt.Printf("Complexity sum [Part 2]: %d\n", complexity_sum2)
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

func get_initial_pair_counts(instructions string) map[string]int {
	pair_counts := make(map[string]int)
	for key, _ := range pair_sequences {
		pair_counts[key] = 0
	}

	points := []string{"A"}
	points = append(points, strings.Split(instructions, "")...)
	for i := 0; i < len(points)-1; i++ {
		pair_counts[points[i]+points[i+1]]++
	}

	return pair_counts
}

func get_next_pair_counts(pair_counts map[string]int) map[string]int {
	new_counts := make(map[string]int)
	for key, _ := range pair_sequences {
		new_counts[key] = 0
	}

	for pair, count := range pair_counts {
		seq := pair_sequences[pair]
		for key, val := range seq {
			new_counts[key] += val * count
		}
	}

	return new_counts
}
