package binary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecode(t *testing.T) {
	module, err := DecodeFile("./testdata/hw_rust.wasm")
	require.NoError(t, err)
	require.Equal(t, uint32(MagicNumber), module.Magic)
	require.Equal(t, uint32(Version), module.Version)
	require.Equal(t, 2, len(module.CustomSecs))
	require.Equal(t, 15, len(module.TypeSec))
	require.Equal(t, 0, len(module.ImportSec))
	require.Equal(t, 171, len(module.FuncSec))
	require.Equal(t, 1, len(module.TableSec))
	require.Equal(t, 1, len(module.MemSec))
	require.Equal(t, 4, len(module.GlobalSec))
	require.Equal(t, 5, len(module.ExportSec))
	require.Nil(t, module.StartSec)
	require.Equal(t, 1, len(module.ElemSec))
	require.Equal(t, 171, len(module.CodeSec))
	require.Equal(t, 4, len(module.DataSec))
}
