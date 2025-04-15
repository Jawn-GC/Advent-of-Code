package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"
)

type GardenMap struct {
	Destination int
	Source      int
	Range       int
}

type RangeInfo struct {
	SourceStart int
	SourceRange int
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Reading almanac...")
	mode := "seeds"
	seeds := []int{}
	garden_maps := [][]GardenMap{}
	garden_maps_row := []GardenMap{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Swap modes
		r, _ := utf8.DecodeLastRuneInString(line)
		if r == ':' {
			if mode != "seeds" {
				garden_maps = append(garden_maps, garden_maps_row)
			}
			mode = "maps"
			garden_maps_row = []GardenMap{}
			continue
		}

		// Parse integers
		if mode == "seeds" && strings.TrimSpace(line) != "" {
			temp := strings.TrimPrefix(line, "seeds: ")
			temp2 := strings.Fields(temp)

			for _, str := range temp2 {
				num, _ := strconv.Atoi(str)
				seeds = append(seeds, num)
			}
		} else if mode == "maps" && strings.TrimSpace(line) != "" {
			temp := strings.Fields(line) // {Destination, Source, Range}
			des, _ := strconv.Atoi(temp[0])
			src, _ := strconv.Atoi(temp[1])
			rng, _ := strconv.Atoi(temp[2])
			garden_maps_row = append(garden_maps_row, GardenMap{Destination: des, Source: src, Range: rng})
		}

	}
	garden_maps = append(garden_maps, garden_maps_row) // Add the final row of GardenMaps

	fmt.Println("Adding missing ranges to the almanac...")
	garden_maps = getAllGardenMaps(garden_maps)

	fmt.Println("Getting seed locations...")
	destinations := []int{}
	for _, seed := range seeds {
		destinations = append(destinations, getFinalDestination(seed, garden_maps))
	}

	// The "seeds" slice is assumed to have an even number of values
	//location_ranges := []RangeInfo{}
	//for i := 0; i < len(seeds); i += 2 {
	//	initial_range_info := RangeInfo{SourceStart: seeds[i], SourceRange: seeds[i+1]}
	//	ranges := []RangeInfo{}
	//	ranges = append(ranges, initial_range_info)

	//	for j, _ := range garden_maps {
	//		next_ranges := []RangeInfo{}

	// Previous source ranges are split into new ranges as determined by the ranges of the GardenMaps
	//		for _, current_range := range ranges {
	//			next_ranges = append(next_ranges, getNextRanges(current_range, garden_maps, j)...)
	//		}

	//		ranges = next_ranges
	//	}

	//	location_ranges = append(location_ranges, ranges...)
	//}

	sort.Ints(destinations)
	fmt.Printf("[Part 1] The nearest seed has a location value of %d\n", destinations[0])

	//sort.Slice(location_ranges, func(i, j int) bool {
	//	return location_ranges[i].SourceStart < location_ranges[j].SourceStart
	//})
	//fmt.Printf("[Part 2 (WIP)] The nearest seed has a location value of %d\n", location_ranges[0].SourceStart)
}

// Row 0: seed-to-soil
// Row 1: soil-to-fertilizer
// Row 2: fertilizer-to-water
// Row 3: water-to-light
// Row 4: light-to-tempertature
// Row 5: temperature-to-humidity
// Row 6: humidity-to-location
func getFinalDestination(seed int, garden_maps [][]GardenMap) int {
	current_val := seed

	for _, map_row := range garden_maps {
		for _, garden_map := range map_row {
			difference := current_val - garden_map.Source
			if difference >= 0 && difference < garden_map.Range {
				current_val = garden_map.Destination + difference
				break
			}
		}
	}

	return current_val
}

func getNextRanges(current_range RangeInfo, garden_maps [][]GardenMap, current_row int) []RangeInfo {
	current_maps := garden_maps[current_row]
	thresholds := []int{}
	next_ranges := []RangeInfo{}
	source_splice := []RangeInfo{}

	// Any source values that exceed the maps are mapped to themselves
	last_map := current_maps[len(current_maps)-1]
	last_source := last_map.Source + last_map.Range
	current_maps = append(current_maps, GardenMap{Destination: last_source, Source: last_source, Range: math.MaxInt})

	for i, m := range current_maps {
		if m.Source != 0 {
			thresholds = append(thresholds, m.Source)
		}

		// Add the last threshold value
		if i == len(current_maps)-1 {
			thresholds = append(thresholds, m.Source+m.Range)
		}
	}
	sort.Ints(thresholds)

	// Split the source range
	current_value := current_range.SourceStart
	final_value := current_range.SourceStart + current_range.SourceRange
	for i, t := range thresholds {
		if current_value < t {
			rng := t - current_value
			source_splice = append(source_splice, RangeInfo{SourceStart: current_value, SourceRange: rng})
			current_value = t
		}

		if i == len(thresholds)-1 && final_value >= t {
			source_splice = append(source_splice, RangeInfo{SourceStart: current_value, SourceRange: final_value - t + 1})
		}
	}

	// Find the destination ranges
	for _, r := range source_splice {
		for _, garden_map := range current_maps {
			difference := r.SourceStart - garden_map.Source
			if difference >= 0 && difference < garden_map.Range {
				next_ranges = append(next_ranges, RangeInfo{SourceStart: garden_map.Destination + difference, SourceRange: r.SourceRange})
				break
			}
		}
	}
	return next_ranges
}

// There may be gaps in the source ranges for a particular map type.
// If a source value is within a gap, its destination value is just equal to its source value.
// This function adds these gaps to the map collection.
// All source and destination values are assumed to be nonnegative.
func getAllGardenMaps(garden_maps [][]GardenMap) [][]GardenMap {
	complete_garden_maps := [][]GardenMap{}
	for _, row := range garden_maps {
		new_row := []GardenMap{}
		new_row = append(new_row, row...)
		sort.Slice(new_row, func(i, j int) bool { // Sort by ascensding source values
			return new_row[i].Source < new_row[j].Source
		})

		current_source_value := 0
		n := len(row)
		final_source := row[n-1].Source + row[n-1].Range - 1 // Final original source value
		new_maps := []GardenMap{}

		// Find where the gaps are and create new GardenMaps for them.
		for {
			for _, m := range new_row {
				if current_source_value > m.Source {
					continue
				} else if current_source_value == m.Source {
					current_source_value += m.Range
					break
				} else if current_source_value < m.Source {
					rng := m.Source - current_source_value
					new_maps = append(new_maps, GardenMap{Destination: current_source_value, Source: current_source_value, Range: rng})
					current_source_value = m.Source
					break
				}
			}

			if current_source_value > final_source {
				break
			}
		}

		new_row = append(new_row, new_maps...)
		sort.Slice(new_row, func(i, j int) bool { // Sort by ascending source values after adding the gaps
			return new_row[i].Source < new_row[j].Source
		})
		complete_garden_maps = append(complete_garden_maps, new_row)
	}

	return complete_garden_maps
}
