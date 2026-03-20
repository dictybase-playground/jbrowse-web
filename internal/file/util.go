package file

import (
	"archive/zip"
)

func IsDir(file *zip.File) bool {
	return file.FileInfo().IsDir()
}
