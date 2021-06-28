package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/bomoko/lagoon-facts/gatherers"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
)

var gatheredFileName string

// gatherCmd represents the gather command
var filegatherCmd = &cobra.Command{
	Use:   "filegather",
	Short: "Running this command will invoke only invoke the filesystem gatherer",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		var facts []gatherers.GatheredFact


		//we can just unmarshall the file data into the facts ...
		if gatheredFileName == "" {
			fmt.Errorf("Filename should be passed as argument")
			os.Exit(1)
		}

		file, _ := ioutil.ReadFile(gatheredFileName)
		_ = json.Unmarshal([]byte(file), &facts)

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

func init() {
	filegatherCmd.PersistentFlags().StringVarP(&gatheredFileName, "file-name", "f", "", "The file containing the gathered facts")
	gatherCmd.AddCommand(filegatherCmd)
}
