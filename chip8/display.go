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
	for range 10 {
		g.VM.Step()
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	for y := 0; y < 32; y++ {
		line := y * 64
		for x := 0; x < 64; x++ {
			drawColor := color.Black
			if g.VM.display[line+x] {
				drawColor = color.White
			}
			vector.FillRect(screen, float32(x*10), float32(y*10), 10, 10, drawColor, false)
		}
	}
}

func (g *Game) Layout(outsideW, outsideH int) (int, int) {
	return 64 * 10, 32 * 10
}
