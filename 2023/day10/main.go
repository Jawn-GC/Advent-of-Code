package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Point struct {
	Row int
	Col int
}

func (p Point) add(other Point) Point {
	return Point{
		Row: p.Row + other.Row,
		Col: p.Col + other.Col,
	}
}

var deltas = map[string]Point{
	"north": {Row: -1, Col: 0},
	"east":  {Row: 0, Col: 1},
	"south": {Row: 1, Col: 0},
	"west":  {Row: 0, Col: -1},
}

var pipe_map = map[string][]string{
	"|": {"north", "south"},
	"-": {"east", "west"},
	"L": {"north", "east"},
	"J": {"north", "west"},
	"7": {"south", "west"},
	"F": {"south", "east"},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Reading tilemap...")
	grid := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		new_row := strings.Split(line, "")
		grid = append(grid, new_row)
	}

	fmt.Println("Identifying starting tile...")
	var start Point
	for i, row := range grid {
		for j, tile := range row {
			if tile == "S" {
				start = Point{Row: i, Col: j}
			}
		}
	}
	starting_tile := identifyTile(grid, start)
	grid[start.Row][start.Col] = starting_tile

	fmt.Println("Finding tiles in pipe loop...")
	pipe_loop := getPipeLoop(grid, start)
	fmt.Printf("[Part 1] The furthest point in the loop is %d tiles from the start\n", len(pipe_loop)/2)
}

// Each pipe tile connects to exactly two other pipe tiles.
func identifyTile(grid [][]string, point Point) string {
	connected_directions := []string{}
	for dir, delta := range deltas {
		adj_point := point.add(delta)
		if isOOB(grid, adj_point) {
			continue
		}
		adj_tile := grid[adj_point.Row][adj_point.Col]

		// If-statement is separated into else-ifs for clarity
		if dir == "north" && (adj_tile == "|" || adj_tile == "F" || adj_tile == "7") {
			connected_directions = append(connected_directions, dir)
		} else if dir == "east" && (adj_tile == "-" || adj_tile == "7" || adj_tile == "J") {
			connected_directions = append(connected_directions, dir)
		} else if dir == "south" && (adj_tile == "|" || adj_tile == "J" || adj_tile == "L") {
			connected_directions = append(connected_directions, dir)
		} else if dir == "west" && (adj_tile == "-" || adj_tile == "L" || adj_tile == "F") {
			connected_directions = append(connected_directions, dir)
		}
	}

	for tile, dirs := range pipe_map {
		if equalSets(connected_directions, dirs) {
			return tile
		}
	}

	return "."
}

// The starting point is assumed to be part of a pipe loop.
// OOB is not checked for because the direction of motion is always toward
// a connected pipe tile in the loop and thus within the bounds of the grid.
func getPipeLoop(grid [][]string, start Point) []Point {
	loop_points := []Point{}
	current_point := start
	loop_points = append(loop_points, current_point)
	forward_direction := ""

	for {
		current_tile := grid[current_point.Row][current_point.Col]
		dirs := pipe_map[current_tile]

		// Choose either direction along the pipe loop to start and continue
		// forward until the start is reached again.
		if forward_direction == "" {
			forward_direction = dirs[0]
		} else {
			backward_direction := getOppositeDir(forward_direction)
			for _, dir := range dirs {
				if dir != backward_direction {
					forward_direction = dir
				}
			}
		}

		delta := deltas[forward_direction]
		current_point = current_point.add(delta)
		loop_points = append(loop_points, current_point)

		if current_point == start {
			break
		}
	}

	return loop_points
}

// Assumes grid is non-empty and rectangular
func isOOB(grid [][]string, point Point) bool {
	height := len(grid)
	width := len(grid[0])

	return (point.Row < 0 || point.Row >= height || point.Col < 0 || point.Col >= width)
}

// Checks if two slices have the same entries regardless of order
func equalSets(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	counts := make(map[string]int)

	for _, item := range a {
		counts[item]++
	}
	for _, item := range b {
		counts[item]--
		if counts[item] < 0 {
			return false
		}
	}
	return true
}

func getOppositeDir(dir string) string {
	if dir == "north" {
		return "south"
	} else if dir == "east" {
		return "west"
	} else if dir == "south" {
		return "north"
	} else if dir == "west" {
		return "east"
	}
	return ""
}
