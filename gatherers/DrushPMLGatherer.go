package gatherers

import (
	"encoding/json"
	"fmt"
	"github.com/bomoko/lagoon-facts/utils"
	"log"
)

type drushPmlGatherer struct {
	GatheredFacts []GatheredFact
}

type drushPmlEntry struct {
	Package string
	Name string
	Type string
	Status string
	Version string
}


func (p *drushPmlGatherer) AppliesToEnvironment() bool {

	err, stdOut, stdErr := utils.Shellout("drush pml --format=json 2> /dev/null")
	if err != nil {
		log.Printf("Drush gatherer cannot be applied: %v", stdErr)
		return false
	}

	var result map[string]drushPmlEntry

	if err = json.Unmarshal([]byte(stdOut), &result); err != nil {
		log.Println(err.Error())
		return false
	}

	for key, element := range result {
		p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
			Name:        key,
			Value:       element.Version,
			Source:      "drush_pml",
			Description: "Drupal module/theme status: " + element.Status,
			Category:  fmt.Sprintf("Drupal %s", element.Type),
		})
	}

	return true
}

func (p *drushPmlGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init()  {
	RegisterGatherer("Drupal Module List Gatherer", &drushPmlGatherer{})
}
