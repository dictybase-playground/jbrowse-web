package file

import (
	"archive/zip"
	"dictybase-playground/jbrowse-web/internal/functional"
	"errors"
	"fmt"
	A "github.com/IBM/fp-go/v2/array"
	E "github.com/IBM/fp-go/v2/either"
	F "github.com/IBM/fp-go/v2/function"
	IOE "github.com/IBM/fp-go/v2/ioeither"
	IOEF "github.com/IBM/fp-go/v2/ioeither/file"
	O "github.com/IBM/fp-go/v2/option"
	PAIR "github.com/IBM/fp-go/v2/pair"
	S "github.com/IBM/fp-go/v2/string"
	"io"
	"os"
	"path/filepath"
)

var (
	IOECreateParentDir = F.Curry2(IOE.Eitherize2(createParentDir))
	IOECopyZipFile     = F.Curry2(copyZipFile)
)

func getFileInfo(file *os.File) IOEither[error, os.FileInfo] {
	return IOE.TryCatchError(func() (os.FileInfo, error) {
		return file.Stat()
	})
}

func getZipReaderParameters(file *os.File) IOEither[error, PAIR.Pair[io.ReaderAt, int64]] {
	return F.Pipe1(
		getFileInfo(file),
		IOE.Map[error](func(info os.FileInfo) PAIR.Pair[io.ReaderAt, int64] {
			return PAIR.MakePair[io.ReaderAt, int64](file, info.Size())
		}),
	)
}

var getZipReader = F.Curry2(F.Swap(IOE.Eitherize2(zip.NewReader)))

func getZipFiles(r *zip.Reader) []*zip.File {
	return r.File
}

func createParentDir(destDir string, file *zip.File) (*zip.File, error) {
	if err := os.MkdirAll(filepath.Join(destDir, filepath.Dir(file.Name)), 0755); err != nil {
		return nil, fmt.Errorf("could not create directory: %w", err)
	}
	return file, nil
}

func ReadZip(zf *zip.File) IOEither[error, io.ReadCloser] {
	return IOE.TryCatchError(func() (io.ReadCloser, error) {
		return zf.Open()
	})
}

func copyZipFile(destDir string, zf *zip.File) IOEither[error, *os.File] {
	withDest := F.Pipe3(
		destDir,
		S.Append(zf.Name),
		IOEF.Create,
		IOEF.Read[*os.File],
	)

	return F.Pipe2(zf, ReadZip, IOE.Chain(F.Flow2(copyToFile, withDest)))
}

func openZip(file *os.File) IOEither[error, []*zip.File] {
	return F.Pipe3(
		file,
		getZipReaderParameters,
		IOE.Chain(PAIR.Merge(getZipReader)),
		IOE.Map[error](getZipFiles))
}

// func processZipFile(destDir string, zf []*zip.File) IOEither[error, []*os.File] {
// 	return F.Pipe1(
// 		zf,
// 		IOE.TraverseArraySeq[error, *zip.File, *os.File](IOECopyZipFile(destDir)),
// 	)
// }

// FunctionalExtractZipFile copies a single zip file a destination folder, preserving the directory stucture of the copied file.
func extractZipFile(dest string, zf *zip.File) O.Option[IOEither[error, *os.File]] {
	return F.Pipe2(
		zf,                     // the zip file
		O.FromPredicate(IsDir), // Check if the file is a directory
		O.Map(
			F.Flow2(
				IOECreateParentDir(dest), // Create the parent directory
				IOE.Chain(IOECopyZipFile(dest)),
			),
		),
	)
}

func extractZipArchive(destDir string, file *os.File) {
	curriedExtractZipFile := F.Curry2(extractZipFile)
	a := F.Pipe2(
		file,
		openZip,
	)
}
