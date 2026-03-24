package file

import (
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	FILE "github.com/IBM/fp-go/v2/ioeither/file"
	"io"
	"os"
)

type IOEither[E, A any] = IOE.IOEither[E, A]

var createTemp = IOE.TryCatchError(func() (*os.File, error) {
	return os.CreateTemp("", "jbrowse-web-*.zip")
})

// createTempWriter accepts a Kleisli and ccomposes it with the createTemp Kleisli. Basically, if you have a function that takes a WriteCloser and returns an IOEither of a file,
// you can use createTempWriter to create a temporary file and use that as the argument to your function.
var createTempWriter = F.Pipe1(createTemp, FILE.Write[*os.File])

// copyToFile is a curried function that copies data from a reader to a file
var copyToFile = F.Curry2(
	func(rc io.ReadCloser, file *os.File) IOEither[error, *os.File] {
		return IOE.TryCatchError(func() (*os.File, error) {
			_, err := io.Copy(file, rc)
			return file, err
		})
	},
)

var IOEcopyToFile = IOE.Eitherize2(
	func(rc io.ReadCloser, file *os.File) (*os.File, error) {
		_, err := io.Copy(file, rc)
		return file, err
	},
)

// SaveTemp takes a ReadCloser and saves it to a temporary file, returning a pointer to the file.
var SaveTemp = F.Flow2(copyToFile, createTempWriter)
