package main

import (
	"bufio"
	"fmt"
	"os"
	"unicode"
)

type Block struct {
	Index int
	Id    int
	Size  int
}

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
	file_blocks := []Block{}
	free_blocks := []Block{}
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

		// Alternate between creating file blocks and free blocks
		if unicode.IsDigit(r) {
			size := int(r - '0')
			if is_file {
				block := Block{Index: len(disk), Id: id, Size: size}
				file_blocks = append(file_blocks, block)
				for range size {
					disk = append(disk, id)
				}
				id++
				is_file = false
			} else {
				block := Block{Index: len(disk), Id: -1, Size: size}
				free_blocks = append(free_blocks, block)
				for range size {
					disk = append(disk, -1) // -1 will represent a free space since valid ids are not negative
				}
				is_file = true
			}
		}
	}

	fmt.Printf("Moving file blocks...\n")
	disk_copy1 := copyDisk(disk)
	disk_copy2 := copyDisk(disk)

	moveBlocks1(disk_copy1)
	moveBlocks2(disk_copy2, file_blocks, free_blocks)

	fmt.Printf("Calculating checksums...\n")
	checksum1 := getChecksum(disk_copy1)
	checksum2 := getChecksum(disk_copy2)

	fmt.Printf("[Part 1] Filesystem Checksum: %d\n", checksum1)
	fmt.Printf("[Part 2] Filesystem Checksum: %d\n", checksum2)
}

func moveBlocks1(disk []int) {
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
		for j := right_block_index; j >= 0; j-- {
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

func moveBlocks2(disk []int, files []Block, gaps []Block) {
	for i := len(files) - 1; i >= 0; i-- { // Iterate over the file blocks in order of decreasing id #
	fileLoop:
		for j := 0; j < len(gaps); j++ { // Iterate over the gaps starting from the left
			if gaps[j].Size != 0 && gaps[j].Index < files[i].Index && gaps[j].Size >= files[i].Size {
				for k := 0; k < files[i].Size; k++ { // Fill the jth gap with the ith file
					disk[gaps[j].Index+k] = files[i].Id
					disk[files[i].Index+k] = -1
				}
				gaps[j].Size -= files[i].Size
				gaps[j].Index += files[i].Size

				break fileLoop
			}
		}
	}
}

func getChecksum(disk []int) int {
	checksum := 0
	for i := 0; i < len(disk); i++ {
		id := disk[i]
		if id != -1 {
			checksum += id * i
		}
	}
	return checksum
}

func copyDisk(disk []int) []int {
	disk_copy := []int{}
	disk_copy = append(disk_copy, disk...)
	return disk_copy
}
