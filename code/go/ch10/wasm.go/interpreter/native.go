package interpreter

import "fmt"

func printChar(args []interface{}) []interface{} {
	fmt.Printf("%c", args[0].(int32))
	return nil
}

func assertTrue(args []interface{}) []interface{} {
	assertEq(args[0].(int32), int32(1))
	return nil
}
func assertFalse(args []interface{}) []interface{} {
	assertEq(args[0].(int32), int32(0))
	return nil
}

func assertEqI32(args []interface{}) []interface{} {
	assertEq(args[0].(int32), args[1].(int32))
	return nil
}
func assertEqI64(args []interface{}) []interface{} {
	assertEq(args[0].(int64), args[1].(int64))
	return nil
}
func assertEqF32(args []interface{}) []interface{} {
	assertEq(args[0].(float32), args[1].(float32))
	return nil
}
func assertEqF64(args []interface{}) []interface{} {
	assertEq(args[0].(float64), args[1].(float64))
	return nil
}

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}
