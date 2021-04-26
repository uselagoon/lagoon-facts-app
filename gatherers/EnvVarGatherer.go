package gatherers

import (
	"os"
)

type envVarGatherer struct {
	GatheredFacts []GatheredFact
}

func (p *envVarGatherer) GetGathererCmdType() string {
	return Static
}

func (p *envVarGatherer) AppliesToEnvironment() bool {

	return true
}

func (p *envVarGatherer) GatherFacts() ([]GatheredFact, error) {
	var lagoonVersion = os.Getenv("LAGOON_VERSION")
	if(lagoonVersion == "") {
		lagoonVersion = "N/A - UNSET"
	}
	return []GatheredFact{
		{
			Name: "lagoon",
			Value: lagoonVersion,
			Source: "env",
			Description: "This is the current version of Lagoon",
			Category:  EnvVar,
		},
	}, nil
}

func init()  {
	RegisterGatherer("EnvVar Gatherer", &envVarGatherer{})
}