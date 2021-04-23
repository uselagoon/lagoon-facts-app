package gatherers

import (
	"log"
)

type dockerComposeGatherer struct {
	GatheredFacts []GatheredFact
}

type DockerComposeConfig struct {
	Project      string                 `yaml:"x-lagoon-project,omitempty"`
	Volumes      map[string]interface{} `yaml:"x-volumes,omitempty"`
	Environment  map[string]interface{} `yaml:"x-environment,omitempty"`
	Services     map[string]interface{} `yaml:"services,omitempty"`
}

type dockerComposeService struct {
	Name string
	Description string
	Image string
	Labels map[interface {}]interface {}
	Build map[interface {}]interface {}
}

func (p *dockerComposeGatherer) AppliesToEnvironment() bool {
	applies := false

	lagoonConfigBytestream, err := LoadYamlConfig("docker-compose.yml")
	if err != nil {
		log.Printf("Couldn't load docker-compose.yml file: %s", err.Error())
	}

	config, err := UnmarshallDockerComposeYamlToStructure(lagoonConfigBytestream)
	if err != nil {
		log.Fatalf("There was an issue unmarshalling the docker-compose.yml file: %s", err)
	}

	if config.Services != nil {
		applies = true

		for k, element := range config.Services {
			serviceItems := element.(map[interface {}]interface {})

			var service dockerComposeService
			service.Name = k
			service.Description = "Services found in docker-compose.yml"

			for j, val := range serviceItems {
				if str, ok := val.(string); ok {
					if j == "image" {
						service.Image = str
					}
				}
			}

			p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
				Name:        service.Name,
				Value:       service.Image,
				Source:      "docker-compose",
				Description: service.Description,
				Category:    "Docker configuration",
			})
		}
	}

	return applies
}

func (p *dockerComposeGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Docker-compose gatherer", &dockerComposeGatherer{})
}
