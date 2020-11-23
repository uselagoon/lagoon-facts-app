package gatherers

import (
	"encoding/json"
	"fmt"
	"log"
	"github.com/bomoko/lagoon-facts/utils"
)

type drushGatherer struct {
	DrupalVersion string
}

type drushStatusJson struct {

}

func (p drushGatherer) AppliesToEnvironment() bool {
	err, stdOut, stdErr := utils.Shellout("drush status --format=json 2> /dev/null")
	if err != nil {
		log.Printf("Drush gatherer cannot be applied: %v", stdErr)
		return false
	}

	//the following unmarshalls a flat JSON object into key/value accessible structures.
	var result map[string]interface{}
	json.Unmarshal([]byte(stdOut), &result)

	p.DrupalVersion = fmt.Sprint(result["drupal-version"])

	return true
}

func (p drushGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Environment: "test",
			Name: "drush-version",
			Value: p.DrupalVersion,
			Source: "drush-status",
			Description: "",
		},
	}, nil
}

func init()  {
	RegisterGatherer(drushGatherer{})
}
