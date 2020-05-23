package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"wasm.go/binary"
)

func TestMemSizeAndGrow(t *testing.T) {
	vm := &vm{memory: newMemory(binary.MemType{Min: 2})}
	instrTable[binary.MemorySize](vm, nil)
	require.Equal(t, uint64(2), vm.popU64())

	vm.pushU32(3)
	instrTable[binary.MemoryGrow](vm, nil)
	require.Equal(t, uint64(2), vm.popU64())

	instrTable[binary.MemorySize](vm, nil)
	require.Equal(t, uint64(5), vm.popU64())
}

func TestMemOps(t *testing.T) {
	vm := &vm{memory: newMemory(binary.MemType{Min: 1})}
	testMemOp(t, vm, binary.I32Store, binary.I32Load, 0x10, 0x01, int32(100))
	testMemOp(t, vm, binary.I64Store, binary.I64Load, 0x20, 0x02, int64(123))
	testMemOp(t, vm, binary.F32Store, binary.F32Load, 0x30, 0x03, float32(1.5))
	testMemOp(t, vm, binary.F64Store, binary.F64Load, 0x40, 0x04, 1.5)
	testMemOp(t, vm, binary.I32Store8, binary.I32Load8S, 0x50, 0x05, int32(-100))
	testMemOp(t, vm, binary.I32Store8, binary.I32Load8U, 0x60, 0x06, int32(100))
	testMemOp(t, vm, binary.I32Store16, binary.I32Load16S, 0x70, 0x07, int32(-10000))
	testMemOp(t, vm, binary.I32Store16, binary.I32Load16U, 0x80, 0x08, int32(10000))
	testMemOp(t, vm, binary.I64Store8, binary.I64Load8S, 0x90, 0x09, int32(-100))
	testMemOp(t, vm, binary.I64Store8, binary.I64Load8U, 0xA0, 0x0A, int32(100))
	testMemOp(t, vm, binary.I64Store16, binary.I64Load16S, 0xB0, 0x0B, int32(-10000))
	testMemOp(t, vm, binary.I64Store16, binary.I64Load16U, 0xC0, 0x0C, int32(10000))
	testMemOp(t, vm, binary.I64Store32, binary.I64Load32S, 0xD0, 0x0D, int32(-1000000))
	testMemOp(t, vm, binary.I64Store32, binary.I64Load32U, 0xE0, 0x0E, int32(1000000))
}

func testMemOp(t *testing.T, vm *vm, storeOp, loadOp byte,
	offset, i uint32, val interface{}) {

	memArg := binary.MemArg{Offset: offset}

	// store
	vm.pushU32(i)
	pushVal(vm, val)
	instrTable[storeOp](vm, memArg)

	// load
	vm.pushU32(i)
	instrTable[loadOp](vm, memArg)

	// check
	require.Equal(t, val, popVal(vm, val))
}
