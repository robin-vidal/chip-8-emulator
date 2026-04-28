package chip8

const (
	OpMisc     = 0x0 // Children: OpClear, OpReturn (not implemented)
	OpClear    = 0xE0
	OpJump     = 0x1
	OpSet      = 0x6
	OpAdd      = 0x7
	OpSetIndex = 0xA
	OpDisplay  = 0xD
)

type instruction struct {
	kind    uint8
	x, y, n uint8
	nn      uint8
	nnn     uint16
}

func (vm *VM) Step() {
	opcode := vm.fetch()
	instr := decode(opcode)
	vm.execute(instr)
}

func (vm *VM) fetch() uint16 {
	res := uint16(vm.memory[vm.PC])<<8 | uint16(vm.memory[vm.PC+1])
	vm.PC += 2
	return res
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

func (vm *VM) execute(instr instruction) {
	switch instr.kind {
	case OpMisc:
		switch instr.nn {
		case OpClear:
			for i := range vm.display {
				vm.display[i] = false
			}
		}
	case OpJump:
		vm.PC = instr.nnn
	case OpSet:
		vm.V[instr.x] = instr.nn
	case OpAdd:
		vm.V[instr.x] += instr.nn
	case OpSetIndex:
		vm.I = instr.nnn
	case OpDisplay:
		x, y := vm.V[instr.x]%64, vm.V[instr.y]%32
		vm.V[0xF] = 0

		for line := range instr.n {
			octet := vm.memory[vm.I+uint16(line)]
			for n := range 8 {
				shouldToggle := (octet>>(7-n))&1 == 1
				if shouldToggle {
					idx := (uint16(y)+uint16(line))*64 + uint16(x) + uint16(n)
					if vm.display[idx] {
						vm.V[0xF] = 1
					}

					vm.display[idx] = !vm.display[idx]
				}
			}
		}
	}
}
