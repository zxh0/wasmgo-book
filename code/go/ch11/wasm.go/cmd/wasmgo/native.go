package main

import (
	"fmt"

	"wasm.go/instance"
)

func newEnv() instance.Module {
	env := instance.NewNativeInstance()
	env.RegisterFunc("print_char(i32)->()", printChar)
	env.RegisterFunc("assert_true(i32)->()", assertTrue)
	env.RegisterFunc("assert_false(i32)->()", assertFalse)
	env.RegisterFunc("assert_eq_i32(i32,i32)->()", assertEqI32)
	env.RegisterFunc("assert_eq_i64(i64,i64)->()", assertEqI64)
	env.RegisterFunc("assert_eq_f32(f32,f32)->()", assertEqF32)
	env.RegisterFunc("assert_eq_f64(f64,f64)->()", assertEqF64)
	return env
}

func printChar(args []interface{}) ([]interface{}, error) {
	fmt.Printf("%c", args[0].(int32))
	return nil, nil
}

func assertTrue(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(int32), int32(1))
	return nil, nil
}
func assertFalse(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(int32), int32(0))
	return nil, nil
}

func assertEqI32(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(int32), args[1].(int32))
	return nil, nil
}
func assertEqI64(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(int64), args[1].(int64))
	return nil, nil
}
func assertEqF32(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(float32), args[1].(float32))
	return nil, nil
}
func assertEqF64(args []interface{}) ([]interface{}, error) {
	assertEq(args[0].(float64), args[1].(float64))
	return nil, nil
}

func assertEq(a, b interface{}) {
	if a != b {
		panic(fmt.Errorf("%v != %v", a, b))
	}
}
