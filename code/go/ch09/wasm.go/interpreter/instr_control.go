package interpreter

import "wasm.go/binary"

func unreachable(vm *vm, _ interface{}) {
	panic(errTrap)
}

func nop(vm *vm, _ interface{}) {
	// do nothing
}

func block(vm *vm, args interface{}) {
	blockArgs := args.(binary.BlockArgs)
	bt := vm.module.GetBlockType(blockArgs.BT)
	vm.enterBlock(binary.Block, bt, blockArgs.Instrs)
}

func loop(vm *vm, args interface{}) {
	blockArgs := args.(binary.BlockArgs)
	bt := vm.module.GetBlockType(blockArgs.BT)
	vm.enterBlock(binary.Loop, bt, blockArgs.Instrs)
}

func _if(vm *vm, args interface{}) {
	ifArgs := args.(binary.IfArgs)
	bt := vm.module.GetBlockType(ifArgs.BT)
	if vm.popBool() {
		vm.enterBlock(binary.If, bt, ifArgs.Instrs1)
	} else {
		vm.enterBlock(binary.If, bt, ifArgs.Instrs2)
	}
}

func br(vm *vm, args interface{}) {
	labelIdx := int(args.(uint32))
	for i := 0; i < labelIdx; i++ {
		vm.popControlFrame()
	}
	if cf := vm.topControlFrame(); cf.opcode != binary.Loop {
		vm.exitBlock()
	} else {
		vm.resetBlock(cf)
		cf.pc = 0
	}
}

func brIf(vm *vm, args interface{}) {
	if vm.popBool() {
		br(vm, args)
	}
}

func brTable(vm *vm, args interface{}) {
	brTableArgs := args.(binary.BrTableArgs)
	n := int(vm.popU32())
	if n < len(brTableArgs.Labels) {
		br(vm, brTableArgs.Labels[n])
	} else {
		br(vm, brTableArgs.Default)
	}
}

func _return(vm *vm, _ interface{}) {
	_, labelIdx := vm.topCallFrame()
	br(vm, uint32(labelIdx))
}

func call(vm *vm, args interface{}) {
	f := vm.funcs[args.(uint32)]
	callFunc(vm, f)
}

func callFunc(vm *vm, f vmFunc) {
	if f.goFunc != nil {
		callExternalFunc(vm, f)
	} else {
		callInternalFunc(vm, f)
	}
}

func callExternalFunc(vm *vm, f vmFunc) {
	args := popArgs(vm, f._type)
	results := f.goFunc(args)
	pushResults(vm, f._type, results)
}

func popArgs(vm *vm, ft binary.FuncType) []interface{} {
	paramCount := len(ft.ParamTypes)
	args := make([]interface{}, paramCount)
	for i := paramCount - 1; i >= 0; i-- {
		args[i] = wrapU64(ft.ParamTypes[i], vm.popU64())
	}
	return args
}

func pushResults(vm *vm, ft binary.FuncType, results []interface{}) {
	if len(ft.ResultTypes) != len(results) {
		panic("TODO")
	}
	for _, result := range results {
		vm.pushU64(unwrapU64(ft.ResultTypes[0], result))
	}
}

/*
operand stack:

+~~~~~~~~~~~~~~~+
|               |
+---------------+
|     stack     |
+---------------+
|     locals    |
+---------------+
|     params    |
+---------------+
|  ............ |
*/
func callInternalFunc(vm *vm, f vmFunc) {
	vm.enterBlock(binary.Call, f._type, f.code.Expr)

	// alloc locals
	localCount := int(f.code.GetLocalCount())
	for i := 0; i < localCount; i++ {
		vm.pushU64(0)
	}
}

func callIndirect(vm *vm, args interface{}) {
	typeIdx := args.(uint32)
	ft := vm.module.TypeSec[typeIdx]

	i := vm.popU32()
	if i >= vm.table.Size() {
		panic(errUndefinedElem)
	}

	f := vm.table.GetElem(i)
	if f._type.GetSignature() != ft.GetSignature() {
		panic(errTypeMismatch)
	}

	callFunc(vm, f)
}
