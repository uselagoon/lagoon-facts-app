package cmd

import (
	"github.com/bomoko/lagoon-facts/gatherers"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var projectName string
var environment string

// gatherCmd represents the gather command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Running this command will invoke the registered gatherers",
	Long: `Running all the registered gatherers will inspect the system and write FACT data back to the Lagoon insights system`,
	Run: func(cmd *cobra.Command, args []string) {

		//get the basic env vars
		if projectName == "" {
			projectName = os.Getenv("LAGOON_PROJECT")
		}
		if projectName == "" {
			projectName = os.Getenv("LAGOON_SAFE_PROJECT")
		}
		if environment == "" {
			environment = os.Getenv("LAGOON_GIT_BRANCH")
		}

		if environment == "" || projectName == "" {
			log.Fatalf("PROJECT OR ENVIRONMENT NOT SET - exiting")
			os.Exit(1)
		}

		//run the gatherers...
		gathererSlice := gatherers.GetGatherers()

		var facts []gatherers.GatheredFact


		for _, e := range gathererSlice {
			if e.AppliesToEnvironment() {
				gatheredFacts, err := e.GatherFacts()
				if err != nil {
					log.Println(err.Error())
					continue
				}
				facts = append(facts, gatheredFacts...)
			}
		}

		err := gatherers.Writefacts(projectName, environment, facts)

		if err != nil {
			log.Println(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(gatherCmd)
}
