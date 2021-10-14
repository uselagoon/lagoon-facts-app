package gatherers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/uselagoon/lagoon-facts-app/utils"
)

type drushGatherer struct {
	DrupalVersion string
	DrushVersion  string
}

func (p *drushGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

func (p *drushGatherer) AppliesToEnvironment() bool {
	err, stdOut, stdErr := utils.Shellout("drush status --format=json 2> /dev/null")
	if err != nil {
		log.Printf("Drush status gatherer cannot be applied: %v", stdErr)
		return false
	}

	//the following unmarshalls a flat JSON object into key/value accessible structures.
	var result map[string]interface{}
	json.Unmarshal([]byte(stdOut), &result)

	p.DrupalVersion = fmt.Sprint(result["drupal-version"])
	log.Println("Found Drupal Version: " + p.DrupalVersion)
	p.DrushVersion = fmt.Sprint(result["drush-version"])
	log.Println("Found Drush Version: " + p.DrushVersion)

	return true
}

func (p *drushGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name:        "drupal-version",
			Value:       p.DrupalVersion,
			Source:      "drush_status",
			Description: "Currently installed version of Drupal on the Environment",
			Category:    Framework,
			KeyFact: 	 true,
		},
		{
			Name:        "drush-version",
			Value:       p.DrushVersion,
			Source:      "drush_status",
			Description: "Currently installed version of Drush on the Environment",
			Category:    Drupal,
		},
	}, nil
}

func init() {
	RegisterGatherer("Drush gatherer", &drushGatherer{})
}
