package gatherers

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/uselagoon/lagoon-facts-app/utils"
)

type drushPmlGatherer struct {
	GatheredFacts []GatheredFact
}

func (p *drushPmlGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

type drushPmlEntry struct {
	Package string
	Name    string
	Type    string
	Status  string
	Version interface{}
}

func (p *drushPmlGatherer) AppliesToEnvironment() bool {

	err, stdOut, stdErr := utils.Shellout("drush pml --core --format=json 2> /dev/null")
	if err != nil {
		log.Printf("Drush pml gatherer cannot be applied: %v", stdErr)
		return false
	}

	var result map[string]drushPmlEntry

	if err = json.Unmarshal([]byte(stdOut), &result); err != nil {
		log.Println(err.Error())
		return false
	}

	for key, element := range result {
		if element.Version != nil {
			p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
				Name:        key,
				Value:       fmt.Sprintf("%v", element.Version),
				Source:      "drush_pml",
				Description: "Status: " + element.Status,
				Category:    Drupal,
			})
		}
	}

	return true
}

func (p *drushPmlGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Drupal Module List Gatherer", &drushPmlGatherer{})
}
