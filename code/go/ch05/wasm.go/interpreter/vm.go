package interpreter

import (
	"wasm.go/binary"
)

type vm struct {
	operandStack
	module binary.Module
}

func ExecMainFunc(module binary.Module) {
	idx := int(*module.StartSec) - len(module.ImportSec)
	vm := &vm{module: module}
	vm.execCode(idx)
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
