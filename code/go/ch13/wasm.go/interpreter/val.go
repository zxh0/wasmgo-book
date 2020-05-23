package interpreter

import (
	"math"

	"wasm.go/binary"
)

func wrapU64(vt binary.ValType, val uint64) interface{} {
	switch vt {
	case binary.ValTypeI32:
		return int32(val)
	case binary.ValTypeI64:
		return int64(val)
	case binary.ValTypeF32:
		return math.Float32frombits(uint32(val))
	case binary.ValTypeF64:
		return math.Float64frombits(val)
	default:
		panic("unreachable") // TODO
	}
}

func unwrapU64(vt binary.ValType, val interface{}) uint64 {
	switch vt {
	case binary.ValTypeI32:
		return uint64(val.(int32))
	case binary.ValTypeI64:
		return uint64(val.(int64))
	case binary.ValTypeF32:
		return uint64(math.Float32bits(val.(float32)))
	case binary.ValTypeF64:
		return math.Float64bits(val.(float64))
	default:
		panic("unreachable") // TODO
	}
}
