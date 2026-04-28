package chip8

// Hardware specification constants.
const (
	ScreenWidth  = 64
	ScreenHeight = 32
)

// Rendering constants. Will move to main once the Game adapter is extracted.
const (
	PixelScale     = 10
	cpuHz          = 600
	CyclesPerFrame = cpuHz / 60
)
