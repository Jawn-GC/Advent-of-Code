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
	"UP":    {Row: -1, Col: 0},
	"RIGHT": {Row: 0, Col: 1},
	"DOWN":  {Row: 1, Col: 0},
	"LEFT":  {Row: 0, Col: -1},
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Parsing racetrack...\n")
	racetrack := [][]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		row := []string{}
		row = append(row, strings.Split(line, "")...)
		racetrack = append(racetrack, row)
	}

	fmt.Printf("Finding trackpoints...\n")
	start, end := getEndpoints(racetrack)
	trackpoints := getTrackPoints(racetrack, start, end)
	fmt.Printf("Calculating cheat times...\n")
	cheat_times := getCheatTimes(racetrack, trackpoints)

	max_time := len(trackpoints) - 1
	num_times := 0
	for _, time := range cheat_times {
		if max_time-time >= 100 {
			num_times++
		}
	}

	fmt.Printf("Number of cheats that save at least 100 picoseconds: %d\n", num_times)
}

func getEndpoints(grid [][]string) (Point, Point) {
	start, end := Point{}, Point{}
	start_found, end_found := false, false

	for i := 0; i < len(grid); i++ {
		for j := 0; j < len(grid[0]); j++ {
			if grid[i][j] == "S" {
				start_found = true
				start.Row = i
				start.Col = j
			} else if grid[i][j] == "E" {
				end_found = true
				end.Row = i
				end.Col = j
			}

			if start_found && end_found {
				return start, end
			}
		}
	}

	return Point{}, Point{}
}

func getTrackPoints(grid [][]string, start, end Point) map[Point]int {
	height, width := len(grid), len(grid[0])
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}
	visited[start.Row][start.Col] = true

	current := start
	points := map[Point]int{}
	points[start] = 0
	index := 1 // Initialized at 1 because "start" was already added
	for current != end {
		for _, delta := range Directions {
			adj_point := current.Add(delta)

			// The racetrack has no branches, so each trackpoint only touches up to
			// two other trackpoints, the one ahead and the one behind.
			if !visited[adj_point.Row][adj_point.Col] && grid[adj_point.Row][adj_point.Col] != "#" {
				points[adj_point] = index
				current = adj_point
				visited[adj_point.Row][adj_point.Col] = true
				index++
				break
			}
		}
	}

	return points
}

func getCheatTimes(grid [][]string, trackpoints map[Point]int) []int {
	cheat_times := []int{}
	max_time := len(trackpoints) - 1

	height, width := len(grid), len(grid[0])
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			// Ignore the border tiles
			if i > 0 && i < height-1 && j > 0 && j < width-1 && grid[i][j] == "#" {
				// If up and down tiles are part of the racetrack
				if grid[i-1][j] != "#" && grid[i+1][j] != "#" {
					up := Point{Row: i - 1, Col: j}
					down := Point{Row: i + 1, Col: j}

					if trackpoints[up] < trackpoints[down] {
						time := trackpoints[up] + 2 + (max_time - trackpoints[down])
						cheat_times = append(cheat_times, time)
					} else {
						time := trackpoints[down] + 2 + (max_time - trackpoints[up])
						cheat_times = append(cheat_times, time)
					}
				}

				// If left and right tiles are part of the racetrack
				if grid[i][j-1] != "#" && grid[i][j+1] != "#" {
					left := Point{Row: i, Col: j - 1}
					right := Point{Row: i, Col: j + 1}

					if trackpoints[left] < trackpoints[right] {
						time := trackpoints[left] + 2 + (max_time - trackpoints[right])
						cheat_times = append(cheat_times, time)
					} else {
						time := trackpoints[right] + 2 + (max_time - trackpoints[left])
						cheat_times = append(cheat_times, time)
					}
				}
			}
		}
	}

	return cheat_times
}
