package main

import (
	"flag"
	"fmt"
	"image/color"
	"os"
	"path/filepath"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/robin-vidal/chip-8-emulator/chip8"
)

const (
	pixelScale     = 10
	cpuHz          = 600
	cyclesPerFrame = cpuHz / 60
)

var keyMapping = [16]ebiten.Key{
	ebiten.KeyX,
	ebiten.Key1,
	ebiten.Key2,
	ebiten.Key3,
	ebiten.KeyQ,
	ebiten.KeyW,
	ebiten.KeyE,
	ebiten.KeyA,
	ebiten.KeyS,
	ebiten.KeyD,
	ebiten.KeyZ,
	ebiten.KeyC,
	ebiten.Key4,
	ebiten.KeyR,
	ebiten.KeyF,
	ebiten.KeyV,
} // Key-mapping

type game struct {
	vm *chip8.VM
}

func (g *game) Update() error {
	for idx, key := range keyMapping {
		g.vm.Keys[idx] = ebiten.IsKeyPressed(key)
	}

	for range cyclesPerFrame {
		if err := g.vm.Step(); err != nil {
			return err
		}
	}

	g.vm.DecrementTimers()

	return nil
}

func (g *game) Draw(screen *ebiten.Image) {
	for y := range chip8.ScreenHeight {
		for x := range chip8.ScreenWidth {
			vector.FillRect(screen, float32(x*pixelScale), float32(y*pixelScale), pixelScale, pixelScale, pixelColor(g.vm.Pixel(x, y)), false)
		}
	}
}

func pixelColor(lit bool) color.Color {
	if lit {
		return color.White
	}
	return color.Black
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
	shiftInPlace := flag.Bool("shift-in-place", false, "shift VX directly instead of copying VY first (CHIP-48)")
	jumpOffsetVX := flag.Bool("jump-offset-vx", false, "BNNN jumps to XNN+VX instead of NNN+V0 (CHIP-48)")
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
	emulator.ShiftInPlace = *shiftInPlace
	emulator.JumpOffsetVX = *jumpOffsetVX
	if err := emulator.LoadROM(f); err != nil {
		return err
	}

	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowSize(chip8.ScreenWidth*pixelScale, chip8.ScreenHeight*pixelScale)
	ebiten.SetWindowTitle(filepath.Base(*romPath))
	return ebiten.RunGame(&game{vm: emulator})
}
