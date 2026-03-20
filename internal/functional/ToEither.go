// Package functional provides utility functions for working with functional programming constructs.
package functional

import (
	E "github.com/IBM/fp-go/v2/either"
	IOE "github.com/IBM/fp-go/v2/ioeither"
)

// ToEither converts an IOEither to an Either by executing the IO action and returning the result as an Either.
func ToEither[ERR, A any](ioe IOE.IOEither[ERR, A]) E.Either[ERR, A] {
	return ioe()
}
