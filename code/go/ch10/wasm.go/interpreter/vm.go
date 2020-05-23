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
	funcs     []vmFunc
	table     *table
	local0Idx uint32
}

func ExecMainFunc(module binary.Module) {
	vm := &vm{module: module}
	vm.initMem()
	vm.initGlobals()
	vm.initFuncs()
	vm.initTable()
	if module.StartSec != nil {
		call(vm, *module.StartSec)
	} else {
		call(vm, getMainFuncIdx(module))
	}
	vm.loop()
}

func getMainFuncIdx(module binary.Module) uint32 {
	for _, exp := range module.ExportSec {
		if exp.Desc.Tag == binary.ImportTagFunc && exp.Name == "main" {
			return exp.Desc.Idx
		}
	}
	panic("main func not found")
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
func (vm *vm) initFuncs() {
	vm.linkNativeFuncs()
	for i, ftIdx := range vm.module.FuncSec {
		ft := vm.module.TypeSec[ftIdx]
		code := vm.module.CodeSec[i]
		vm.funcs = append(vm.funcs, newInternalFunc(ft, code))
	}
}
func (vm *vm) linkNativeFuncs() {
	for _, imp := range vm.module.ImportSec {
		if imp.Desc.Tag == binary.ImportTagFunc && imp.Module == "env" {
			ft := vm.module.TypeSec[imp.Desc.FuncType]
			switch imp.Name {
			case "print_char":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, printChar))
			case "assert_true":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertTrue))
			case "assert_false":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertFalse))
			case "assert_eq_i32":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertEqI32))
			case "assert_eq_i64":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertEqI64))
			case "assert_eq_f32":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertEqF32))
			case "assert_eq_f64":
				vm.funcs = append(vm.funcs, newExternalFunc(ft, assertEqF64))
			default:
				panic("TODO")
			}
		}
	}
}
func (vm *vm) initTable() {
	if len(vm.module.TableSec) > 0 {
		vm.table = newTable(vm.module.TableSec[0])
	}
	for _, elem := range vm.module.ElemSec {
		for _, instr := range elem.Offset {
			vm.execInstr(instr)
		}
		offset := vm.popU32()
		for i, funcIdx := range elem.Init {
			vm.table.SetElem(offset+uint32(i), vm.funcs[funcIdx])
		}
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
