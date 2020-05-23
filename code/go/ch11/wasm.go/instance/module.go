package instance

import "wasm.go/binary"

type WasmVal = interface{}
type Map = map[string]Module

type Module interface {
	GetMember(name string) interface{}
	InvokeFunc(name string, args ...WasmVal) ([]WasmVal, error)
	GetGlobalVal(name string) (WasmVal, error)
	SetGlobalVal(name string, val WasmVal) error
}

type Function interface {
	Type() binary.FuncType
	Call(args ...WasmVal) ([]WasmVal, error)
}

type Table interface {
	Type() binary.TableType
	Size() uint32
	Grow(n uint32)
	GetElem(idx uint32) Function
	SetElem(idx uint32, elem Function)
}

type Memory interface {
	Type() binary.MemType
	Size() uint32 // page count
	Grow(n uint32) uint32
	Read(offset uint64, buf []byte)
	Write(offset uint64, buf []byte)
}

type Global interface {
	Type() binary.GlobalType
	GetAsU64() uint64
	SetAsU64(val uint64)
	Get() WasmVal
	Set(val WasmVal)
}
