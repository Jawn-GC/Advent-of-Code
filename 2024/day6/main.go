package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type State struct {
	Row int
	Col int
	Dir string
}

type Delta struct {
	Dy int
	Dx int
}

var Directions = map[string]Delta{
	"UP":    {-1, 0},
	"RIGHT": {0, 1},
	"DOWN":  {1, 0},
	"LEFT":  {0, -1},
}

var DirectionsOrder = map[int]string{
	0: "UP",
	1: "RIGHT",
	2: "DOWN",
	3: "LEFT",
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v", filename, err)
	}
	defer file.Close()

	lab_grid := [][]string{}
	guard_found := false
	guard := State{
		Row: 0,
		Col: 0,
		Dir: "UP", // The guard is assumed to be facing up at the start
	}
	fmt.Printf("Parsing grid...\n")
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []string{}
		row = append(row, strings.Split(line, "")...)
		lab_grid = append(lab_grid, row)

		if !guard_found {
			col := getIndex("^", row) // Different symbols need to be searched for if the guard isn't facing up at the start
			if col != -1 {
				guard.Col = col
				guard_found = true
				continue
			}
			guard.Row++
		}
	}

	grid_copy := getGridCopy(lab_grid)
	guard_copy := guard

	// Path for original lab grid
	fmt.Printf("Predicting path...\n")
	x_states, _ := predictPath(grid_copy, guard_copy)
	fmt.Printf("The guard visited %d unique tiles.\n", len(x_states))

	// Paths for modified lab grid
	num_loops := 0
	fmt.Printf("\nFinding potential infinite loops...\n")
	for _, state := range x_states {
		if state.Row == guard.Row && state.Col == guard.Col {
			continue
		}
		grid_copy = getGridCopy(lab_grid)
		guard_copy = guard
		grid_copy[state.Row][state.Col] = "#"
		_, loop_found := predictPath(grid_copy, guard_copy)
		if loop_found {
			num_loops++
		}
	}
	fmt.Printf("There are %d potential loops.\n", num_loops)
}

func getIndex[T comparable](val T, slice []T) int {
	for i, v := range slice {
		if v == val {
			return i
		}
	}
	return -1
}

// Marks positions that the guard has occupied with "X".
// Returns the states of the guard at X-marked positions.
// Returns true/false if a loop is found/not found.
func predictPath(lab_grid [][]string, guard State) ([]State, bool) {
	dirTracker := 0 // The guard starts facing "UP"
	states := []State{}
	height := len(lab_grid)
	width := len(lab_grid[0])
	tile_ahead := ""
	delta := Directions[guard.Dir]
	loop_found := false

outerLoop:
	for {
		if lab_grid[guard.Row][guard.Col] != "X" {
			lab_grid[guard.Row][guard.Col] = "X"
			states = append(states, State{Row: guard.Row, Col: guard.Col, Dir: guard.Dir})
		}
	innerLoop:
		for {
			ahead_row := guard.Row + delta.Dy
			ahead_col := guard.Col + delta.Dx

			// Check if the tile in front of the guard is out of bounds
			if ahead_row < 0 || ahead_row >= height || ahead_col < 0 || ahead_col >= width {
				break outerLoop
			}

			tile_ahead = lab_grid[ahead_row][ahead_col]

			if tile_ahead == "." || tile_ahead == "X" {
				break innerLoop
			}

			// Rotate the guard until the path directly in front is unobstructed
			if tile_ahead == "#" {
				dirTracker = (dirTracker + 1) % 4 // Rotates the guard 90 degrees clockwise
				guard.Dir = DirectionsOrder[dirTracker]
				delta = Directions[guard.Dir]
			}
		}

		guard.Row += delta.Dy
		guard.Col += delta.Dx

		if getIndex(guard, states) != -1 {
			loop_found = true
			break outerLoop
		}
	}

	return states, loop_found
}

func getGridCopy(grid [][]string) [][]string {
	grid_copy := [][]string{}
	for _, row := range grid {
		row_copy := make([]string, len(row))
		copy(row_copy, row)
		grid_copy = append(grid_copy, row_copy)
	}
	return grid_copy
}
