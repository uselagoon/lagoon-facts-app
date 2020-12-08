package gatherers

import "log"

// GatheredFact maps to the Lagoon GraphQL AddFactsInput
type GatheredFact struct {
	Environment int  `json:"environment"`
	Name string `json:"name"`
	Value string  `json:"value"`
	Source string  `json:"source"`
	Description string  `json:"description"`
}

type Gatherer interface {
	AppliesToEnvironment() bool //Whether this gatherer can run in the local environment
	GatherFacts() ([]GatheredFact, error)
}

var gathererInternalMap []Gatherer

func RegisterGatherer(name string, gatherer Gatherer) {
	log.Print("registering: " + name)
	gathererInternalMap = append(gathererInternalMap, gatherer)
}

func GetGatherers() []Gatherer {
	return gathererInternalMap
}