package chip8

func (vm *VM) Step() {
	vm.fetch()
}

func (vm *VM) fetch() uint16 {
	res := uint16(vm.memory[vm.PC])<<8 | uint16(vm.memory[vm.PC])
	vm.PC += 2
	return res
}
