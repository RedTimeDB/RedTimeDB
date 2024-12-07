package filer

import (
	"errors"
	"os"
)

func ReadFileData(FileUri string) (fileData []byte, err error) {

	if _, err = os.Stat(FileUri); os.IsNotExist(err) {

		return nil, errors.New("the file path does not exist.")
	}
	fileData, err = os.ReadFile(FileUri)
	return
}
