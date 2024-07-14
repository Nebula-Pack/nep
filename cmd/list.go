package cmd

import (
	"fmt"
	"nep/utils"

	"github.com/spf13/cobra"
)

//work in progress

var listCmd = &cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List packages",
	Run: func(cmd *cobra.Command, args []string) {
		headers := []string{"Name", "Age", "City"}
		rows := [][]string{
			{"Alice", "30", "New York"},
			{"Bob", "25", "San Francisco"},
			{"Charlie", "35", "London"},
		}

		tableString := utils.DisplayTable(headers, rows)
		fmt.Println(tableString)
	},
}

func init() {
	listCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	listCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(listCmd)
}
