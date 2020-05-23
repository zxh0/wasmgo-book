package interpreter

import (
	"testing"

	"github.com/stretchr/testify/require"
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
