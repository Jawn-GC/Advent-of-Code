package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// The directional fields determine if the plot
// has a wall on a particular side.
type Plot struct {
	Type  string
	Up    bool
	Right bool
	Down  bool
	Left  bool
}

func NewPlot(plot_type string) Plot {
	return Plot{
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
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		plots := strings.Split(line, "")
		garden_row := []Plot{}
		for _, p := range plots {
			garden_row = append(garden_row, NewPlot(p))
		}
		garden = append(garden, garden_row)
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
	fmt.Printf("The total price of fencing (by perimeter) for all regions is %d\n", total_price_perimeter)
	fmt.Printf("The total price of fencing (by sides) for all regions is %d\n", total_price_sides)
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

func calcFencePricingBySides(region []*Plot) int {
	price := 0

	return price
}

func isInSlice[T comparable](val T, slice []T) bool {
	for _, v := range slice {
		if v == val {
			return true
		}
	}
	return false
}
