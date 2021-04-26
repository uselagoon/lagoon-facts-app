package gatherers

import (
	"github.com/bomoko/lagoon-facts/utils"
	"log"
	"os"
)

type technologyGatherer struct {
	Name     	string
	Value  		string
	Version  	string
	Description string
	Category 	FactCategory
}

func (p *technologyGatherer) AppliesToEnvironment() bool {
	applies := false

	homepage, hasLagoonRoute := os.LookupEnv("LAGOON_ROUTE")
	if !hasLagoonRoute {
		log.Printf("LAGOON_ROUTE cannot be determined")
	}

	var findWebServer utils.HTTPHeaderOutput
	findWebServer, err := utils.GetURLHeaderByKey(homepage, "Server")
	if err != nil {
		log.Printf("Error getting HTTP header: %s", err)
	}

	if err == nil {
		if findWebServer.Name != "" {
			p.Name = findWebServer.Name
			p.Description = findWebServer.Name
			p.Category = ApplicationTechnology
			p.Value = findWebServer.Value
			log.Printf("Found web server: '%s'", p.Name)
			applies = true
		}
	}

	return applies
}

func (p *technologyGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name:        p.Value,
			Value:       p.Value,
			Source:      "http_header",
			Description: p.Description,
			Category:  	 p.Category,
		},
	}, nil
}

func init() {
	RegisterGatherer("Technology gatherer", &technologyGatherer{})
}
