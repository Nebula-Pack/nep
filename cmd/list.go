package cmd

import (
	"fmt"
	"nep/utils"
	"os"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List packages",
	Run: func(cmd *cobra.Command, args []string) {
		// Change working directory if path is set
		projectPath := utils.Prepare(false, path)

		keys := [][]string{
			{"dependencies"},
		}

		results, err := utils.ReadConfig(projectPath, keys)
		if err != nil {
			fmt.Printf("Error reading config: %v\n", err)
			os.Exit(1)
		}

		// Convert the result to a map
		dependencies := map[string]string{}
		if len(results) > 0 {
			depMap, ok := results[0].(map[string]interface{})
			if !ok {
				fmt.Println("Error: dependencies are not in the expected format")
				os.Exit(1)
			}

			for pkg, version := range depMap {
				dependencies[pkg] = fmt.Sprintf("%v", version)
			}
		}

		headers := []string{"Package", "Version"}
		rows := [][]string{}
		for pkg, version := range dependencies {
			rows = append(rows, []string{pkg, version})
		}

		// Display table
		err = utils.DisplayTable(headers, rows)
		if err != nil {
			fmt.Printf("Error displaying table: %v\n", err)
			os.Exit(1)
		}

	},
}

func init() {
	listCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(listCmd)
}
