package gatherers

import (
	"log"
)

// GatheredFact maps to the Lagoon GraphQL AddFactsInput
type GatheredFact struct {
	Name string `json:"name"`
	Value string  `json:"value"`
	Source string  `json:"source"`
	Environment int  `json:"environment"`
	Description string  `json:"description"`
	Category FactCategory `json:"category,omitempty"`
}

const (
	GATHERER_TYPE_STATIC  string = "static"
	GATHERER_TYPE_DYNAMIC string = "dynamic"
)

type Gatherer interface {
	GetGathererCmdType() string
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