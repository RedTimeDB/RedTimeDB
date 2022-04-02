/*
 * @Author: gitsrc
 * @Date: 2022-04-02 11:57:10
 * @LastEditors: gitsrc
 * @LastEditTime: 2022-04-02 13:20:38
 * @FilePath: /RedTimeDB/app/rtserver/rtserverconf/yamlParser.go
 */
package confer

import (
	"errors"

	"github.com/RedTimeDB/RedTimeDB/lib/filer"

	"gopkg.in/yaml.v2"
)

func ParseYamlConfFromBytes(yamlOrignData []byte) (conf RTServerConfS, err error) {

	if len(yamlOrignData) == 0 {
		return RTServerConfS{}, errors.New("The configuration content cannot be empty.")
	}

	err = yaml.Unmarshal(yamlOrignData, &conf)
	if err != nil {
		return RTServerConfS{}, err
	}

	return
}

func ParseYamlFromFile(yamlFileUri string) (conf RTServerConfS, err error) {

	var fileData []byte

	fileData, err = filer.ReadFileData(yamlFileUri)

	if err != nil {
		return RTServerConfS{}, err
	}

	conf, err = ParseYamlConfFromBytes(fileData)
	return
}
