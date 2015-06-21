package main

import (
	"fmt"
	"log"
	"strconv"
)

// contiguous memory
var main_memory = make([]byte, 2048)

// cache - 2d array composed of 16 blocks (16 bytes per block)
var cache = make([][]byte, 16)

var exit_requested = false

const (
	// BLOCK_OFFSET_MASK uint16 = 0x000F
	// SLOT_OFFSET_MASK  uint16 = 0x00F0
	SLOT    uint16 = 0
	VALID   uint16 = 1
	TAG     uint16 = 2
	DATA    uint16 = 3
	DATA_B2 uint16 = 5
	DATA_B4 uint16 = 7
	DATA_B6 uint16 = 9
	DATA_B8 uint16 = 11
	DATA_BA uint16 = 13
	DATA_BC uint16 = 15
)

func main() {
	// Initialize and Zero cache with 20 byte blocks
	// 16 bytes / data block
	// 1 byte / slot #
	// 1 byte / valid flag
	// 2 bytes / tag
	for i := range cache {
		block := make([]byte, 20)
		block[SLOT] = byte(i)
		cache[i] = block
	}

	max := 256
	mem_val := -1

	// Initialize Main Memory so that Address 0x100 = byte value 00, and so on.
	for i := range main_memory {
		mem_val = (mem_val + 1) % max
		main_memory[i] = byte(mem_val)
	}

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
	fmt.Println("====================================================================================")
	fmt.Println("================================ C A C H E =========================================")
	fmt.Println("Slot  Valid   Tag    Data")

	for entry := range cache {
		row := cache[entry]

		fmt.Printf("%2.1X     %2.1d     %2.1X     %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X  %2.2X\n", row[SLOT], row[VALID], row[TAG], row[3:4], row[4:5], row[5:6], row[6:7], row[7:8], row[8:9], row[9:10], row[10:11], row[11:12], row[12:13], row[13:14], row[14:15], row[15:16], row[16:17], row[17:18], row[18:19])
	}
}

func checkCache(tag byte, slot byte, block byte) bool {
	// TODO: implement complete dump for Cache MISS
	// TODO: implement write to memory
	var cache_line = cache[int(slot)]
	var line_slot = cache_line[SLOT]
	var line_tag = cache_line[TAG]
	var line_valid = cache_line[VALID]
	var line_data = make([]byte, 16)

	fmt.Printf("Cache line contains slot[%X] tag[%X] \n", line_slot, line_tag)

	if line_slot == slot && line_tag == tag && line_valid == 1 {
		copy(line_data[:], cache_line[DATA:])
		fmt.Printf("Cache line COPIED: [%X] \n", line_data)
		return true
	}

	return false
}

func readData() {
	// Get user input
	// TODO: validate input length
	var input string
	fmt.Print("What address would you like to read? ")
	fmt.Scanln(&input)

	// convert to main_memory index
	var idx, err = strconv.ParseInt(input, 16, 16)
	if err != nil {
		log.Fatal(err)
	}

	var addr = uint16(idx)
	var block byte = byte(addr & 0x000F)
	var slot byte = byte((addr & 0x00F0) >> 4)
	var tag byte = byte(addr >> 8)
	fmt.Printf("Tag: [%2.2X]    Slot: [%X]   Block: [%X] \n", tag, slot, block)

	if checkCache(tag, slot, block) {
		// IN THE CACHE
		fmt.Println("Cache HIT")
	} else {
		// NOT IN THE CACHE
		fmt.Println("Cache MISS")
		readAndCacheData(idx, tag, slot, block)
	}
}

func writeData() {
	// Get user input
	// TODO: validate input length

	/////// GET ADDRESS TO WRITE
	var addr_requested string
	fmt.Print("What address would you like to write to? ")
	fmt.Scanln(&addr_requested)

	// convert to main_memory index
	var addr_to_write, err_addr = strconv.ParseInt(addr_requested, 16, 16)
	if err_addr != nil {
		log.Fatal(err_addr)
	}
	fmt.Println(addr_to_write)

	//////// GET DATA TO WRITE
	var data_requested string
	fmt.Print("What data would you like to write at that address? ")
	fmt.Scanln(&data_requested)

	// convert to main_memory index
	var data_to_write, err_data = strconv.ParseInt(data_requested, 16, 16)
	if err_data != nil {
		log.Fatal(err_data)
	}
	fmt.Println(data_to_write)

	// check the cache to see if we already have it
	var addr = uint16(addr_to_write)
	var block byte = byte(addr & 0x000F)
	var slot byte = byte((addr & 0x00F0) >> 4)
	var tag byte = byte(addr >> 8)
	fmt.Printf("Tag: [%2.2X]    Slot: [%X]   Block: [%X] \n", tag, slot, block)

	if checkCache(tag, slot, block) {
		// IN THE CACHE
		fmt.Println("Cache HIT")
		// TODO: already in Cache? write data to cache, and flush to main_memory
	} else {
		// NOT IN THE CACHE
		fmt.Println("Cache MISS")
		// TODO: not in Cache? pull data from main_memory, write data to cache, and flush back to main_memory
	}
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
		writeData()
		break
	}
}

func readAndCacheData(idx int64, tag byte, slot byte, block byte) {
	// NOT IN THE CACHE
	// get the data at that location
	var mem byte = main_memory[idx]
	fmt.Printf("Reading Main_Memory[%d] \n", idx)
	fmt.Printf("Data: [%4.4X] \n", mem)

	fmt.Printf("At that byte there is the value %2.2X \n", mem)

	// read the entire data block
	var addr_begin = idx - int64(block)
	var addr_end = addr_begin + 16
	var data_block = make([]byte, 16)
	data_block = main_memory[addr_begin:addr_end]

	fmt.Printf("Data retrieved: %X \n", data_block[:])

	// write it to the cache
	var slot_idx = int(slot)
	var cache_row = cache[slot_idx]
	cache_row[VALID] = 0x01
	cache_row[TAG] = tag
	copy(cache_row[DATA:], data_block[:])
}
