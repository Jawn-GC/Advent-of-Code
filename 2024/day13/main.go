package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

type ClawMachine struct {
	ButtonA       []string
	ButtonB       []string
	PrizeLocation []string
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
	}
	defer file.Close()

	fmt.Printf("Reading claw machine info...\n")
	claw_machines := []ClawMachine{}
	claw_machine := ClawMachine{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// Create a new claw machine struct if an empty
		// line is encountered in the file.
		if len(line) == 0 {
			claw_machine = ClawMachine{}
		}

		// Get the numeric values from the claw machine info.
		temp := strings.Split(line, ":")
		if temp[0] == "Button A" {
			temp = strings.Split(temp[1], ",")
			x := strings.TrimLeft(temp[0], " X+")
			y := strings.TrimLeft(temp[1], " Y+")
			claw_machine.ButtonA = []string{x, y}
		} else if temp[0] == "Button B" {
			temp = strings.Split(temp[1], ",")
			x := strings.TrimLeft(temp[0], " X+")
			y := strings.TrimLeft(temp[1], " Y+")
			claw_machine.ButtonB = []string{x, y}
		} else if temp[0] == "Prize" {
			temp = strings.Split(temp[1], ",")
			x := strings.TrimLeft(temp[0], " X=")
			y := strings.TrimLeft(temp[1], " Y=")
			claw_machine.PrizeLocation = []string{x, y}
			claw_machines = append(claw_machines, claw_machine)
		}
	}

	fmt.Printf("Calculating number of tokens required...\n")
	total_tokens1 := 0
	total_tokens2 := 0
	for _, claw_machine := range claw_machines {
		total_tokens1 += numTokensRequired(claw_machine, 1)
		total_tokens2 += numTokensRequired(claw_machine, 2)
	}
	fmt.Printf("[Part 1] Total number of tokens required: %d\n", total_tokens1)
	fmt.Printf("[Part 2] Total number of tokens required: %d\n", total_tokens2)
}

func numTokensRequired(machine ClawMachine, mode int) int {
	tokens := 0
	var tol float64
	var matrix [][]float64
	if mode == 1 {
		tol = 0.0000001
		matrix = getMachineMatrix1(machine)
	} else if mode == 2 {
		tol = 0.01
		matrix = getMachineMatrix2(machine)
	}
	rref(matrix)

	if isValidSolution(matrix, tol) {
		a := math.Round(matrix[0][2]) // Number of presses for Button A
		b := math.Round(matrix[1][2]) // Number of presses for Button B

		tokens += 3*int(a) + int(b)
	}

	return tokens
}

// The matrix will be 2x3.
// Column 1 represents movement due to pushing Button A.
// Column 2 represents movement due to pushing Button B.
// Column 3 represents the goal position.
func getMachineMatrix1(machine ClawMachine) [][]float64 {
	matrix := [][]float64{}

	// Convert strings to ints and store them in rows of
	// a matrix. Ignore potential conversion errors.
	xA, _ := strconv.Atoi(machine.ButtonA[0])
	xB, _ := strconv.Atoi(machine.ButtonB[0])
	xP, _ := strconv.Atoi(machine.PrizeLocation[0])
	rowX := []float64{float64(xA), float64(xB), float64(xP)}
	matrix = append(matrix, rowX)

	yA, _ := strconv.Atoi(machine.ButtonA[1])
	yB, _ := strconv.Atoi(machine.ButtonB[1])
	yP, _ := strconv.Atoi(machine.PrizeLocation[1])
	rowY := []float64{float64(yA), float64(yB), float64(yP)}
	matrix = append(matrix, rowY)

	return matrix
}

func getMachineMatrix2(machine ClawMachine) [][]float64 {
	matrix := [][]float64{}

	// Convert strings to ints and store them in rows of
	// a matrix. Ignore potential conversion errors.
	xA, _ := strconv.Atoi(machine.ButtonA[0])
	xB, _ := strconv.Atoi(machine.ButtonB[0])
	xP, _ := strconv.Atoi(machine.PrizeLocation[0])
	rowX := []float64{float64(xA), float64(xB), float64(10000000000000 + xP)}
	matrix = append(matrix, rowX)

	yA, _ := strconv.Atoi(machine.ButtonA[1])
	yB, _ := strconv.Atoi(machine.ButtonB[1])
	yP, _ := strconv.Atoi(machine.PrizeLocation[1])
	rowY := []float64{float64(yA), float64(yB), float64(10000000000000 + yP)}
	matrix = append(matrix, rowY)

	return matrix
}

// Row Reduced Echelon Form (n x m matrix, n < m)
// It is assumed that the initial entries of the matrix
// are all positive values.
func rref(matrix [][]float64) {
	num_rows := len(matrix)
	num_cols := len(matrix[0])

	for i := 0; i < num_rows; i++ {
		pivot := matrix[i][i]
		if pivot != 0 {
			for m := 0; m < num_cols; m++ {
				matrix[i][m] /= pivot
			}
		}

		for n := 0; n < num_rows; n++ {
			if n != i {
				factor := matrix[n][i]
				for m := 0; m < num_cols; m++ {
					matrix[n][m] -= factor * matrix[i][m]
				}
			}
		}
	}
}

// The matrix is assumed to be in rref with a unique solution.
// This function simply checks if the solution vector has
// whole-number values.
func isValidSolution(matrix [][]float64, tol float64) bool {
	a := matrix[0][2]
	b := matrix[1][2]

	// The number of button presses should not be negative.
	if a < 0 || b < 0 {
		return false
	}

	ar := math.Round(a)
	br := math.Round(b)

	if math.Abs(ar-a) > tol || math.Abs(br-b) > tol {
		return false
	}

	return true
}
