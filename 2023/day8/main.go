package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Println("Parsing documents...")
	mode := "directions"
	var directions []string
	nodes := map[string]map[string]string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		if len(line) == 0 {
			mode = "nodes"
			continue
		}

		if mode == "directions" {
			directions = strings.Split(scanner.Text(), "")
		}

		if mode == "nodes" {
			var current_node string
			next_nodes_map := map[string]string{}

			temp := strings.Split(line, " = ") // {Current_Node, (Left_Node, Right_Node)}
			current_node = temp[0]

			temp2 := strings.TrimLeft(temp[1], "(")
			temp2 = strings.TrimRight(temp2, ")")
			next_nodes := strings.Split(temp2, ", ")
			next_nodes_map["L"] = next_nodes[0]
			next_nodes_map["R"] = next_nodes[1]
			nodes[current_node] = next_nodes_map
		}
	}

	fmt.Println("Pathing from AAA to ZZZ...")
	num_steps := getNumSteps(nodes, directions, "AAA", "ZZZ")
	num_steps2 := getNumSteps2(nodes, directions)
	fmt.Printf("[Part 1] Number of steps: %d\n", num_steps)
	fmt.Printf("[Part 2] Number of steps: %d\n", num_steps2)
}

func getNumSteps(nodes map[string]map[string]string, directions []string, start, end string) int {
	num_steps := 0
	current_node := start
	for {
		if current_node == end {
			break
		}

		direction := directions[num_steps%len(directions)]
		current_node = nodes[current_node][direction]

		num_steps += 1
	}
	return num_steps
}

// Find all nodes that end with A.
// Each node "__A" is assigned a single node "__Z".
// Keep track of how many steps it takes each starting node to reach its "__Z".
// It takes the same number of steps to go from "__Z" to "__A".
// Find and return the LCM of these step counts.
func getNumSteps2(nodes map[string]map[string]string, directions []string) int {
	num_steps := 0
	current_nodes := []string{}

	for node := range nodes {
		if node[len(node)-1] == 'A' {
			current_nodes = append(current_nodes, node)
		}
	}
	fmt.Println(current_nodes)
	loop_steps := make([]int, len(current_nodes))

	for {
		done := true

		for i, node := range current_nodes {
			if node[len(node)-1] == 'Z' && loop_steps[i] == 0 {
				loop_steps[i] = num_steps
			}
		}

		for _, steps := range loop_steps {
			if steps == 0 {
				done = false
				break
			}
		}

		if done {
			break
		}

		direction := directions[num_steps%len(directions)]
		for i := 0; i < len(current_nodes); i++ {
			current_nodes[i] = nodes[current_nodes[i]][direction]
		}

		num_steps += 1
	}

	return lcmOfSlice(loop_steps)
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return abs(a*b) / gcd(a, b)
}

func lcmOfSlice(nums []int) int {
	if len(nums) == 0 {
		return 0
	}

	result := nums[0]
	for _, n := range nums[1:] {
		result = lcm(result, n)
	}
	return result
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
