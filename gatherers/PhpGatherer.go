package gatherers

import (
	"log"
	"github.com/bomoko/lagoon-facts/utils"
)

type phpGatherer struct {
	PhpVersion string
}

func (p *phpGatherer) AppliesToEnvironment() bool {
	err, stdOut, stdErr := utils.Shellout("php -r \"echo phpversion();\" | sed -ne 's/[^0-9]*\\(\\([0-9]\\.\\)\\{0,4\\}[0-9][^.]\\).*/\\1/p'")
	if err != nil {
		log.Printf("PhpVersion gatherer cannot be applied: %v", stdErr)
		return false
	}
	p.PhpVersion = stdOut
	log.Printf("Found PHP version: %v", p.PhpVersion)
	return true
}

func (p *phpGatherer) GatherFacts() ([]GatheredFact, error) {

	return []GatheredFact{
		{
			Name: "php-version",
			Value: p.PhpVersion,
			Source: "php-version",
			Description: "This is the current running php version on the system",
			Category: ProgrammingLanguage,
		},
	}, nil
}

func init()  {
	RegisterGatherer("Php Version Gatherer", &phpGatherer{})
}