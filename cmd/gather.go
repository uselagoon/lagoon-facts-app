package cmd

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uselagoon/lagoon-facts-app/gatherers"
)

var projectName string
var environmentName string
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
			project, exists := os.LookupEnv("LAGOON_PROJECT")
			if exists {
				projectName = strings.Replace(project, "_", "-", -1)
			}
		}
		if projectName == "" {
			project, exists := os.LookupEnv("LAGOON_SAFE_PROJECT")
			if exists {
				projectName = strings.Replace(project, "_", "-", -1)
			}
		}
		if environmentName == "" {
			environmentName = os.Getenv("LAGOON_GIT_BRANCH")
		}

		if environmentName == "" || projectName == "" {
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
					for _, f := range gatheredFacts {
						if verbose := viper.Get("verbose"); verbose == true {
							log.Printf("Registering %s", f.Name)
						}
					}
					facts = append(facts, gatheredFacts...)
				}
			}
		}

		if !dryRun {
			err := gatherers.Writefacts(projectName, environmentName, facts)

			if err != nil {
				log.Println(err.Error())
			}
		}

		if dryRun {
			if facts != nil {
				log.Println("---- Dry run ----")
				log.Printf("Would post the follow facts to '%s:%s'", projectName, environmentName)
				s, _ := json.MarshalIndent(facts, "", "\t")
				log.Println(string(s))
			}
		}
	},
}

var GatherCommand = gatherCmd

func init() {
	gatherCmd.PersistentFlags().StringVarP(&projectName, "project", "p", "", "The Lagoon project name")
	gatherCmd.PersistentFlags().StringVarP(&environmentName, "environment", "e", "", "The Lagoon environment name")
	gatherCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "run gathers and print to screen without running write methods")
	rootCmd.AddCommand(gatherCmd)
}
