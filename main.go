package main

import (
	"flag"
	"fmt"
	"os"

	vm "github.com/robin-vidal/chip-8-emulator/chip8"
)

func main() {
	var romPathFlag = flag.String("romPath", "", "Path to the ROM you want to load")
	flag.Parse()

	if *romPathFlag == "" {
		fmt.Fprintln(os.Stderr, "you must specify a ROM")
		os.Exit(1)
	}

	emulator := vm.New()

	err := emulator.LoadROM(*romPathFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Println("ROM", *romPathFlag, "loaded!")
}
