package cmd

import (
	"github.com/bomoko/lagoon-facts/utils"
	"log"
	"os"
	"fmt"
	"github.com/bomoko/lagoon-facts/gatherers"
	"github.com/spf13/cobra"
)

var projectName string
var environment string

// gatherCmd represents the gather command
var gatherCmd = &cobra.Command{
	Use:   "gather",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
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
		fmt.Println(gathererSlice)
		var facts []gatherers.GatheredFact


		for _, e := range gathererSlice {
			if e.AppliesToEnvironment() {
				gatheredFacts, err := e.GatherFacts()
				if err != nil {
					fmt.Println(err.Error())
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

// gatherCmd represents the gather command
var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Grab an ssh token",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		token, err := utils.GetToken()
		if err != nil {
			log.Printf(err.Error())
			return
		}
		fmt.Println(token)
		return
	},
}

func init() {
	rootCmd.AddCommand(gatherCmd)
	rootCmd.AddCommand(tokenCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// gatherCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// gatherCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
