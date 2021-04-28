package gatherers

import (
	"fmt"
	"log"
)

type lagoonYamlGatherer struct {
	GatheredFacts []GatheredFact
}

func (p *lagoonYamlGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

type LagoonYamlConfig struct {
	Tasks        map[string]postRollout `yaml:"tasks,omitempty"`
	Routes       map[string]interface{} `yaml:"routes,omitempty"`
	Environments map[string]interface{} `yaml:"environments,omitempty"`
}

type (
	postRollout []map[string]Task
	Task        struct {
		Name    string `yaml:"name"`
		Command string `yaml:"command"`
		Service string `yaml:"service"`
	}
)

func (p *lagoonYamlGatherer) AppliesToEnvironment() bool {
	applies := false

	lagoonConfigBytestream, err := LoadYamlConfig(".lagoon.yml")
	if err != nil {
		log.Printf("Couldn't load lagoon.yml file: ", err.Error())
	}

	config, err := UnmarshallLagoonYamlToStructure(lagoonConfigBytestream)
	if err != nil {
		log.Fatalf("There was an issue unmarshalling the lagoon.yml file: %s", err)
	}

	if config.Tasks != nil {
		applies = true

		for _, element := range config.Tasks["post-rollout"] {
			task := element["run"]
			p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
				Name:       task.Name,
				Value:       task.Command,
				Source:      "lagoon-yml",
				Description: fmt.Sprintf("Post deployment task for '%s' service", task.Service),
				Category:    Lagoon,
			})
		}
	}

	return applies
}

func (p *lagoonYamlGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Lagoon yaml gatherer", &lagoonYamlGatherer{})
}
