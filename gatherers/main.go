package gatherers

import "io/ioutil"

func init() {
}

func LoadDockerComposeConfig(dockerComposePath string) ([]byte, error) {
	var data, err = ioutil.ReadFile(dockerComposePath)
	if err != nil {
		return []byte{}, err
	}
	return data, nil
}