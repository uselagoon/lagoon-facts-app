package gatherers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	insightRemoteLib "github.com/uselagoon/insights-remote-lib"
	"log"
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
			Category:    string(fact.Category),
		}

		insightsRemoteFacts.Facts = append(insightsRemoteFacts.Facts, f)
	}

	bodyString, err := json.Marshal(insightsRemoteFacts)
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Printf("Sending %v fact(s) to insights core\n", len(facts))

	serviceEndpoint := viper.GetString("insights-remote-endpoint")
	req, _ := http.NewRequest(http.MethodPost, serviceEndpoint, bytes.NewBuffer(bodyString))
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	response, err := client.Do(req)

	if err != nil {
		return err
	}

	if response.StatusCode != 200 {
		log.Fatalf("There was an error sending the facts to '%s' : %s\n", serviceEndpoint, response.Body)
	}

	defer response.Body.Close()

	return err
}
