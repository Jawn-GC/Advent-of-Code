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

func (p Point) Add(other Point) Point {
	return Point{
		Row: p.Row + other.Row,
		Col: p.Col + other.Col,
	}
}

var Directions = map[string]Point{
	"^": {Row: -1, Col: 0},
	">": {Row: 0, Col: 1},
	"v": {Row: 1, Col: 0},
	"<": {Row: 0, Col: -1},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading grid...\n")
	grid := [][]string{}
	moves := []string{}
	input_type := "grid"
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if len(line) == 0 {
			input_type = "moves"
			fmt.Printf("Reading instructions...\n")
			continue
		}
		if input_type == "grid" {
			grid_row := strings.Split(line, "")
			grid = append(grid, grid_row)
		} else if input_type == "moves" {
			moves = append(moves, strings.Split(line, "")...)
		}
	}

	fmt.Printf("Following instructions...\n")
	followInstructions(grid, moves)

	fmt.Printf("Calculating GPS coordinates...\n")
	gps_sum := calcGPSSum(grid)
	fmt.Printf("The sum of GPS coordinates is %d\n", gps_sum)
}

func calcGPSSum(grid [][]string) int {
	gps_sum := 0

	barrel_points := findBarrels(grid)
	for _, point := range barrel_points {
		gps_sum += 100*point.Row + point.Col
	}

	return gps_sum
}

func followInstructions(grid [][]string, moves []string) {
	for _, move := range moves {
		r_pos := findRobot(grid)
		tryMove(r_pos, move, grid)
	}
}

func tryMove(entity_point Point, dir string, grid [][]string) {
	entity := grid[entity_point.Row][entity_point.Col]

	adj_point := entity_point.Add(Directions[dir])
	adj_entity := grid[adj_point.Row][adj_point.Col]
	if entity == "@" || entity == "O" {
		if adj_entity == "O" {
			tryMove(adj_point, dir, grid)
			// The adjacent entity may have changed after the move attempt,
			// so we need to update it.
			adj_entity = grid[adj_point.Row][adj_point.Col]
		}

		// Move the entity into the adjacent space and free up the previous space.
		if adj_entity == "." {
			grid[adj_point.Row][adj_point.Col] = entity
			grid[entity_point.Row][entity_point.Col] = "."
		}
	}
}

func findRobot(grid [][]string) Point {
	for i, grid_row := range grid {
		for j, entity := range grid_row {
			if entity == "@" {
				return Point{Row: i, Col: j}
			}
		}
	}
	return Point{Row: -1, Col: -1}
}

func findBarrels(grid [][]string) []Point {
	barrels := []Point{}
	for i, grid_row := range grid {
		for j, entity := range grid_row {
			if entity == "O" {
				barrels = append(barrels, Point{Row: i, Col: j})
			}
		}
	}
	return barrels
}
