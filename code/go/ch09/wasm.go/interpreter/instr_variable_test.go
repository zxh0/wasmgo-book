package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"wasm.go/binary"
)

func TestLocal(t *testing.T) {
	vm := &vm{operandStack: operandStack{slots: []uint64{123, 456, 789}}}
	vm.local0Idx = 1

	instrTable[binary.LocalGet](vm, uint32(1))
	require.Equal(t, vm.operandStack.slots[2], vm.popU64())

	vm.pushU64(246)
	instrTable[binary.LocalTee](vm, uint32(1))
	require.Equal(t, vm.operandStack.slots[3], vm.operandStack.slots[2])
	instrTable[binary.LocalSet](vm, uint32(0))
	require.Equal(t, vm.operandStack.slots[2], vm.operandStack.slots[1])
}

func TestGlobal(t *testing.T) {
	vm := &vm{globals: []*globalVar{
		newGlobal(binary.GlobalType{Mut: 1}, 100),
		newGlobal(binary.GlobalType{Mut: 1}, 200),
		newGlobal(binary.GlobalType{Mut: 1}, 300),
	}}

	instrTable[binary.GlobalGet](vm, uint32(0)) // [100
	instrTable[binary.GlobalGet](vm, uint32(1)) // [100[200
	instrTable[binary.GlobalGet](vm, uint32(2)) // [100[200[300
	instrTable[binary.GlobalSet](vm, uint32(1))
	instrTable[binary.GlobalSet](vm, uint32(0))
	instrTable[binary.GlobalSet](vm, uint32(2))

	require.Equal(t, uint64(200), vm.globals[0].GetAsU64())
	require.Equal(t, uint64(300), vm.globals[1].GetAsU64())
	require.Equal(t, uint64(100), vm.globals[2].GetAsU64())
}
