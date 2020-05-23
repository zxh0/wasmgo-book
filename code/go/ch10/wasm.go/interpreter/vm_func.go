package interpreter

import (
	"fmt"

	"wasm.go/binary"
	"wasm.go/instance"
)

type vmFunc struct {
	vm    *vm
	_type binary.FuncType
	code  binary.Code
	_func instance.Function
}

func newExternalFunc(ft binary.FuncType,
	f instance.Function) vmFunc {

	return vmFunc{
		_type: ft,
		_func: f,
	}
}
func newInternalFunc(vm *vm, ft binary.FuncType,
	code binary.Code) vmFunc {

	return vmFunc{
		vm:    vm,
		_type: ft,
		code:  code,
	}
}

func (f vmFunc) Type() binary.FuncType {
	return f._type
}
func (f vmFunc) Call(args ...WasmVal) ([]WasmVal, error) {
	if f._func != nil {
		return f._func.Call(args...)
	}
	return f.safeCall(args)
}

func (f vmFunc) safeCall(args []WasmVal) (results []WasmVal, err error) {
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

	results = f.call(args)
	return
}

func (f vmFunc) call(args []interface{}) []interface{} {
	pushArgs(f.vm, f._type, args)
	callFunc(f.vm, f)
	if f._func == nil {
		f.vm.loop()
	}
	return popResults(f.vm, f._type)
}

func pushArgs(vm *vm, ft binary.FuncType, args []interface{}) {
	if len(ft.ParamTypes) != len(args) {
		panic(fmt.Errorf("param count: %d, arg count: %d",
			len(ft.ParamTypes), len(args)))
	}
	for i, vt := range ft.ParamTypes {
		vm.pushU64(unwrapU64(vt, args[i]))
	}
}
func popResults(vm *vm, ft binary.FuncType) []interface{} {
	results := make([]interface{}, len(ft.ResultTypes))
	for n := len(ft.ResultTypes) - 1; n >= 0; n-- {
		results[n] = wrapU64(ft.ResultTypes[n], vm.popU64())
	}
	return results
}
