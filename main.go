package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/robin-vidal/chip-8-emulator/chip8"
)

const (
	pixelScale     = 10
	cpuHz          = 600
	cyclesPerFrame = cpuHz / 60
)

type game struct {
	vm *chip8.VM
}

func (g *game) Update() error {
	for range cyclesPerFrame {
		g.vm.Step()
	}
	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	for y := range chip8.ScreenHeight {
		for x := range chip8.ScreenWidth {
			drawColor := color.Color(color.Black)
			if g.vm.Pixel(x, y) {
				drawColor = color.White
			}
			vector.FillRect(screen, float32(x*pixelScale), float32(y*pixelScale), pixelScale, pixelScale, drawColor, false)
		}
	}
}

func (g *game) Layout(_, _ int) (int, int) {
	return chip8.ScreenWidth * pixelScale, chip8.ScreenHeight * pixelScale
}

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
	ebiten.SetWindowSize(chip8.ScreenWidth*pixelScale, chip8.ScreenHeight*pixelScale)
	ebiten.SetWindowTitle(*romPathFlag)
	if err := ebiten.RunGame(&game{vm: emulator}); err != nil {
		log.Fatal(err)
	}
}
