package binary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeVarUint(t *testing.T) {
	data := []byte{
		0b1_0111111,
		0b1_0011111,
		0b1_0001111,
		0b1_0000111,
		0b1_0000011,
		0b0_0000001}
	testDecodeVarUint32(t, data[5:], 0b0000001, 1)
	testDecodeVarUint32(t, data[4:], 0b1_0000011, 2)
	testDecodeVarUint32(t, data[3:], 0b1_0000011_0000111, 3)
	testDecodeVarUint32(t, data[2:], 0b1_0000011_0000111_0001111, 4)
	testDecodeVarUint32(t, data[1:], 0b1_0000011_0000111_0001111_0011111, 5)
	//testDecodeVarUint32(t, data[0:], 0, 0)
}

func TestDecodeVarInt(t *testing.T) {
	data := []byte{0xC0, 0xBB, 0x78}
	testDecodeVarInt32(t, data, int32(-123456), 3)
}

func testDecodeVarUint32(t *testing.T, data []byte, n uint32, w int) {
	_n, _w := decodeVarUint(data, 32)
	require.Equal(t, n, uint32(_n))
	require.Equal(t, w, _w)
}
func testDecodeVarInt32(t *testing.T, data []byte, n int32, w int) {
	_n, _w := decodeVarInt(data, 32)
	require.Equal(t, n, int32(_n))
	require.Equal(t, w, _w)
}
