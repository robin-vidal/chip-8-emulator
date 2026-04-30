package chip8

import "fmt"

const (
	v0 = 0x0
	vf = 0xF // VF flag register
)

const (
	OpSet             = 0x6
	OpAdd             = 0x7
	OpSetIndex        = 0xA
	OpDisplay         = 0xD
	OpJump            = 0x1
	OpJumpOffset      = 0xB
	OpCallSubroutine  = 0x2
	OpSkipEqualXNN    = 0x3
	OpSkipNotEqualXNN = 0x4
	OpSkipEqualXY     = 0x5
	OpSkipNotEqualXY  = 0x9

	OpSys              = 0x0 // Parent (NN vary)
	OpClear            = 0xE0
	OpReturnSubroutine = 0xEE

	OpLogicalArithmetic = 0x8 // Parent (N vary)
	OpSetXY             = 0x0
	OpBinaryOr          = 0x1
	OpBinaryAnd         = 0x2
	OpLogicalXor        = 0x3
	OpAddXY             = 0x4
	OpSubstractXY       = 0x5
	OpSubstractYX       = 0x7
	OpShiftRight        = 0x6
	OpShiftLeft         = 0xE
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
	case OpSet:
		vm.v[instr.x] = instr.nn
	case OpAdd:
		vm.v[instr.x] += instr.nn
	case OpSetIndex:
		vm.i = instr.nnn
	case OpDisplay:
		vm.executeDisplay(instr)
	case OpJump:
		vm.jump(instr.nnn)
	case OpJumpOffset:
		vm.jump(instr.nnn + uint16(vm.v[v0]))
	case OpCallSubroutine:
		vm.stack = append(vm.stack, vm.pc)
		vm.jump(instr.nnn)
	case OpSkipEqualXNN:
		if vm.v[instr.x] == instr.nn {
			vm.pc += 2
		}
	case OpSkipNotEqualXNN:
		if vm.v[instr.x] != instr.nn {
			vm.pc += 2
		}
	case OpSkipEqualXY:
		if vm.v[instr.x] == vm.v[instr.y] {
			vm.pc += 2
		}
	case OpSkipNotEqualXY:
		if vm.v[instr.x] != vm.v[instr.y] {
			vm.pc += 2
		}
	case OpSys:
		switch instr.nn {
		case OpClear:
			vm.clearDisplay()
		case OpReturnSubroutine:
			popped := vm.stack[len(vm.stack)-1]
			vm.stack = vm.stack[:len(vm.stack)-1]
			vm.jump(popped)
		}
	case OpLogicalArithmetic:
		switch instr.n {
		case OpSetXY:
			vm.v[instr.x] = vm.v[instr.y]
		case OpBinaryOr:
			vm.v[instr.x] |= vm.v[instr.y]
		case OpBinaryAnd:
			vm.v[instr.x] &= vm.v[instr.y]
		case OpLogicalXor:
			vm.v[instr.x] ^= vm.v[instr.y]
		case OpAddXY:
			res := uint16(vm.v[instr.x]) + uint16(vm.v[instr.y])
			vm.v[vf] = uint8(res >> 8)
			vm.v[instr.x] = uint8(res)
		case OpSubstractXY:
			borrow := uint8(0)
			if vm.v[instr.x] < vm.v[instr.y] {
				borrow = 1
			}
			vm.v[instr.x] -= vm.v[instr.y]
			vm.v[vf] = 1 - borrow
		case OpSubstractYX:
			borrow := uint8(0)
			if vm.v[instr.x] > vm.v[instr.y] {
				borrow = 1
			}
			vm.v[instr.x] = vm.v[instr.y] - vm.v[instr.x]
			vm.v[vf] = 1 - borrow
		case OpShiftRight:
			vm.v[instr.x] = vm.v[instr.y]
			vm.v[vf] = vm.v[instr.x] & 1
			vm.v[instr.x] >>= 1
		case OpShiftLeft:
			vm.v[instr.x] = vm.v[instr.y]
			vm.v[vf] = (vm.v[instr.x] >> 7) & 1
			vm.v[instr.x] <<= 1
		}
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
