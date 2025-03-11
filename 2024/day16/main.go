package main

import (
	"bufio"
	"container/heap"
	"fmt"
	"os"
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

var Directions = map[int]Point{
	0: {Row: -1, Col: 0},
	1: {Row: 0, Col: 1},
	2: {Row: 1, Col: 0},
	3: {Row: 0, Col: -1},
}

type State struct {
	Pos       Point
	Direction int
	Score     int
}

// Implement the heap.Interface with the following functions
type PriorityQueue []State

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].Score < pq[j].Score
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x interface{}) {
	*pq = append(*pq, x.(State))
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
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

	fmt.Printf("Reading maze...\n")
	maze := [][]rune{}
	start := Point{}
	end := Point{}
	row_index := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		maze_row := []rune{}
		for col, r := range line {
			if r == 'S' {
				start.Row = row_index
				start.Col = col
				maze_row = append(maze_row, '.')
			} else if r == 'E' {
				end.Row = row_index
				end.Col = col
				maze_row = append(maze_row, '.')
			} else {
				maze_row = append(maze_row, r)
			}
		}
		maze = append(maze, maze_row)
		row_index++
	}

	fmt.Printf("Calculating mininum score...\n")
	min_score := getMinScore(maze, start, end)
	fmt.Printf("The lowest possible score is %d\n", min_score)
}

func getMinScore(maze [][]rune, start, end Point) int {
	height, width := len(maze), len(maze[0])

	visited := make([][][]bool, height)
	for i := range visited {
		visited[i] = make([][]bool, width)
		for j := range visited[i] {
			visited[i][j] = make([]bool, len(Directions))
		}
	}

	initial_state := State{Pos: start, Direction: 1, Score: 0}
	pq := &PriorityQueue{}
	heap.Init(pq)

	// Add Initial States
	for new_dir, _ := range Directions {
		var rotations int
		difference := new_dir - initial_state.Direction
		if difference < 0 {
			difference *= -1
		}

		if difference != 2 && difference != 0 {
			rotations = 1
		} else {
			rotations = difference
		}

		heap.Push(pq, State{Pos: start, Direction: new_dir, Score: 1000 * rotations})
		visited[start.Row][start.Col][new_dir] = true
	}

	for pq.Len() > 0 {
		current := heap.Pop(pq).(State)

		if current.Pos == end {
			return current.Score
		}

		for new_dir, delta := range Directions {
			nextPos := current.Pos.Add(delta)

			if nextPos.Row < 0 || nextPos.Row >= height || nextPos.Col < 0 || nextPos.Col >= width || maze[nextPos.Row][nextPos.Col] == '#' {
				continue
			}

			var rotations int
			var new_score int
			difference := new_dir - current.Direction
			if difference < 0 {
				difference *= -1
			}

			if difference != 2 && difference != 0 {
				rotations = 1
			} else {
				rotations = difference
			}
			new_score = current.Score + 1000*rotations + 1

			if !visited[nextPos.Row][nextPos.Col][new_dir] {
				visited[nextPos.Row][nextPos.Col][new_dir] = true
				heap.Push(pq, State{Pos: nextPos, Direction: new_dir, Score: new_score})
			}
		}
	}

	return -1 // Invalid Path
}
