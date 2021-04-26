package gatherers

import (
	"errors"
	"gopkg.in/yaml.v2"
	"log"
)

// GatheredFact maps to the Lagoon GraphQL AddFactsInput
type GatheredFact struct {
	Name string `json:"name"`
	Value string  `json:"value"`
	Source string  `json:"source"`
	Environment int  `json:"environment"`
	Description string  `json:"description"`
	Category FactCategory `json:"category,omitempty"`
}

const (
	Static string = "static"
	Dynamic string = "dynamic"
)

type Gatherer interface {
	GetGathererCmdType() string
	AppliesToEnvironment() bool //Whether this gatherer can run in the local environment
	GatherFacts() ([]GatheredFact, error)
}

var gathererInternalMap []Gatherer

func RegisterGatherer(name string, gatherer Gatherer) {
	log.Print("registering: " + name)
	gathererInternalMap = append(gathererInternalMap, gatherer)
}

func GetGatherers() []Gatherer {
	return gathererInternalMap
}

func UnmarshallDockerComposeYamlToStructure(data []byte) (DockerComposeConfig, error) {
	config := DockerComposeConfig{}
	err := yaml.Unmarshal(data, &config)

	if err != nil {
		return DockerComposeConfig{}, errors.New("Unable to parse docker-compose.yml config")
	}
	return config, nil
}

func UnmarshallLagoonYamlToStructure(data []byte) (LagoonYamlConfig, error) {
	config := LagoonYamlConfig{}
	err := yaml.Unmarshal(data, &config)

	if err != nil {
		return LagoonYamlConfig{}, errors.New("Unable to parse lagoon.yml config")
	}
	return config, nil
}