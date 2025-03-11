package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Computer struct {
	Halted    bool
	RegisterA int
	RegisterB int
	RegisterC int
	Program   []int
	Pointer   int
}

func (c Computer) getComboOperand(operand int) int {
	switch operand {
	case 0, 1, 2, 3:
		return operand
	case 4:
		return c.RegisterA
	case 5:
		return c.RegisterB
	case 6:
		return c.RegisterC
	}
	return -1
}

func (c *Computer) ADV() {
	exponent := c.getComboOperand(c.Program[c.Pointer+1])
	denominator := 1
	for exponent > 0 && denominator <= c.RegisterA {
		denominator *= 2
		exponent--
	}
	c.RegisterA /= denominator
	c.Pointer += 2
}

func (c *Computer) BXL() {
	c.RegisterB = c.RegisterB ^ c.Program[c.Pointer+1]
	c.Pointer += 2
}

func (c *Computer) BST() {
	c.RegisterB = c.getComboOperand(c.Program[c.Pointer+1]) % 8
	c.Pointer += 2
}

func (c *Computer) JNZ() {
	if c.RegisterA == 0 {
		c.Pointer += 2
	} else {
		c.Pointer = c.Program[c.Pointer+1]
	}
}

func (c *Computer) BXC() {
	c.RegisterB = c.RegisterB ^ c.RegisterC
	c.Pointer += 2
}

func (c *Computer) OUT() int {
	val := c.getComboOperand(c.Program[c.Pointer+1]) % 8
	c.Pointer += 2
	return val
}

func (c *Computer) BDV() {
	exponent := c.getComboOperand(c.Program[c.Pointer+1])
	denominator := 1
	for exponent > 0 && denominator <= c.RegisterA {
		denominator *= 2
		exponent--
	}
	c.RegisterB = c.RegisterA / denominator
	c.Pointer += 2
}

func (c *Computer) CDV() {
	exponent := c.getComboOperand(c.Program[c.Pointer+1])
	denominator := 1
	for exponent > 0 && denominator <= c.RegisterA {
		denominator *= 2
		exponent--
	}
	c.RegisterC = c.RegisterA / denominator
	c.Pointer += 2
}

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	fmt.Printf("Reading computer info...\n")
	computer := Computer{Halted: false, Program: []int{}}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		temp := strings.Split(line, ":")
		if temp[0] == "Register A" {
			num, _ := strconv.Atoi(strings.Trim(temp[1], " "))
			computer.RegisterA = num
		} else if temp[0] == "Register B" {
			num, _ := strconv.Atoi(strings.Trim(temp[1], " "))
			computer.RegisterB = num
		} else if temp[0] == "Register C" {
			num, _ := strconv.Atoi(strings.Trim(temp[1], " "))
			computer.RegisterC = num
		} else if temp[0] == "Program" {
			nums_str := strings.Trim(temp[1], " ")
			nums_slice := strings.Split(nums_str, ",")
			for _, num_str := range nums_slice {
				num, _ := strconv.Atoi(num_str)
				computer.Program = append(computer.Program, num)
			}
		}
	}

	fmt.Printf("Running program...\n")
	output := runProgram(computer)
	fmt.Printf("Output: %v\n", output)
}

func runProgram(computer Computer) []int {
	output := []int{}
	for !computer.Halted {
		instruction := computer.Program[computer.Pointer]

		switch instruction {
		case 0:
			computer.ADV()
		case 1:
			computer.BXL()
		case 2:
			computer.BST()
		case 3:
			computer.JNZ()
		case 4:
			computer.BXC()
		case 5:
			output = append(output, computer.OUT())
		case 6:
			computer.BDV()
		case 7:
			computer.CDV()
		default:
			return []int{}
		}

		if computer.Pointer >= len(computer.Program) {
			computer.Halted = true
		}
	}
	return output
}
