package main

import (
	"fmt"

	vm "github.com/robin-vidal/chip-8-emulator/chip8"
)

func main() {
	vm := vm.New()
	fmt.Println("Hello CHIP-8!", vm)
}
