package confer

import (
	"errors"

	"github.com/RedEpochDB/RedEpochDB/lib/filer"

	"gopkg.in/yaml.v2"
)

func ParseYamlConfFromBytes(yamlOrignData []byte) (conf REServerConfS, err error) {

	if len(yamlOrignData) == 0 {
		return REServerConfS{}, errors.New("The configuration content cannot be empty.")
	}

	err = yaml.Unmarshal(yamlOrignData, &conf)
	if err != nil {
		return REServerConfS{}, err
	}

	return
}

func ParseYamlFromFile(yamlFileUri string) (conf REServerConfS, err error) {

	var fileData []byte

	fileData, err = filer.ReadFileData(yamlFileUri)

	if err != nil {
		return REServerConfS{}, err
	}

	conf, err = ParseYamlConfFromBytes(fileData)
	return
}
