package chip8

const (
	OpMisc     = 0x0 // Children: OpClear, OpReturn (not implemented)
	OpClear    = 0xE0
	OpJump     = 0x1
	OpSet      = 0x6
	OpAdd      = 0x7
	OpSetIndex = 0xA
)

type instruction struct {
	kind    uint8
	x, y, n uint8
	nn      uint8
	nnn     uint16
}

func (vm *VM) Step() {
	opcode := vm.fetch()
	instruction := decode(opcode)
	vm.execute(instruction)
}

func (vm *VM) fetch() uint16 {
	res := uint16(vm.memory[vm.PC])<<8 | uint16(vm.memory[vm.PC+1])
	vm.PC += 2
	return res
}

func decode(opcode uint16) *instruction {
	instruction := new(instruction)
	instruction.kind = uint8(opcode >> 12)
	instruction.x = uint8((opcode >> 8) & 0x000F)
	instruction.y = uint8((opcode >> 4) & 0x000F)
	instruction.n = uint8((opcode) & 0x000F)
	instruction.nn = uint8((opcode) & 0x00FF)
	instruction.nnn = uint16((opcode) & 0x0FFF)

	return instruction
}

func (vm *VM) execute(instruction *instruction) {
	switch instruction.kind {
	case OpMisc:
		switch instruction.nn {
		case OpClear:
			for i := range vm.display {
				vm.display[i] = false
			}
		}
	case OpJump:
		vm.PC = instruction.nnn
	case OpSet:
		vm.V[instruction.x] = instruction.nn
	case OpAdd:
		vm.V[instruction.x] += instruction.nn
	case OpSetIndex:
		vm.I = instruction.nnn
	}
}
