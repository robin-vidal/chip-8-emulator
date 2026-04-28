package main

import (
	"flag"
	"fmt"
	"image/color"
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
		if err := g.vm.Step(); err != nil {
			return err
		}
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
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	romPath := flag.String("rom", "", "path to the CHIP-8 ROM file")
	flag.Parse()

	if *romPath == "" {
		return fmt.Errorf("you must specify a ROM path via -rom")
	}

	f, err := os.Open(*romPath)
	if err != nil {
		return err
	}
	defer f.Close()

	emulator := chip8.New()
	if err := emulator.LoadROM(f); err != nil {
		return err
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(chip8.ScreenWidth*pixelScale, chip8.ScreenHeight*pixelScale)
	ebiten.SetWindowTitle(*romPath)
	return ebiten.RunGame(&game{vm: emulator})
}
