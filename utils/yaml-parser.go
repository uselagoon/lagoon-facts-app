package utils

import (
	"errors"
	"gopkg.in/yaml.v2"
)

type DockerComposeConfig struct {
	Project       	string   				  `yaml:"x-lagoon-project"`
	Volumes       	map[string]interface{}    `yaml:"x-volumes"`
	Environment     map[string]interface{}    `yaml:"x-environment"`
	Services    	map[string]interface{} 	  `yaml:"services"`
}

//take bytes and return parsed yaml config struct
func UnmarshallYamlToStructure(data []byte) (DockerComposeConfig, error) {
	config := DockerComposeConfig{}
	err := yaml.Unmarshal(data, &config)

	if err != nil {
		return DockerComposeConfig{}, errors.New("Unable to parse yaml config")
	}
	return config, nil
}