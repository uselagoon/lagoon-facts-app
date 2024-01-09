package cmd

import (
	"encoding/json"
	"github.com/spf13/cobra"
	"github.com/uselagoon/lagoon-facts-app/gatherers"
	"io/ioutil"
	"log"
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
			log.Fatalf("Filename should be passed as argument")
		}

		file, err := ioutil.ReadFile(gatheredFileName)

		if err != nil {
			log.Fatalf("Error reading file: %v", err.Error())
		}

		err = json.Unmarshal([]byte(file), &facts)

		if err != nil {
			log.Fatalf("Error unmarshalling file %v: %v", gatheredFileName, err.Error())
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

func init() {
	filegatherCmd.PersistentFlags().StringVarP(&gatheredFileName, "file-name", "f", "", "The file containing the gathered facts")
	gatherCmd.AddCommand(filegatherCmd)
}
