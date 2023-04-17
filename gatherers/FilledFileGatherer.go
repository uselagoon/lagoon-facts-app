package gatherers

import (
	"path/filepath"
)

// This Gatherer will search a particular directory for .json files that contain arrays of facts

type filledFileGatherer struct {
	GatheredFacts []GatheredFact
}

func (p *filledFileGatherer) GetGathererCmdType() string {
	return GATHERER_TYPE_STATIC
}

func (p *filledFileGatherer) AppliesToEnvironment() bool {

	applies := false

	// Look into /tmp/facts for any .json file
	filepath.Glob("/tmp/facts/*.json")

	return applies
}

func (p *filledFileGatherer) GatherFacts() ([]GatheredFact, error) {
	return p.GatheredFacts, nil
}

func init() {
	RegisterGatherer("Application gatherer", &filledFileGatherer{})
}
