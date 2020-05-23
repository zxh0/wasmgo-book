package interpreter

import "errors"

var (
	errTrap              = errors.New("unreachable")
	errCallStackOverflow = errors.New("call stack exhausted")
	errTypeMismatch      = errors.New("indirect call type mismatch")
	errUndefinedElem     = errors.New("undefined element")
	errUninitializedElem = errors.New("uninitialized element")
	errMemOutOfBounds    = errors.New("out of bounds memory access")
	errImmutableGlobal   = errors.New("immutable global")
	errIntOverflow       = errors.New("integer overflow")
	errConvertToInt      = errors.New("invalid conversion to integer")
)
