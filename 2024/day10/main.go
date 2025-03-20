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

var Dir = map[string]Point{
	"UP":    {Row: -1, Col: 0},
	"RIGHT": {Row: 0, Col: 1},
	"DOWN":  {Row: 1, Col: 0},
	"LEFT":  {Row: 0, Col: -1},
}

var Endpoints = map[Point][]Point{}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading topographic map...\n")
	grid := [][]int{}
	trailheads := []Point{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		line_splice := strings.Split(line, "")
		row := []int{}
		for i := 0; i < len(line_splice); i++ {
			num, err := strconv.Atoi(line_splice[i])
			if err != nil {
				fmt.Printf("String conversion error: %v", err)
			}
			if num == 0 { // Locate potential trailheads
				p := Point{Row: len(grid), Col: i}
				trailheads = append(trailheads, p)
			}
			row = append(row, num)
		}
		grid = append(grid, row)
	}

	fmt.Printf("Calculating trail scores...\n")
	total_score := 0
	total_rating := 0
	for i := 0; i < len(trailheads); i++ {
		Endpoints[trailheads[i]] = []Point{}
		score, rating := followPath(grid, trailheads[i], trailheads[i], 0)
		total_score += score
		total_rating += rating
	}

	fmt.Printf("[Part 1] Trail scores total: %d\n", total_score)
	fmt.Printf("[Part 2] Trail ratings total: %d\n", total_rating)
}

func followPath(grid [][]int, trailhead Point, current_point Point, current_elevation int) (int, int) {
	total_score := 0
	total_rating := 0

	for _, delta := range Dir {
		next_point := current_point.Add(delta)
		if !isOutOfBounds(grid, next_point) {
			next_elevation := grid[next_point.Row][next_point.Col]
			if current_elevation == 8 && next_elevation == 9 {
				if !isInSlice(next_point, Endpoints[trailhead]) {
					total_score++
					Endpoints[trailhead] = append(Endpoints[trailhead], next_point)
				}
				total_rating++
			} else if next_elevation == current_elevation+1 {
				score, rating := followPath(grid, trailhead, next_point, next_elevation)
				total_score += score
				total_rating += rating
			}
		}
	}

	return total_score, total_rating
}

func isOutOfBounds(grid [][]int, p Point) bool {
	grid_height := len(grid)
	grid_width := len(grid[0])

	if p.Row < 0 || p.Row >= grid_height || p.Col < 0 || p.Col >= grid_width {
		return true
	}

	return false
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
