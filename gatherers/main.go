package gatherers

import "io/ioutil"

func init() {
}

func LoadYamlConfig(yamlFilePath string) ([]byte, error) {
	var data, err = ioutil.ReadFile(yamlFilePath)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}
