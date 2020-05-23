package interpreter

import (
	"wasm.go/binary"
	"wasm.go/instance"
)

type table struct {
	_type binary.TableType
	elems []instance.Function
}

func NewTable(min, max uint32) instance.Table {
	tt := binary.TableType{
		ElemType: binary.FuncRef,
		Limits:   binary.Limits{Min: min, Max: max},
	}
	return newTable(tt)
}

func newTable(tt binary.TableType) *table {
	return &table{
		_type: tt,
		elems: make([]instance.Function, tt.Limits.Min),
	}
}

func (t *table) Type() binary.TableType {
	return t._type
}

func (t *table) Size() uint32 {
	return uint32(len(t.elems))
}
func (t *table) Grow(n uint32) {
	// TODO: check max
	t.elems = append(t.elems, make([]instance.Function, n)...)
}

func (t *table) GetElem(idx uint32) instance.Function {
	t.checkIdx(idx)
	elem := t.elems[idx]
	if elem == nil {
		panic(errUninitializedElem)
	}
	return elem
}
func (t *table) SetElem(idx uint32, elem instance.Function) {
	t.checkIdx(idx)
	t.elems[idx] = elem
}

func (t *table) checkIdx(idx uint32) {
	if idx >= uint32(len(t.elems)) {
		panic(errUndefinedElem)
	}
}
