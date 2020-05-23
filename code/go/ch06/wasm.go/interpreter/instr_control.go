package interpreter

import "fmt"

// hack!
func call(vm *vm, args interface{}) {
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
