package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/uselagoon/lagoon-facts-app/gatherers"
	"log"
	"os"
)

var tokenValue string
var tokenFile string
var insightsRemoteEndpoint string

// gatherCmd represents the gather command
var gatherInClusterCmd = &cobra.Command{
	Use:   "gather-in-cluster",
	Short: "Running this command will invoke the registered gatherers in cluster",
	Long:  `Running all the registered gatherers will inspect the system and write FACT data back to the Lagoon insights system via insights-remote`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//get the basic env vars
		if argStatic && argDynamic {
			log.Fatalf("Cannot use both 'static' and 'dynamic' only gatherers - exiting")
			os.Exit(1)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {

		if tokenValue != "" && tokenFile != "" {
			log.Fatal("Either a token or a token file needs to be passed as an argument, not both")
		}

		if tokenValue == "" {
			//tokenValue = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJFbnZpcm9ubWVudElkIjoiMiIsIkVudmlyb25tZW50TmFtZSI6Im1haW4iLCJOYW1lc3BhY2UiOiJkZWxldGVtZSIsIlByb2plY3ROYW1lIjoiZGVsZXRlbWUifQ.r5hoa7-459Kl87fOinCIJPTkjbBw9rxOl2R8KKKQahQ"

			if tokenFile == "" {
				log.Fatal("Either a token or a token file needs to be passed as an argument")
			}

			_, err := os.Stat(tokenFile)
			if err != nil {
				log.Fatal(err)
			}

			ba, err := os.ReadFile(tokenFile)
			if err != nil {
				log.Fatal(err)
			}
			tokenValue = string(ba)

		}

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
			fmt.Print(facts)
			//err := gatherers.Writefacts(projectName, environmentName, facts)
			err := gatherers.WriteFactsToInsightsRemote(tokenValue, facts)
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

//var GatherCommand = gatherCmd

func init() {
	gatherInClusterCmd.PersistentFlags().StringVarP(&tokenValue, "token", "t", "", "The Lagoon insights remote token")
	gatherInClusterCmd.PersistentFlags().StringVarP(&tokenFile, "token-file", "", "", "Read the Lagoon insights remote token from a file")
	gatherInClusterCmd.PersistentFlags().BoolVarP(&dryRun, "dry-run", "d", false, "run gathers and print to screen without running write methods")
	gatherInClusterCmd.PersistentFlags().StringVar(&insightsRemoteEndpoint, "insights-remote-endpoint", "http://localhost:10999/facts", "The Lagoon insights remote endpoint")
	viper.BindPFlag("insights-remote-endpoint", gatherInClusterCmd.PersistentFlags().Lookup("insights-remote-endpoint"))
	rootCmd.AddCommand(gatherInClusterCmd)
}
