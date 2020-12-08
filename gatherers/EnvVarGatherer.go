package gatherers

import (
	"os"
)

type envVarGatherer struct {
	GatheredFacts []GatheredFact
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
			Name: "LAGOON_VERSION",
			Value: lagoonVersion,
			Source: "env",
			Description: "This is the current LAGOON_VERSION running in the environment",
		},
	}, nil
}

func init()  {
	RegisterGatherer("EnvVar Gatherer", &envVarGatherer{})
}