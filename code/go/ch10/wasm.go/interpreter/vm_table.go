package interpreter

import (
	"wasm.go/binary"
)

type table struct {
	_type binary.TableType
	elems []vmFunc
}

func newTable(tt binary.TableType) *table {
	return &table{
		_type: tt,
		elems: make([]vmFunc, tt.Limits.Min),
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
	t.elems = append(t.elems, make([]vmFunc, n)...)
}

func (t *table) GetElem(idx uint32) vmFunc {
	t.checkIdx(idx)
	elem := t.elems[idx]
	return elem
}
func (t *table) SetElem(idx uint32, elem vmFunc) {
	t.checkIdx(idx)
	t.elems[idx] = elem
}

func (t *table) checkIdx(idx uint32) {
	if idx >= uint32(len(t.elems)) {
		panic(errUndefinedElem)
	}
}
