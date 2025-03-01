package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

func main() {
	filename := "input.txt"
	file, err := os.Open(filename)
	if err != nil {
		fmt.Printf("Could not open %s: %v\n", filename, err)
		return
	}
	defer file.Close()

	is_file := true
	disk := []int{}
	id := 0
	reader := bufio.NewReader(file)

	// Read the diskmap string character by character
	fmt.Printf("Parsing disk map...\n")
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			fmt.Printf("Error reading file: %v\n", err)
			break
		}

		// Alternate between creating file blocks and free spaces
		if unicode.IsDigit(r) {
			size := int(r - '0')
			if is_file {
				for range size {
					disk = append(disk, id)
				}
				id++
				is_file = false
			} else {
				for range size {
					disk = append(disk, -1) // -1 will represent a free space since valid ids are not negative
				}
				is_file = true
			}
		}
	}

	fmt.Printf("Moving file blocks...\n")
	moveBlocks(disk)

	fmt.Printf("Calculating checksum...\n")
	checksum := getChecksum(disk)

	fmt.Printf("Filesystem checksum: %d\n", checksum)
}

func moveBlocks(disk []int) {
	disk_space := len(disk)
	left_space_index := 0
	right_block_index := disk_space - 1

	space_found := false
	block_found := false

	for {
		// Find the next leftmost free space
		for i := left_space_index; i < disk_space; i++ {
			if disk[i] == -1 {
				left_space_index = i
				space_found = true
				break
			}
		}

		// Find the rightmost file block
		for j := right_block_index; j > 0; j-- {
			if disk[j] != -1 {
				right_block_index = j
				block_found = true
				break
			}
		}

		// Edge cases: break if no file blocks exist or no free spaces exist
		if !space_found || !block_found {
			break
		}
		// End the loop once there are no gaps remaining
		if left_space_index >= right_block_index {
			break
		}

		// Move file block to the left and free its original space
		temp := disk[left_space_index]
		disk[left_space_index] = disk[right_block_index]
		disk[right_block_index] = temp

		space_found = false
		block_found = false
	}
}

func getChecksum(disk []int) int {
	checksum := 0
	for i := 0; i < len(disk); i++ {
		if disk[i] == -1 {
			break
		}
		checksum += disk[i] * i
	}
	return checksum
}
