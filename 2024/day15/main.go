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

type Entity struct {
	Ent_Type string
	Position Point
}

var Directions = map[string]Point{
	"^": {Row: -1, Col: 0},
	">": {Row: 0, Col: 1},
	"v": {Row: 1, Col: 0},
	"<": {Row: 0, Col: -1},
}

var mode = "normal" // normal or wide
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
			if mode == "normal" {
				grid = append(grid, grid_row)
			}
		} else if input_type == "moves" {
			moves = append(moves, strings.Split(line, "")...)
		}
	}

	fmt.Printf("Following instructions...\n")
	followInstructions(grid, moves)

	fmt.Printf("Calculating GPS coordinates...\n")
	gps_sum := calcGPSSum(grid)
	fmt.Printf("The sum of GPS coordinates is %d\n", gps_sum)

	// Print final grid
	//for _, grid_row := range grid {
	//	for _, str := range grid_row {
	//		fmt.Printf("%s", str)
	//	}
	//	fmt.Printf("\n")
	//}
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
	for i := 0; i < len(moves); i++ {
		robot := Entity{Ent_Type: "@", Position: findRobot(grid)}
		entities_to_move := []Entity{}
		if canMove(robot, moves[i], grid, &entities_to_move) {
			moveEntities(grid, entities_to_move, moves[i])
		}
	}
}

func canMove(entity Entity, dir string, grid [][]string, entities_to_move *[]Entity) bool {
	adj_point := entity.Position.Add(Directions[dir])
	adj_ent_type := grid[adj_point.Row][adj_point.Col]
	adj_entity := Entity{Ent_Type: adj_ent_type, Position: adj_point}

	// Robot
	if entity.Ent_Type == "@" {
		if adj_entity.Ent_Type == "." {
			*entities_to_move = append(*entities_to_move, entity)
			return true
		} else if adj_entity.Ent_Type == "O" {
			if canMove(adj_entity, dir, grid, entities_to_move) {
				*entities_to_move = append(*entities_to_move, entity)
				return true
			}
		}
	}

	// Normal Barrels
	if entity.Ent_Type == "O" {
		if adj_entity.Ent_Type == "." {
			*entities_to_move = append(*entities_to_move, entity)
			return true
		} else if adj_entity.Ent_Type == "O" {
			if canMove(adj_entity, dir, grid, entities_to_move) {
				*entities_to_move = append(*entities_to_move, entity)
				return true
			}
		}
	}

	return false
}

func moveEntities(grid [][]string, entities []Entity, dir string) {
	// Erase symbols at entity positions
	for _, entity := range entities {
		grid[entity.Position.Row][entity.Position.Col] = "."
	}

	// Add symbols to new entity positions
	for _, entity := range entities {
		adj_point := entity.Position.Add(Directions[dir])
		grid[adj_point.Row][adj_point.Col] = entity.Ent_Type
	}
}

func findRobot(grid [][]string) Point {
	for i, grid_row := range grid {
		for j, entity_type := range grid_row {
			if entity_type == "@" {
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
			if mode == "normal" && entity == "O" {
				barrels = append(barrels, Point{Row: i, Col: j})
			} else if mode == "wide" && entity == "[" {
				barrels = append(barrels, Point{Row: i, Col: j})
			}
		}
	}
	return barrels
}
