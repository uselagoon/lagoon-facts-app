package gatherers

import (
	"context"
	"fmt"
	"github.com/bomoko/lagoon-facts/utils"
	"github.com/shurcooL/graphql"
	"golang.org/x/oauth2"
	"log"
	"os"
)

const lagoonAPIEndpoint = "https://api.lagoon.amazeeio.cloud/graphql"

func Writefacts(environmentName string, facts []GatheredFact) error {

	//let's grab all the environments we'll need (for now it'll almost certainly just be a single one
	//but we support multiple environments

	return nil
}

func getGraphqlClient() (*graphql.Client, error) {
	ctx := context.Background()

	token, err := utils.GetToken()
	if err != nil {
		return nil, err
	}

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}))

	client := graphql.NewClient(lagoonAPIEndpoint, httpClient)
	return client, nil
}

func GetEnvironmentId(environmentName string) (int, error) {

	client, err := getGraphqlClient()
	if err != nil {
		return 0, err
	}

	var projectQuery struct {
		ProjectName struct {
			Id graphql.Int
		} `graphql:"projectByName(name: $name)"`
	}

	err = client.Query(context.Background(), &projectQuery, map[string]interface{}{
		"name": graphql.String(environmentName),
	})
	if err != nil {
		log.Printf(err.Error())
		os.Exit(1)
	}
	fmt.Printf("id : %v", projectQuery.ProjectName.Id)
	return int(projectQuery.ProjectName.Id), nil
}
