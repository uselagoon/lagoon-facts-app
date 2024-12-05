package gatherers

import (
	"context"
	"fmt"
	"github.com/machinebox/graphql"
	"github.com/uselagoon/lagoon-facts-app/utils"
	"golang.org/x/oauth2"
	"log"
)

const lagoonAPIEndpoint = "https://api.lagoon.amazeeio.cloud/graphql"
const lagoonUIEndpoint = "https://dashboard.amazeeio.cloud"

func Writefacts(projectName string, environmentName string, facts []GatheredFact) error {

	projectId, err := GetProjectId(projectName)
	if err != nil {
		return err
	}

	environmentId, err := GetEnvironmentId(projectId, environmentName)

	if err != nil {
		return err
	}

	sources := map[string]string{}

	for i, e := range facts {
		e.Environment = environmentId
		facts[i] = e
		if sources[e.Source] == "" {
			sources[e.Source] = e.Source
		}
	}

	for _, e := range sources {
		log.Println("Clearing source: " + e)
		err = DeleteFactsBySource(environmentId, e)
		if err != nil {
			log.Println(err.Error())
		}
	}

	// let's split the array into multidimensional arrays in case we have huge fact payloads
	var factMD [][]GatheredFact

	const maxFactWrite = 100
	for i := 0; i < len(facts); i += maxFactWrite {
		end := i + maxFactWrite
		if end > len(facts) {
			end = len(facts)
		}
		factMD = append(factMD, facts[i:end])
	}

	for _, f := range factMD {
		err := writeFactsToGraphql(err, f)
		if err != nil {
			return err
		}
	}

	var factsUIUrl = fmt.Sprintf("%s/projects/%s/%s/facts", lagoonUIEndpoint, projectName, fmt.Sprintf("%s-%s", projectName, environmentName))
	log.Printf("Successfully added facts to %s:%s \n %s", projectName, environmentName, factsUIUrl)

	return nil
}

func writeFactsToGraphql(err error, facts []GatheredFact) error {
	client, err := getGraphqlClient()
	if err != nil {
		log.Println(err)
		return err
	}

	var addFactMutation struct {
		AddFacts []struct {
			Id int
		}
	}

	req := graphql.NewRequest(`
	mutation addFactMutation($facts: AddFactsInput!) {
  addFacts(input: $facts) {
    id
  }
}
`)

	var factInput struct {
		Facts []GatheredFact `json:"facts"`
	}
	factInput.Facts = facts

	req.Var("facts", factInput)

	ctx := context.Background()

	if err := client.Run(ctx, req, &addFactMutation); err != nil {
		log.Fatal(err)
	}
	return nil
}

func getGraphqlClient() (*graphql.Client, error) {
	ctx := context.Background()

	token, err := utils.GetToken()
	if err != nil {
		log.Fatal(err)
	}

	httpClient := oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}))

	client := graphql.NewClient(lagoonAPIEndpoint, graphql.WithHTTPClient(httpClient))
	return client, nil
}

func GetProjectId(projectName string) (int, error) {

	client, err := getGraphqlClient()
	if err != nil {
		return 0, err
	}

	var projectQuery struct {
		ProjectByName struct {
			Id int
		}
	}

	req := graphql.NewRequest(`
	query getProjectByName($name: String!) {
		projectByName (name: $name) {
			id
		}
	}
`)
	req.Var("name", projectName)

	ctx := context.Background()

	if err := client.Run(ctx, req, &projectQuery); err != nil {
		log.Fatal(err)
	}

	return int(projectQuery.ProjectByName.Id), nil
}

func GetEnvironmentId(projectId int, environmentName string) (int, error) {

	client, err := getGraphqlClient()
	if err != nil {
		return 0, err
	}

	var environmentQuery struct {
		EnvironmentByName struct {
			Id   int
			Name string
		}
	}

	req := graphql.NewRequest(`
	query getEnvironmentByName($project: Int!, $name: String!) {
		environmentByName (project: $project, name: $name) {
			id
			name
		}
	}
`)

	req.Var("project", projectId)
	req.Var("name", environmentName)

	ctx := context.Background()

	if err := client.Run(ctx, req, &environmentQuery); err != nil {
		return 0, err
	}

	return int(environmentQuery.EnvironmentByName.Id), nil
}

func DeleteFactsBySource(environmentId int, source string) error {

	client, err := getGraphqlClient()
	if err != nil {
		return err
	}

	var responseText struct {
		Data string
	}

	req := graphql.NewRequest(`
	mutation deleteFactsFromSourceMutation($environment: Int!, $source: String!) {
		deleteFactsFromSource(input: {environment: $environment, source: $source})
	}
`)

	req.Var("environment", environmentId)
	req.Var("source", source)

	ctx := context.Background()

	if err := client.Run(ctx, req, &responseText); err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}
