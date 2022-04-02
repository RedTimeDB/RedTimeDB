/*
 * @Author: gitsrc
 * @Date: 2022-04-02 12:06:42
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 12:06:52
 * @FilePath: /RedTimeDB/lib/filer/main.go
 */

package filer

import (
	"errors"
	"io/ioutil"
	"os"
)

func ReadFileData(FileUri string) (fileData []byte, err error) {

	if _, err = os.Stat(FileUri); os.IsNotExist(err) {

		return nil, errors.New("The file path does not exist.")
	}
	fileData, err = ioutil.ReadFile(FileUri)
	return
}
