package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Robot struct {
	X  int
	Y  int
	Vx int
	Vy int
}

func (r *Robot) Move() {
	r.X += r.Vx
	r.Y += r.Vy
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Parsing initial conditions...\n")
	robots := []Robot{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Fields(line)
		pos_str := strings.TrimLeft(temp[0], "p=")
		vel_str := strings.TrimLeft(temp[1], "v=")

		pos_strs := strings.Split(pos_str, ",")
		vel_strs := strings.Split(vel_str, ",")

		x, err1 := strconv.Atoi(pos_strs[0])
		y, err2 := strconv.Atoi(pos_strs[1])
		vx, err3 := strconv.Atoi(vel_strs[0])
		vy, err4 := strconv.Atoi(vel_strs[1])

		if err1 != nil || err2 != nil || err3 != nil || err4 != nil {
			fmt.Printf("String conversion error")
			return
		}

		robots = append(robots, Robot{X: x, Y: y, Vx: vx, Vy: vy})
	}

	grid_height := 103
	grid_width := 101
	steps := 100
	cur_step := 1
	fmt.Printf("Moving robots...\n")
	for i := 0; i < steps; i++ {
		grid := newGrid(grid_height, grid_width)
		for j := 0; j < len(robots); j++ {
			robot := &robots[j]
			robot.Move()
			robot.X %= grid_width
			if robot.X < 0 {
				robot.X += grid_width
			}
			robot.Y %= grid_height
			if robot.Y < 0 {
				robot.Y += grid_height
			}
			grid[robot.Y][robot.X] = "*"
		}
		//if (cur_step-89)%grid_height == 0 || (cur_step-11)%grid_width == 0 {
		//	fmt.Printf("Step number: %d\n", cur_step)
		//	printGrid(grid)
		//}
		cur_step++
	}

	fmt.Printf("Locating robots...\n")
	ul, ur, lr, ll := 0, 0, 0, 0
	for _, robot := range robots {
		if robot.X < grid_width/2 && robot.Y < grid_height/2 {
			ul++
		}
		if robot.X > grid_width/2 && robot.Y < grid_height/2 {
			ur++
		}
		if robot.X > grid_width/2 && robot.Y > grid_height/2 {
			lr++
		}
		if robot.X < grid_width/2 && robot.Y > grid_height/2 {
			ll++
		}
	}
	safety_factor := ul * ur * lr * ll

	fmt.Printf("Safety factor after %d steps: %d\n", steps, safety_factor)
}

func newGrid(height, width int) [][]string {
	grid := [][]string{}
	for range height {
		grid_row := []string{}
		for range width {
			grid_row = append(grid_row, ".")
		}
		grid = append(grid, grid_row)
	}
	return grid
}

func printGrid(grid [][]string) {
	height := len(grid)
	width := len(grid[0])

	for i := range height {
		for j := range width {
			fmt.Printf(grid[i][j])
		}
		fmt.Printf("\n")
	}
}
