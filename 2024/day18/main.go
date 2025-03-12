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

func (p Point) Add(other Point) Point {
	return Point{
		Row: p.Row + other.Row,
		Col: p.Col + other.Col,
	}
}

var Directions = []Point{
	{Row: -1, Col: 0},
	{Row: 0, Col: 1},
	{Row: 1, Col: 0},
	{Row: 0, Col: -1},
}

type State struct {
	Pos   Point
	Steps int
}

type Queue []State

func (q *Queue) Push(x State) {
	*q = append(*q, x)
}

func (q *Queue) Pop() State {
	old := *q
	item := old[0]
	*q = old[1:]
	return item
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading byte locations...\n")
	byte_positions := []Point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ",")
		x, _ := strconv.Atoi(temp[0])
		y, _ := strconv.Atoi(temp[1])
		byte_positions = append(byte_positions, Point{Row: y, Col: x})
	}

	height, width := 71, 71
	grid := makeGrid(height, width)
	start := Point{Row: 0, Col: 0}
	end := Point{Row: height - 1, Col: width - 1}

	fmt.Printf("Corrupting grid...\n")
	var block_pos Point
	for _, p := range byte_positions {
		corruptSpace(grid, p)
		min_steps := getNumSteps(grid, start, end)
		if min_steps == -1 {
			block_pos = p
			break
		}
	}

	fmt.Printf("No path available after (%d,%d) is corrupted.\n", block_pos.Col, block_pos.Row)
}

func getNumSteps(grid [][]string, start, end Point) int {
	queue := Queue{}
	height, width := len(grid), len(grid[0])

	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	// Initial Point
	queue.Push(State{Pos: start, Steps: 0})
	visited[start.Row][start.Col] = true

	for len(queue) > 0 {
		current := queue.Pop()
		if current.Pos == end {
			return current.Steps
		}

		for _, delta := range Directions {
			new_point := current.Pos.Add(delta)

			if new_point.Row < 0 || new_point.Row >= height || new_point.Col < 0 || new_point.Col >= width || grid[new_point.Row][new_point.Col] == "#" {
				continue
			}

			if !visited[new_point.Row][new_point.Col] {
				visited[new_point.Row][new_point.Col] = true
				queue.Push(State{Pos: new_point, Steps: current.Steps + 1})
			}
		}
	}

	return -1 // No valid path
}

func makeGrid(height, width int) [][]string {
	grid := [][]string{}

	for i := 0; i < height; i++ {
		grid_row := []string{}
		for j := 0; j < width; j++ {
			grid_row = append(grid_row, ".")
		}
		grid = append(grid, grid_row)
	}

	return grid
}

// It is assumed that the corrupted points are in bounds
func corruptGrid(grid [][]string, bytes []Point, num_bytes int) {
	for i := 0; i < num_bytes; i++ {
		grid[bytes[i].Row][bytes[i].Col] = "#"
	}
}

func corruptSpace(grid [][]string, p Point) {
	grid[p.Row][p.Col] = "#"
}
