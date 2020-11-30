/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/bomoko/lagoon-facts/gatherers"
	"github.com/spf13/cobra"
)

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

		//run the gatherers...
		gathererSlice := gatherers.GetGatherers()
		fmt.Println(gathererSlice)
		//
		var facts []gatherers.GatheredFact


		//var facts gatherers.GatheredFact

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
		for _, e := range facts {
			fmt.Println(e.Value)
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
		//token, err := utils.GetToken()
		//if err != nil {
		//	log.Printf(err.Error())
		//	return
		//}
		//fmt.Println(token)
		//id, err := gatherers.GetEnvironmentId("heks_ch")
		//if err != nil {
		//	log.Println(err.Error())
		//	return
		//}
		//fmt.Printf("ID : %v", id)

		var gatheredFacts []gatherers.GatheredFact

		gatheredFacts = append(gatheredFacts, gatherers.GatheredFact{
			Environment: 64091,
			Name:        "test1",
			Value:       "1",
			Source:      "test",
			Description: "a test",
		})

		gatheredFacts = append(gatheredFacts, gatherers.GatheredFact{
			Environment: 64091,
			Name:        "test2",
			Value:       "2",
			Source:      "test",
			Description: "a second test",
		})

		//gatherers.Writefacts("amazeelabsv4_com", "dev", gatheredFacts)
		//gatherers.GetProjectId("amazeelabsv4_com")
		//gatherers.GetEnvironmentId(333, "dev")
		gatherers.DeleteFactsBySource(64091, "test")
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
