package gatherers

import (
	"bytes"
	"encoding/json"
	"github.com/spf13/viper"
	"github.com/uselagoon/insights-remote-lib"
	"net/http"
)

func WriteFactsToInsightsRemote(token string, facts []GatheredFact) error {

	insightsRemoteFacts := insightRemoteLib.Facts{
		Facts: []insightRemoteLib.Fact{},
	}

	for _, fact := range facts {
		f := insightRemoteLib.Fact{
			Name:        fact.Name,
			Value:       fact.Value,
			Source:      fact.Source,
			Description: fact.Description,
			Type:        "type",
			Category:    string(fact.Category),
			//Service:         fact.,
		}

		insightsRemoteFacts.Facts = append(insightsRemoteFacts.Facts, f)
	}

	bodyString, _ := json.Marshal(insightsRemoteFacts)
	req, _ := http.NewRequest(http.MethodPost, viper.GetString("insights-remote-endpoint"), bytes.NewBuffer(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	_, err := client.Do(req)
	return err
}
