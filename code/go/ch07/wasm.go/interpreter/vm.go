package interpreter

import (
	"wasm.go/binary"
)

type vm struct {
	operandStack
	module binary.Module
	memory *memory
}

func ExecMainFunc(module binary.Module) {
	idx := int(*module.StartSec) - len(module.ImportSec)
	vm := &vm{module: module}
	vm.initMem()
	vm.execCode(idx)
}

func (vm *vm) initMem() {
	if len(vm.module.MemSec) > 0 {
		vm.memory = newMemory(vm.module.MemSec[0])
	}
	for _, data := range vm.module.DataSec {
		for _, instr := range data.Offset {
			vm.execInstr(instr)
		}
		vm.memory.Write(vm.popU64(), data.Init)
	}
}

func (vm *vm) execCode(idx int) {
	code := vm.module.CodeSec[idx]
	for _, instr := range code.Expr {
		vm.execInstr(instr)
	}
}

func (vm *vm) execInstr(instr binary.Instruction) {
	//fmt.Printf("%s %v\n", instr.GetOpname(), instr.Args)
	instrTable[instr.Opcode](vm, instr.Args)
}
