package gatherers

import (
	"github.com/uselagoon/lagoon-facts-app/utils"
	"log"
	"os"
)

type platformGatherer struct {
	Name     string
	Version  string
	Category FactCategory
	Description string
}

func (p *platformGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_DYNAMIC
}

func (p *platformGatherer) AppliesToEnvironment() bool {
	applies := false

	homepage, hasLagoonRoute := os.LookupEnv("LAGOON_ROUTE")
	if !hasLagoonRoute {
		log.Printf("LAGOON_ROUTE cannot be determined")
	}

	// find if running on Lagoon
	var findLagoonHttpHeader utils.HTTPHeaderOutput
	findLagoonHttpHeader, err := utils.GetURLHeaderByKey(homepage, "x-lagoon")
	if err != nil {
		log.Printf("Error getting HTTP header: %s", err)
	}

	// find if running on Pantheon
	var findPantheonHTTPHeader utils.HTTPHeaderOutput
	findPantheonHTTPHeader, err = utils.GetURLHeaderByKey(homepage, "x-pantheon-styx-hostname")
	if err != nil {
		log.Printf("Error getting HTTP header: %s", err)
	}

	if err == nil {
		if findLagoonHttpHeader.Name != "" {
			p.Name = "Lagoon Platform"
			p.Category = Platform
			p.Description = "Lagoon is the open-source web hosting platform that enables global teams to scale with ease."
			lagoonVersion, _ := os.LookupEnv("LAGOON_VERSION")
			p.Version = lagoonVersion
			log.Printf("Found PaaS: '%s'", p.Name)
			applies = true
		}
		if findPantheonHTTPHeader.Name != "" {
			p.Name = "Pantheon"
			p.Category = Platform
			p.Version = "-"
			log.Printf("Found PaaS: '%s'", p.Name)
			applies = true
		}
	}

	return applies
}

func (p *platformGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name:        p.Name,
			Value:       p.Version,
			Source:      "http_header",
			Description: p.Description,
			Category:  	 p.Category,
		},
	}, nil
}

func init() {
	RegisterGatherer("PaaS gatherer", &platformGatherer{})
}
