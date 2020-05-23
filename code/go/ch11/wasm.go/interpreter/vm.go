package interpreter

import (
	"fmt"

	"wasm.go/binary"
	"wasm.go/instance"
	"wasm.go/validator"
)

type WasmVal = instance.WasmVal

type vm struct {
	operandStack
	controlStack
	module    binary.Module
	memory    instance.Memory
	table     instance.Table
	globals   []instance.Global
	funcs     []vmFunc
	local0Idx uint32
}

func New(m binary.Module, mm map[string]instance.Module,
) (inst instance.Module, err error) {

	if err := validator.Validate(m); err != nil {
		return nil, err
	}

	defer func() {
		if _err := recover(); _err != nil {
			switch x := _err.(type) {
			case error:
				err = x
			default:
				panic(err)
			}
		}
	}()

	inst = newVM(m, mm)
	return
}

func newVM(m binary.Module, mm map[string]instance.Module) *vm {
	vm := &vm{module: m}
	vm.linkImports(mm)
	vm.initFuncs()
	vm.initTable()
	vm.initMem()
	vm.initGlobals()
	vm.execStartFunc()
	return vm
}

/* linking */

func (vm *vm) linkImports(mm map[string]instance.Module) {
	for _, imp := range vm.module.ImportSec {
		if m := mm[imp.Module]; m == nil {
			panic(fmt.Errorf("module not found: " + imp.Module))
		} else {
			vm.linkImport(m, imp)
		}
	}
}
func (vm *vm) linkImport(m instance.Module, imp binary.Import) {
	exported := m.GetMember(imp.Name)
	if exported == nil {
		panic(fmt.Errorf("unknown import: %s.%s",
			imp.Module, imp.Name))
	}

	typeMatched := false
	switch x := exported.(type) {
	case instance.Function:
		if imp.Desc.Tag == binary.ImportTagFunc {
			expectedFT := vm.module.TypeSec[imp.Desc.FuncType]
			typeMatched = isFuncTypeMatch(expectedFT, x.Type())
			vm.funcs = append(vm.funcs, newExternalFunc(expectedFT, x))
		}
	case instance.Table:
		if imp.Desc.Tag == binary.ImportTagTable {
			typeMatched = isLimitsMatch(imp.Desc.Table.Limits, x.Type().Limits)
			vm.table = x
		}
	case instance.Memory:
		if imp.Desc.Tag == binary.ImportTagMem {
			typeMatched = isLimitsMatch(imp.Desc.Mem, x.Type())
			vm.memory = x
		}
	case instance.Global:
		if imp.Desc.Tag == binary.ImportTagGlobal {
			typeMatched = isGlobalTypeMatch(imp.Desc.Global, x.Type())
			vm.globals = append(vm.globals, x)
		}
	}

	if !typeMatched {
		panic(fmt.Errorf("incompatible import type: %s.%s",
			imp.Module, imp.Name))
	}
}

/* init */

func (vm *vm) initMem() {
	if len(vm.module.MemSec) > 0 {
		vm.memory = newMemory(vm.module.MemSec[0])
	}
	for _, data := range vm.module.DataSec {
		vm.execConstExpr(data.Offset)
		offset := vm.popU64() // TODO: check offset
		vm.memory.Write(offset, data.Init)
	}
}
func (vm *vm) initGlobals() {
	for _, global := range vm.module.GlobalSec {
		vm.execConstExpr(global.Init)
		vm.globals = append(vm.globals,
			newGlobal(global.Type, vm.popU64()))
	}
}
func (vm *vm) initFuncs() {
	for i, ftIdx := range vm.module.FuncSec {
		ft := vm.module.TypeSec[ftIdx]
		code := vm.module.CodeSec[i]
		vm.funcs = append(vm.funcs, newInternalFunc(vm, ft, code))
	}
}
func (vm *vm) initTable() {
	if len(vm.module.TableSec) > 0 {
		vm.table = newTable(vm.module.TableSec[0])
	}
	for _, elem := range vm.module.ElemSec {
		vm.execConstExpr(elem.Offset)
		offset := vm.popU32() // TODO: check offset
		for i, funcIdx := range elem.Init {
			vm.table.SetElem(offset+uint32(i), vm.funcs[funcIdx])
		}
	}
}

func (vm *vm) execConstExpr(expr []binary.Instruction) {
	for _, instr := range expr {
		vm.execInstr(instr)
	}
}
func (vm *vm) execStartFunc() {
	if vm.module.StartSec != nil {
		idx := *vm.module.StartSec
		vm.funcs[idx].call(nil)
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

/* instance.Module */

func (vm *vm) GetMember(name string) interface{} {
	for _, exp := range vm.module.ExportSec {
		if exp.Name == name {
			idx := exp.Desc.Idx
			switch exp.Desc.Tag {
			case binary.ExportTagFunc:
				return vm.funcs[idx]
			case binary.ExportTagTable:
				return vm.table
			case binary.ExportTagMem:
				return vm.memory
			case binary.ExportTagGlobal:
				return vm.globals[idx]
			}
		}
	}
	return nil
}

func (vm *vm) InvokeFunc(name string, args ...WasmVal) ([]WasmVal, error) {
	m := vm.GetMember(name)
	if m != nil {
		if f, ok := m.(instance.Function); ok {
			return f.Call(args...)
		}
	}
	return nil, fmt.Errorf("function not found: " + name)
}
func (vm vm) GetGlobalVal(name string) (WasmVal, error) {
	m := vm.GetMember(name)
	if m != nil {
		if g, ok := m.(instance.Global); ok {
			return g.Get(), nil
		}
	}
	return nil, fmt.Errorf("global not found: " + name)
}
func (vm vm) SetGlobalVal(name string, val WasmVal) error {
	m := vm.GetMember(name)
	if m != nil {
		if g, ok := m.(instance.Global); ok {
			g.Set(val)
			return nil
		}
	}
	return fmt.Errorf("global not found: " + name)
}

/* helpers */

func isFuncTypeMatch(expected, actual binary.FuncType) bool {
	return fmt.Sprintf("%s", expected) == fmt.Sprintf("%s", actual)
}
func isGlobalTypeMatch(expected, actual binary.GlobalType) bool {
	return actual.ValType == expected.ValType &&
		actual.Mut == expected.Mut
}
func isLimitsMatch(expected, actual binary.Limits) bool {
	return actual.Min >= expected.Min &&
		(expected.Max == 0 || actual.Max > 0 && actual.Max <= expected.Max)
}
