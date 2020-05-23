package instance

import (
	"wasm.go/binary"
)

var _ Function = (*nativeFunction)(nil)

type GoFunc = func(args []WasmVal) ([]WasmVal, error)

type nativeFunction struct {
	t binary.FuncType
	f GoFunc
}

func (nf nativeFunction) Type() binary.FuncType {
	return nf.t
}
func (nf nativeFunction) Call(args ...WasmVal) ([]WasmVal, error) {
	return nf.f(args)
}
