package cmd

import (
	"encoding/json"
	"github.com/bomoko/lagoon-facts/gatherers"
	"github.com/spf13/cobra"
	"log"
	"os"
)

var projectName string
var environment string
var gatherer bool
var dryRun bool

// gatherCmd represents the gather command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "Running this command will invoke the registered gatherers",
	Long:  `Running all the registered gatherers will inspect the system and write FACT data back to the Lagoon insights system`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
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

		if argStatic && argDynamic {
			log.Fatalf("Cannot use both 'static' and 'dynamic' only gatherers - exiting")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		//set gatherer type to be static by default
		gathererTypeArg := gatherers.GATHERER_TYPE_STATIC
		if argDynamic {
			gathererTypeArg = gatherers.GATHERER_TYPE_DYNAMIC
		}

		//run the gatherers...
		gathererSlice := gatherers.GetGatherers()

		var facts []gatherers.GatheredFact

		for _, e := range gathererSlice {
			if e.GetGathererCmdType() == gathererTypeArg {
				if e.AppliesToEnvironment() {
					gatheredFacts, err := e.GatherFacts()
					if err != nil {
						log.Println(err.Error())
						continue
					}
					facts = append(facts, gatheredFacts...)
				}
			}
		}

		if !dryRun {
			err := gatherers.Writefacts(projectName, environment, facts)

			if err != nil {
				log.Println(err.Error())
			}
		}

		if dryRun {
			if facts != nil {
				log.Println("---- Dry run ----")
				log.Printf("Would post the follow facts to '%s:%s'", projectName, environment)
				s, _ := json.MarshalIndent(facts, "", "\t")
				log.Println(string(s))
			}
		}
	},
}

var GatherCommand = gatherCmd

func init() {
	gatherCmd.PersistentFlags().StringVarP(&projectName, "project-name", "p", "", "The Lagoon project name")
	gatherCmd.PersistentFlags().StringVarP(&environment, "environment-name", "e", "", "The Lagoon environment name")
	gatherCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "run gathers and print to screen without running write methods")
	rootCmd.AddCommand(gatherCmd)
}
