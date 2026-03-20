package file

import (
	"archive/zip"
	"fmt"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	O "github.com/IBM/fp-go/v2/option"
	"io"
	"os"
	"path/filepath"
)

var (
	ioeCreateParentDir = IOE.Eitherize1(createParentDir)
	ioeCopy            = F.Curry2(IOE.Eitherize2(io.Copy))
	ioeCopyZipFile     = IOE.Eitherize1(copyZipFile)
)

func createParentDir(file *zip.File) (*zip.File, error) {
	if err := os.MkdirAll(filepath.Dir(file.Name), 0755); err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}
	return file, nil
}

// FunctionalExtractZipFile copies a single zip file a destination folder, preserving the directory stucture of the copied file.
func FunctionalExtractZipFile(file *zip.File) error {
	a := F.Pipe2(
		file, // the zip file
		O.FromPredicate(IsDir),
		O.Map(
			F.Flow2(
				ioeCreateParentDir,
				IOE.Chain(ioeCopyZipFile),
			),
		),
	)
}
