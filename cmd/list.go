package cmd

import (
	"encoding/json"
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
		if path != "" {
			err := os.Chdir(path)
			if err != nil {
				fmt.Printf("Error changing directory: %v\n", err)
				os.Exit(1)
			}
			if verbose {
				fmt.Printf("Changed working directory to: %s\n", path)
			}
		}

		projectPath, err := utils.FindProjectDir()
		if err != nil {
			fmt.Printf("Error finding project directory: %v\n", err)
			os.Exit(1)
		}

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

		// Write to file if file flag is set
		if file != "" {
			f, err := os.Create(file)
			if err != nil {
				fmt.Printf("Error creating file: %v\n", err)
				os.Exit(1)
			}
			defer f.Close()

			encoder := json.NewEncoder(f)
			encoder.SetIndent("", "  ")
			err = encoder.Encode(dependencies)
			if err != nil {
				fmt.Printf("Error writing to file: %v\n", err)
				os.Exit(1)
			}

			if verbose {
				fmt.Printf("Wrote package list to file: %s\n", file)
			}
		}
	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	listCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	listCmd.Flags().StringVarP(&file, "file", "f", "", "Set output file name")
	rootCmd.AddCommand(listCmd)
}
