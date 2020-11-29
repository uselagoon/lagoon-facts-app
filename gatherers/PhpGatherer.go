package gatherers

import (
	"log"
	"github.com/bomoko/lagoon-facts/utils"
)

type phpGatherer struct {
	phpVersion string
}

func (p phpGatherer) AppliesToEnvironment() bool {
	err, stdOut, stdErr := utils.Shellout("php -r \"echo phpversion();\"")
	if err != nil {
		log.Printf("PhpVersion gatherer cannot be applied: %v", stdErr)
		return false
	}
	p.phpVersion = stdOut
	return true
}

func (p phpGatherer) GatherFacts() ([]GatheredFact, error) {
	return []GatheredFact{
		{
			Name: "php-version",
			Value: p.phpVersion,
			Source: "php-version",
			Description: "This is the current running php version on the system",
		},
	}, nil
}

func init()  {
	RegisterGatherer(phpGatherer{})
}