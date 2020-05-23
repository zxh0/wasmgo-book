package interpreter

func localGet(vm *vm, args interface{}) {
	idx := args.(uint32)
	val := vm.getOperand(vm.local0Idx + idx)
	vm.pushU64(val)
}
func localSet(vm *vm, args interface{}) {
	idx := args.(uint32)
	val := vm.popU64()
	vm.setOperand(vm.local0Idx+idx, val)
}
func localTee(vm *vm, args interface{}) {
	idx := args.(uint32)
	val := vm.popU64()
	vm.pushU64(val)
	vm.setOperand(vm.local0Idx+idx, val)
}

func globalGet(vm *vm, args interface{}) {
	idx := args.(uint32)
	val := vm.globals[idx].GetAsU64()
	vm.pushU64(val)
}
func globalSet(vm *vm, args interface{}) {
	idx := args.(uint32)
	val := vm.popU64()
	vm.globals[idx].SetAsU64(val)
}
