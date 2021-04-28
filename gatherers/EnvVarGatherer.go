package gatherers

import (
	"os"
)

type envVarGatherer struct {
	GatheredFacts []GatheredFact
}

type envVar struct {
	Key string
	Value string
	Description string
}

func (p *envVarGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

func (p *envVarGatherer) AppliesToEnvironment() bool {
	var lagoonDomain = os.Getenv("LAGOON_DOMAIN")
	if (lagoonDomain == "") {
		lagoonDomain = "N/A - UNSET"
	}

	var lagoonEnvType = os.Getenv("LAGOON_ENVIRONMENT_TYPE")
	if (lagoonEnvType == "") {
		lagoonEnvType = "N/A - UNSET"
	}

	var composerVersion = os.Getenv("COMPOSER_VERSION")
	if (composerVersion == "") {
		composerVersion = "N/A - UNSET"
	}

	envVars := map[string]envVar{
		"LAGOON_DOMAIN": {
			Key: "LAGOON_DOMAIN",
			Value: lagoonDomain,
			Description: "The domain address of this environment",
		},
		"LAGOON_ENVIRONMENT_TYPE": {
			Key: "LAGOON_ENVIRONMENT_TYPE",
			Value: lagoonEnvType,
			Description: "This is a '" + lagoonEnvType + "' environment type",
		},
		"COMPOSER_VERSION": {
			Key: "COMPOSER_VERSION",
			Value: composerVersion,
			Description: "Composer version '" + composerVersion + "' was found",
		},
	}

	for key, element := range envVars {
		p.GatheredFacts = append(p.GatheredFacts, GatheredFact{
			Name: key,
			Value: element.Value,
			Source: "env",
			Description: element.Description,
			Category:  EnvVar,
		})
	}

	return true
}

func (p *envVarGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init()  {
	RegisterGatherer("Environment variables gatherer", &envVarGatherer{})
}