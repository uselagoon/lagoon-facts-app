package gatherers

import (
	"github.com/bomoko/lagoon-facts/utils"
	"log"
	"os"
)

type hostingGatherer struct {
	Name     string
	Version  string
}

func (p *hostingGatherer) AppliesToEnvironment() bool {

	applies := false

	homepage, hasLagoonRoute := os.LookupEnv("LAGOON_ROUTE")
	if !hasLagoonRoute {
		log.Printf("LAGOON_ROUTE cannot be determined")
	}

	var findLagoonHttpHeader utils.HTTPHeaderOutput
	findLagoonHttpHeader = utils.GetURLHeaderByKey(homepage, "x-lagoon")

	if findLagoonHttpHeader.Name != "" {
		p.Name = "Lagoon"
		lagoonVersion, _ := os.LookupEnv("LAGOON_VERSION")
		p.Version = lagoonVersion
		applies = true
	}

	return applies
}

func (p *hostingGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name:        "hosting_name",
			Value:       p.Name,
			Source:      "hosting",
			Description: "The current hosting provider name",
		},
		{
			Name:        "hosting_version",
			Value:       p.Version,
			Source:      "hosting",
			Description: "The current hosting provider version",
		},
	}, nil
}

func init() {
	RegisterGatherer("Hosting data gatherer", &hostingGatherer{})
}
