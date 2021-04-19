package gatherers

import (
	"encoding/json"
	"fmt"
	"github.com/bomoko/lagoon-facts/utils"
)

type applicationGatherer struct {
	Name     string
	Version string
}

type ComposerShowOutput struct {
	Name     string
	Versions []string
}

func (p *applicationGatherer) AppliesToEnvironment() bool {

	a := []string{"drupal/core", "laravel/framework"}

	applies := false

	for _, name := range a {
		fmt.Printf("composer show -i --format=json %v 2> /dev/null", name)
		err, stdOut, _ := utils.Shellout(fmt.Sprintf("composer show -i --format=json %v 2> /dev/null", name))
		if err == nil {
			var result ComposerShowOutput
			json.Unmarshal([]byte(stdOut), &result)
			p.Name = result.Name
			p.Version = result.Versions[0]
			applies = true
		}
	}

	return applies
}

func (p *applicationGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name:        "application_type",
			Value:       p.Name,
			Source:      "system_application",
			Description: "The current application installed on the environment",
		},
		{
			Name:        "application_version",
			Value:       p.Version,
			Source:      "system_application",
			Description: "The current application installed on the environment",
		},
	}, nil
}

func init() {
	RegisterGatherer("Application type Gatherer", &applicationGatherer{})
}
