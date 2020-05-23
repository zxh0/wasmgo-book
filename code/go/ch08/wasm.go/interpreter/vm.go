package interpreter

import (
	"wasm.go/binary"
)

type vm struct {
	operandStack
	controlStack
	module    binary.Module
	memory    *memory
	globals   []*globalVar
	local0Idx uint32
}

func ExecMainFunc(module binary.Module) {
	vm := &vm{module: module}
	vm.initMem()
	vm.initGlobals()
	call(vm, *module.StartSec)
	vm.loop()
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
func (vm *vm) initGlobals() {
	for _, global := range vm.module.GlobalSec {
		for _, instr := range global.Init {
			vm.execInstr(instr)
		}
		vm.globals = append(vm.globals,
			newGlobal(global.Type, vm.popU64()))
	}
}

/* block stack */

func (vm *vm) enterBlock(opcode byte,
	bt binary.FuncType, instrs []binary.Instruction) {

	bp := vm.stackSize() - len(bt.ParamTypes)
	cf := newControlFrame(opcode, bt, instrs, bp)
	vm.pushControlFrame(cf)
	if opcode == binary.Call {
		vm.local0Idx = uint32(bp)
	}
}
func (vm *vm) exitBlock() {
	cf := vm.popControlFrame()
	vm.clearBlock(cf)
}
func (vm *vm) clearBlock(cf *controlFrame) {
	results := vm.popU64s(len(cf.bt.ResultTypes))
	vm.popU64s(vm.stackSize() - cf.bp)
	vm.pushU64s(results)
	if cf.opcode == binary.Call && vm.controlDepth() > 0 {
		lastCallFrame, _ := vm.topCallFrame()
		vm.local0Idx = uint32(lastCallFrame.bp)
	}
}
func (vm *vm) resetBlock(cf *controlFrame) {
	results := vm.popU64s(len(cf.bt.ParamTypes))
	vm.popU64s(vm.stackSize() - cf.bp)
	vm.pushU64s(results)
}

/* loop */

func (vm *vm) loop() {
	depth := vm.controlDepth()
	for vm.controlDepth() >= depth {
		cf := vm.topControlFrame()
		if cf.pc == len(cf.instrs) {
			vm.exitBlock()
		} else {
			instr := cf.instrs[cf.pc]
			cf.pc++
			vm.execInstr(instr)
		}
	}
}

func (vm *vm) execInstr(instr binary.Instruction) {
	//fmt.Printf("%s %v\n", instr.GetOpname(), instr.Args)
	instrTable[instr.Opcode](vm, instr.Args)
}
