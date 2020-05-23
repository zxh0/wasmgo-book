package binary

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSignature(t *testing.T) {
	ft := FuncType{
		ParamTypes:  []ValType{ValTypeI32, ValTypeI64, ValTypeF32, ValTypeF64},
		ResultTypes: nil,
	}
	require.Equal(t, "(i32,i64,f32,f64)->()", ft.GetSignature())

	ft = FuncType{
		ParamTypes:  nil,
		ResultTypes: []ValType{ValTypeI64},
	}
	require.Equal(t, "()->(i64)", ft.GetSignature())
}
