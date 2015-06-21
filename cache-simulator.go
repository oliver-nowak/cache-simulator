package main

import (
	"fmt"
	"log"
	"strconv"
)

// contiguous memory
var main_memory = make([]uint16, 2048)

// cache - 2d array composed of 16 blocks (16 bytes per block)
var cache = make([][]uint16, 16)

var exit_requested = false

const (
	BLOCK_OFFSET_MASK uint16 = 0x000F
	SLOT_OFFSET_MASK  uint16 = 0x00F0 // >> 1
	// TAG_OFFSET_MASK   uint16 = 0xFF00 // >> 2
	SLOT  uint16 = 0
	VALID uint16 = 1
	TAG   uint16 = 2
	DATA  uint16 = 3
)

func main() {
	// Initialize and Zero cache with 20 byte blocks
	// 16 bytes / data block
	// 1 byte / slot #
	// 1 byte / valid flag
	// 2 bytes / tag
	for i := range cache {
		block := make([]uint16, 20)
		cache[i] = block
	}

	max := 256
	mem_val := -1

	// Initialize Main Memory so that Address 0x100 = byte value 00, and so on.
	for i := range main_memory {
		mem_val = (mem_val + 1) % max
		main_memory[i] = uint16(mem_val)

		// fmt.Printf(">> %X \n", main_memory[i])
	}

	// assert 0x7FF = 0xFF
	// fmt.Printf("++ %X \n", main_memory[0x07FF])

	// displayCache()

	// User Input loop
	for !exit_requested {
		displayMenu()
		getMenuInput()
	}
}

func displayMenu() {
	fmt.Println("")
	fmt.Println("")
	fmt.Printf("(R)ead, (W)rite, (D)isplay Cache, (Q)uit? ")
}

func displayCache() {
	fmt.Println("==========================================================================")
	fmt.Println("================================C A C H E ================================")
	fmt.Println("Slot  Valid   Tag    Data")

	for entry := range cache {
		row := cache[entry]

		fmt.Printf("%2.1X     %2.1d     %2.1X     %2.2X \n", row[SLOT], row[VALID], row[TAG], row[DATA:])
	}
}

func readData() {
	// Get user input
	var input string
	fmt.Print("What address would you like to read? ")
	fmt.Scanln(&input)

	// convert to main_memory index
	var idx, err = strconv.ParseInt(input, 16, 16)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Reading Main_Memory[%d] \n", idx)

	// get the data at that location
	var mem uint16 = main_memory[idx]
	fmt.Printf("Data: [%4.4X] \n", mem)

	var block uint16 = mem & BLOCK_OFFSET_MASK
	var slot uint16 = (mem & SLOT_OFFSET_MASK) >> 4
	var tag uint16 = mem >> 8

	fmt.Printf("Tag: [%2.2X]    Slot: [%X]   Block: [%X] \n", tag, slot, block)

	// check if its in the cache

	// read the entire data block if its not in the cache

	fmt.Printf("At that byte there is the value %2.2X \n", mem)
}

func getMenuInput() {
	var input string
	fmt.Scanf("%s", &input)

	switch input {
	case "d", "D":
		displayCache()
		break
	case "r", "R":
		readData()
		break
	case "q", "Q":
		exit_requested = true
		break
	case "w", "W":
		break
	}
}
