package gatherers

import (
	"context"
	"log"

	//"fmt"
	"github.com/bomoko/lagoon-facts/utils"
	"github.com/machinebox/graphql"
	"golang.org/x/oauth2"
	//"log"
	//"os"
)

const lagoonAPIEndpoint = "https://api.lagoon.amazeeio.cloud/graphql"

func Writefacts(projectName string, environmentName string, facts []GatheredFact) error {
	//
	//projectId, err := GetProjectId(projectName)
	//if err != nil {
	//	return err
	//}
	//
	//environmentId, err := GetEnvironmentId(projectId, environmentName)
	//
	//if err != nil {
	//	return err
	//}
	//
	//
	//sources := map[string]string{}
	//
	//for i,e := range facts{
	//	e.Environment = environmentId
	//	facts[i] = e
	//	if sources[e.Source] == "" {
	//		sources[e.Source] = e.Source
	//	}
	//}
	//
	//for _, e := range sources {
	//	log.Println(e)
	//	err = DeleteFactsBySource(environmentId, e)
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//}
	//
	//client, err := getGraphqlClient()
	//if err != nil {
	//	return err
	//}
	//var addFactMutation struct{
	//	AddFacts struct{
	//		Id graphql.Int
	//	} `graphql:"addFact(input:{name: $name, environment: $environment, value: $value, description: $description, source: $source})"`
	//}
	//
	//
	////factsMarshalledString, err := json.Marshal(facts)
	////if err != nil {
	////	return err
	////}
	//for _, e := range facts {
	//	err = client.Mutate(context.Background(), &addFactMutation, map[string]interface{}{
	//		"name": graphql.String(e.Name),
	//		"environment": graphql.Int(e.Environment),
	//		"value": graphql.String(e.Value),
	//		"description": graphql.String(e.Description),
	//		"source": graphql.String(e.Source),
	//	})
	//
	//	if err != nil {
	//		log.Println(err.Error())
	//	}
	//}

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
			Id int
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

	var responseText struct{
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