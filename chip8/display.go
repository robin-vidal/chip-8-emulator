package chip8

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	VM *VM
}

func (g *Game) Update() error {
	for range CyclesPerFrame {
		g.VM.Step()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < ScreenHeight; y++ {
		line := y * ScreenWidth
		for x := 0; x < ScreenWidth; x++ {
			drawColor := color.Black
			if g.VM.display[line+x] {
				drawColor = color.White
			}
			vector.FillRect(screen, float32(x*PixelScale), float32(y*PixelScale), PixelScale, PixelScale, drawColor, false)
		}
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return ScreenWidth * PixelScale, ScreenHeight * PixelScale
}
