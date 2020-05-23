package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
	"wasm.go/binary"
)

func TestOperandStack(t *testing.T) {
	stack := &operandStack{}
	stack.pushBool(true)
	stack.pushBool(false)
	stack.pushU32(1)
	stack.pushS32(-2)
	stack.pushU64(3)
	stack.pushS64(-4)
	stack.pushF32(5.5)
	stack.pushF64(6.5)

	require.Equal(t, 6.5, stack.popF64())
	require.Equal(t, float32(5.5), stack.popF32())
	require.Equal(t, int64(-4), stack.popS64())
	require.Equal(t, uint64(3), stack.popU64())
	require.Equal(t, int32(-2), stack.popS32())
	require.Equal(t, uint32(1), stack.popU32())
	require.Equal(t, false, stack.popBool())
	require.Equal(t, true, stack.popBool())
	require.Equal(t, 0, len(stack.slots))
}

func TestLocalVar(t *testing.T) {
	stack := &operandStack{}
	stack.pushU32(1)
	stack.pushU32(3)
	stack.pushU32(5)

	require.Equal(t, uint64(3), stack.getOperand(1))
	stack.setOperand(1, 7)
	require.Equal(t, uint64(7), stack.getOperand(1))
}

func TestMem(t *testing.T) {
	mem := newMemory(binary.MemType{Min: 1})

	buf := []byte{0x01, 0x02, 0x03}
	mem.Write(10, buf)
	mem.Read(11, buf)
	require.Equal(t, []byte{0x02, 0x03, 0x00}, buf)

	require.Equal(t, uint32(1), mem.Size())
	require.Equal(t, uint32(1), mem.Grow(3))
	require.Equal(t, uint32(4), mem.Size())
}

func TestGlobalVar(t *testing.T) {
	g := newGlobal(binary.GlobalType{
		ValType: binary.ValTypeI32,
		Mut:     1,
	}, 0)
	g.SetAsU64(100)
	require.Equal(t, uint64(100), g.GetAsU64())
}
