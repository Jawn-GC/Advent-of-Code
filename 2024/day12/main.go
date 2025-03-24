package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// The directional fields determine if the plot
// has a wall on a particular side.
type Plot struct {
	Pos   Point
	Type  string
	Up    bool
	Right bool
	Down  bool
	Left  bool
}

func NewPlot(pos Point, plot_type string) Plot {
	return Plot{
		Pos:   pos,
		Type:  plot_type,
		Up:    true,
		Right: true,
		Down:  true,
		Left:  true,
	}
}

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

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Parsing garden plots...\n")
	garden := [][]Plot{}
	row_count := 0
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		plots := strings.Split(line, "")
		garden_row := []Plot{}
		for col_count, p := range plots {
			garden_row = append(garden_row, NewPlot(Point{Row: row_count, Col: col_count}, p))
		}
		garden = append(garden, garden_row)
		row_count++
	}

	fmt.Printf("Finding plot regions...\n")
	regions := findPlotRegions(garden)

	fmt.Printf("Calculating fence prices by perimeter...\n")
	total_price_perimeter := 0
	for _, region_group := range regions {
		for _, region := range region_group {
			total_price_perimeter += calcFencePricingByPerimeter(region)
		}
	}

	total_price_sides := 0
	fmt.Printf("Calculating fence prices by sides...\n")
	for _, region_group := range regions {
		for _, region := range region_group {
			total_price_sides += calcFencePricingBySides(region)
		}
	}
	fmt.Printf("[Part 1] The total price of fencing (by perimeter) for all regions is %d\n", total_price_perimeter)
	fmt.Printf("[Part 2] The total price of fencing (by sides) for all regions is %d\n", total_price_sides)
}

func findPlotRegions(garden [][]Plot) map[string][][]*Plot {
	regions := map[string][][]*Plot{}

	for i, garden_row := range garden {
		for j, plot := range garden_row {
			// If a plot type does not have a group, add one.
			if _, ok := regions[plot.Type]; !ok {
				regions[plot.Type] = [][]*Plot{}
			}

			// Determine if  the current plot is already part of a region.
			// If so, move on to the next plot.
			isNewPlot := true
			for _, region := range regions[plot.Type] {
				if isInSlice(&garden[i][j], region) {
					isNewPlot = false
				}
			}
			if !isNewPlot {
				continue
			}

			// Create a new region that spreads out from the current plot.
			new_region := []*Plot{&garden[i][j]}
			spread(garden, &new_region, Point{Row: i, Col: j})
			regions[plot.Type] = append(regions[plot.Type], new_region)
		}
	}

	return regions
}

func spread(garden [][]Plot, current_region *[]*Plot, current_point Point) {
	current_plot := &garden[current_point.Row][current_point.Col]
	garden_height := len(garden)
	garden_width := len(garden[0])
	for dir_str, dir := range Dir {
		// Check if the adjacent space is OOB.
		adj_point := current_point.Add(dir)
		if adj_point.Row < 0 || adj_point.Row >= garden_height || adj_point.Col < 0 || adj_point.Col >= garden_width {
			continue
		}
		adj_plot := &garden[adj_point.Row][adj_point.Col]

		// If the adjacent plot doesn't match the current plot type,
		// check the next one.
		if adj_plot.Type != current_plot.Type {
			continue
		}

		// Break down the walls between the two plots.
		modifyWalls(current_plot, adj_plot, dir_str)

		// If the adjacent plot is already in the current region,
		// check the next one.
		if isInSlice(&garden[adj_point.Row][adj_point.Col], *current_region) {
			continue
		}

		// Add the adjacent plot to the current region.
		// Move to the adjacent plot and check its adjacent
		// plots for addition to the region.
		*current_region = append(*current_region, adj_plot)
		spread(garden, current_region, adj_point)
	}
}

func modifyWalls(cur_plot, adj_plot *Plot, dir string) {
	switch dir {
	case "UP":
		cur_plot.Up = false
		adj_plot.Down = false
	case "RIGHT":
		cur_plot.Right = false
		adj_plot.Left = false
	case "DOWN":
		cur_plot.Down = false
		adj_plot.Up = false
	case "LEFT":
		cur_plot.Left = false
		adj_plot.Right = false
	default:
		return
	}
}

func calcFencePricingByPerimeter(region []*Plot) int {
	perimeter := 0

	for _, plot := range region {
		if plot.Up {
			perimeter++
		}
		if plot.Right {
			perimeter++
		}
		if plot.Down {
			perimeter++
		}
		if plot.Left {
			perimeter++
		}
	}

	price := len(region) * perimeter
	return price
}

// Order the region's Plots from top-leftmost to bottom-rightmost. Iterate over all plots
// in this new order. Consider the effect on the number of sides of a region as plots are
// added to it one by one. In the figure below, the "A" Plot is the next plot to be added.
// The cells marked with "*" may or may not hold plots of the current region. The cells
// marked with "_" are guaranteed to not have plots that have already been added since plots
// are added left-to-right and top-to-bottom. This function modifies the number of sides of
// the region by considering each of 16 possible states of the surrounding cells.
//
// +---+---+---+
// | * | * | * |
// +---+---+---+
// | * | A | _ |
// +---+---+---+
// | _ | _ | _ |
// +---+---+---+
func calcFencePricingBySides(region []*Plot) int {
	sides := 0

	sort.Slice(region, func(i, j int) bool {
		if region[i].Pos.Row == region[j].Pos.Row {
			return region[i].Pos.Col < region[j].Pos.Col
		}
		return region[i].Pos.Row < region[j].Pos.Row
	})

	visited_cells := map[Point]bool{}
	for _, current_plot := range region {
		// It is not necessary to check if the adjacent cells are OOB of the original grid.
		// They are simply treated as plots that are not a part of the region.
		// This function does not need to access the original grid.
		up_left_cell := current_plot.Pos.Add(Dir["UP"]).Add(Dir["LEFT"])
		up_cell := current_plot.Pos.Add(Dir["UP"])
		up_right_cell := current_plot.Pos.Add(Dir["UP"]).Add(Dir["RIGHT"])
		left_cell := current_plot.Pos.Add(Dir["LEFT"])

		// Check which adjacent points have been visited
		_, up_left_cell_visited := visited_cells[up_left_cell]
		_, up_cell_visited := visited_cells[up_cell]
		_, up_right_cell_visited := visited_cells[up_right_cell]
		_, left_cell_visited := visited_cells[left_cell]

		if !up_left_cell_visited && !up_cell_visited && !up_right_cell_visited && !left_cell_visited {
			sides += 4 // 0000
		} else if up_left_cell_visited && !up_cell_visited && !up_right_cell_visited && !left_cell_visited {
			sides += 4 // 1000
		} else if !up_left_cell_visited && up_cell_visited && !up_right_cell_visited && !left_cell_visited {
			sides += 0 // 0100
		} else if up_left_cell_visited && up_cell_visited && !up_right_cell_visited && !left_cell_visited {
			sides += 2 // 1100
		} else if !up_left_cell_visited && !up_cell_visited && up_right_cell_visited && !left_cell_visited {
			sides += 4 // 0010
		} else if up_left_cell_visited && !up_cell_visited && up_right_cell_visited && !left_cell_visited {
			sides += 4 // 1010
		} else if !up_left_cell_visited && up_cell_visited && up_right_cell_visited && !left_cell_visited {
			sides += 2 // 0110
		} else if up_left_cell_visited && up_cell_visited && up_right_cell_visited && !left_cell_visited {
			sides += 4 // 1110
		} else if !up_left_cell_visited && !up_cell_visited && !up_right_cell_visited && left_cell_visited {
			sides += 0 // 0001
		} else if up_left_cell_visited && !up_cell_visited && !up_right_cell_visited && left_cell_visited {
			sides += 2 // 1001
		} else if !up_left_cell_visited && up_cell_visited && !up_right_cell_visited && left_cell_visited {
			sides -= 2 // 0101
		} else if up_left_cell_visited && up_cell_visited && !up_right_cell_visited && left_cell_visited {
			sides -= 2 // 1101
		} else if !up_left_cell_visited && !up_cell_visited && up_right_cell_visited && left_cell_visited {
			sides += 0 // 0011
		} else if up_left_cell_visited && !up_cell_visited && up_right_cell_visited && left_cell_visited {
			sides += 2 // 1011
		} else if !up_left_cell_visited && up_cell_visited && up_right_cell_visited && left_cell_visited {
			sides += 0 // 0111
		} else if up_left_cell_visited && up_cell_visited && up_right_cell_visited && left_cell_visited {
			sides += 0 // 1111
		}

		visited_cells[current_plot.Pos] = true
	}

	// fmt.Printf("Type: %s, Area: %d, Sides: %d\n", region[0].Type, len(region), sides)
	return sides * len(region)
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
