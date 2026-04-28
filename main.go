package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/robin-vidal/chip-8-emulator/chip8"
)

func main() {
	var romPathFlag = flag.String("romPath", "", "Path to the ROM you want to load")
	flag.Parse()

	if *romPathFlag == "" {
		fmt.Fprintln(os.Stderr, "you must specify a ROM")
		os.Exit(1)
	}

	emulator := chip8.New()

	err := emulator.LoadROM(*romPathFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(chip8.ScreenWidth*chip8.PixelScale, chip8.ScreenHeight*chip8.PixelScale)
	ebiten.SetWindowTitle(*romPathFlag)
	if err := ebiten.RunGame(&chip8.Game{VM: emulator}); err != nil {
		log.Fatal(err)
	}
}
