package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var (
	isolator     string
	seed         int
	multipleArgs []string
)

var compileCmd = &cobra.Command{
	Use:     "compile",
	Aliases: []string{"c"},
	Short:   "Compile the project with a specified type",
	Args:    cobra.ExactArgs(1), // Ensure exactly one argument is provided
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Please specify a compile type. For example: 'nep compile LOVE'")
	},
}

var compileLoveCmd = &cobra.Command{
	Use:     "LOVE [args]",
	Short:   "Compile for LOVE",
	Aliases: []string{"love", "Love", "l", "L"},
	Args:    cobra.ArbitraryArgs, // Allow multiple positional arguments
	Run: func(cmd *cobra.Command, args []string) {
		multipleArgs = args
		compileLove()
	},
}

func compileLove() {
	fmt.Println("Compiling for LOVE...")
	if isolator != "" {
		fmt.Printf("Isolator: %s\n", isolator)
	}
	if seed != 0 {
		fmt.Printf("Seed: %d\n", seed)
	}
	if len(multipleArgs) > 0 {
		fmt.Printf("Additional arguments: %s\n", strings.Join(multipleArgs, ", "))
	}
	// Implement the LOVE compilation logic here
}

func init() {
	// Add global compile command
	rootCmd.AddCommand(compileCmd)

	// Define flags for the LOVE compile type
	compileLoveCmd.Flags().StringVarP(&isolator, "isolator", "i", "", "Description for isolator")
	compileLoveCmd.Flags().IntVarP(&seed, "seed", "s", 0, "Description for seed")

	// Define and add subcommands for specific compile types
	compileCmd.AddCommand(compileLoveCmd)
}
