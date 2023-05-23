package gatherers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/spf13/viper"
	insightRemoteLib "github.com/uselagoon/insights-remote-lib"
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
			Category:    string(fact.Category),
		}

		insightsRemoteFacts.Facts = append(insightsRemoteFacts.Facts, f)
	}

	bodyString, _ := json.Marshal(insightsRemoteFacts)
	req, _ := http.NewRequest(http.MethodPost, viper.GetString("insights-remote-endpoint"), bytes.NewBuffer(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)
	if response.StatusCode != 200 {
		fmt.Println(response.Body)
		os.Exit(1)
	}
	return err
}
