package interpreter

import (
	"fmt"

	"wasm.go/binary"
)

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
	idx := int(args.(uint32))
	importedFuncCount := len(vm.module.ImportSec) // TODO
	if idx < importedFuncCount {
		callAssertFunc(vm, args) // hack!
	} else {
		callInternalFunc(vm, idx-importedFuncCount)
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
func callInternalFunc(vm *vm, idx int) {
	ftIdx := vm.module.FuncSec[idx]
	ft := vm.module.TypeSec[ftIdx]
	code := vm.module.CodeSec[idx]
	vm.enterBlock(binary.Call, ft, code.Expr)

	// alloc locals
	localCount := int(code.GetLocalCount())
	for i := 0; i < localCount; i++ {
		vm.pushU64(0)
	}
}

// hack!
func callAssertFunc(vm *vm, args interface{}) {
	idx := args.(uint32)
	switch vm.module.ImportSec[idx].Name {
	case "assert_true":
		assertEq(vm.popBool(), true)
	case "assert_false":
		assertEq(vm.popBool(), false)
	case "assert_eq_i32":
		assertEq(vm.popU32(), vm.popU32())
	case "assert_eq_i64":
		assertEq(vm.popU64(), vm.popU64())
	case "assert_eq_f32":
		assertEq(vm.popF32(), vm.popF32())
	case "assert_eq_f64":
		assertEq(vm.popF64(), vm.popF64())
	default:
		panic("TODO")
	}
}

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}
