package main

import (
	"fmt"
)

// contiguous memory
var main_memory = make([]byte, 2048)

// cache - 2d array composed of 16 blocks (16 bytes per block)
var cache = make([][]byte, 16)

func main() {
	main_memory[0] = 14

	fmt.Printf(">>> %d \n", main_memory[0])
	fmt.Printf("+++ %d \n", main_memory[1])

	// Zero cache, 16 bytes per block
	for i := range cache {
		block := make([]byte, 16)
		cache[i] = block
	}

	max := 256
	mem_val := -1

	// Initialize Main Memory
	for i := range main_memory {
		mem_val = (mem_val + 1) % max
		main_memory[i] = byte(mem_val)

		fmt.Printf(">> %X \n", main_memory[i])
	}
}
