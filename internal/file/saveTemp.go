// 1. stream ReadCloser to temp zip file
// 2. Open zip file to get ReadCloser
// 3. Create destination file with zip file's path and name.
// 4. copy from reader to writer
// 5. done for 1 file
// 6 map over all files
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

var createTempWriter = F.Pipe1(createTemp, FILE.Write[int64, *os.File])

// copyToFile is a curried function that copies data from a reader to a file
var copyToFile = F.Curry2(
	func(body io.ReadCloser, file *os.File) IOEither[error, int64] {
		return IOE.TryCatchError(func() (int64, error) {
			return io.Copy(file, body)
		})
	},
)

// 1. stream ReadCloser to temp file
func saveTemp(rc io.ReadCloser) IOEither[error, int64] {
	return F.Pipe2(rc, copyToFile, createTempWriter)
}
