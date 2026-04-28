package chip8

import (
	"fmt"
	"io"
)

const (
	memorySize   = 4096
	programStart = 0x200
	fontStart    = 0x50
)

type VM struct {
	memory  [memorySize]uint8
	display [ScreenWidth * ScreenHeight]bool
	pc      uint16
	i       uint16
	stack   []uint16
	delay   uint8
	sound   uint8
	v       [16]uint8
}

func New() *VM {
	vm := new(VM)
	vm.pc = programStart
	copy(vm.memory[fontStart:], font[:])
	return vm
}

func (vm *VM) Pixel(x, y int) bool {
	return vm.display[y*ScreenWidth+x]
}

func (vm *VM) LoadROM(r io.Reader) error {
	data, err := io.ReadAll(r)
	if err != nil {
		return err
	}

	maxSize := memorySize - programStart
	if len(data) > maxSize {
		return fmt.Errorf("rom size %d exceeds maximum %d bytes", len(data), maxSize)
	}

	copy(vm.memory[programStart:], data)

	return nil
}
