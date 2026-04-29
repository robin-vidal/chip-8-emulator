package chip8

import "fmt"

const vf = 0xF // VF flag register

const (
	OpSys              = 0x0 // Children: OpClear, OpReturn (not implemented)
	OpClear            = 0xE0
	OpReturnSubroutine = 0xEE
	OpJump             = 0x1
	OpCallSubroutine   = 0x2
	OpSet              = 0x6
	OpAdd              = 0x7
	OpSetIndex         = 0xA
	OpDisplay          = 0xD
)

type instruction struct {
	kind    uint8
	x, y, n uint8
	nn      uint8
	nnn     uint16
}

func (vm *VM) Step() error {
	opcode := vm.fetch()
	instr := decode(opcode)
	return vm.execute(instr)
}

func (vm *VM) fetch() uint16 {
	opcode := uint16(vm.memory[vm.pc])<<8 | uint16(vm.memory[vm.pc+1])
	vm.pc += 2
	return opcode
}

func decode(opcode uint16) instruction {
	instr := instruction{}
	instr.kind = uint8(opcode >> 12)
	instr.x = uint8((opcode >> 8) & 0x000F)
	instr.y = uint8((opcode >> 4) & 0x000F)
	instr.n = uint8((opcode) & 0x000F)
	instr.nn = uint8((opcode) & 0x00FF)
	instr.nnn = uint16((opcode) & 0x0FFF)

	return instr
}

func (vm *VM) execute(instr instruction) error {
	switch instr.kind {
	case OpSys:
		switch instr.nn {
		case OpClear:
			vm.clearDisplay()
		case OpReturnSubroutine:
			popped := vm.stack[len(vm.stack)-1]
			vm.stack = vm.stack[:len(vm.stack)-1]
			vm.jump(popped)
		}
	case OpJump:
		vm.jump(instr.nnn)
	case OpCallSubroutine:
		vm.stack = append(vm.stack, vm.pc)
		vm.jump(instr.nnn)
	case OpSet:
		vm.v[instr.x] = instr.nn
	case OpAdd:
		vm.v[instr.x] += instr.nn
	case OpSetIndex:
		vm.i = instr.nnn
	case OpDisplay:
		vm.executeDisplay(instr)
	default:
		return fmt.Errorf("unknown opcode: 0x%X", instr.kind)
	}
	return nil
}

func (vm *VM) jump(addr uint16) {
	vm.pc = addr
}

func (vm *VM) clearDisplay() {
	vm.display = [ScreenWidth * ScreenHeight]bool{}
}

func (vm *VM) executeDisplay(instr instruction) {
	x, y := vm.v[instr.x]%ScreenWidth, vm.v[instr.y]%ScreenHeight
	vm.v[vf] = 0

	for line := range instr.n {
		py := uint16(y) + uint16(line)
		if py >= ScreenHeight {
			break // clip vertically
		}
		octet := vm.memory[vm.i+uint16(line)]
		for n := range 8 {
			px := uint16(x) + uint16(n)
			if px >= ScreenWidth {
				break // clip horizontally
			}
			shouldToggle := (octet>>(7-n))&1 == 1
			if shouldToggle {
				idx := py*ScreenWidth + px
				if vm.display[idx] {
					vm.v[vf] = 1
				}
				vm.display[idx] = !vm.display[idx]
			}
		}
	}
}
