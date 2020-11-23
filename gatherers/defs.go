package gatherers

import "log"

type GatheredFact struct {
	Environment string
	Name string
	Value string
	Source string
	Description string
}

type Gatherer interface {
	AppliesToEnvironment() bool //Whether this gatherer can run in the local environment
	GatherFacts() ([]GatheredFact, error)
}

var gathererInternalMap []Gatherer

func RegisterGatherer(gatherer Gatherer) {
	log.Print("registering: ")
	log.Println(gatherer)
	gathererInternalMap = append(gathererInternalMap, gatherer)
}

func GetGatherers() []Gatherer {
	return gathererInternalMap
}