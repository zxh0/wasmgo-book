package binary

import "errors"

var (
	errUnexpectedEnd = errors.New("unexpected end of section or function")
	errIntTooLong    = errors.New("integer representation too long")
	errIntTooLarge   = errors.New("integer too large")
	//errLenOutOfBounds = errors.New("length out of bounds")
)
