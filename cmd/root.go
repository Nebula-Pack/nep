package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	verbose bool
	path    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "nep",
	Version: "0.1.0",
	Short:   "Nebula Pack - A package manager for Lua and LÖVE2D projects",
	Long: `Nebula Pack (nep) is a specialized package manager designed to simplify 
the process of installing and managing libraries for Lua and LÖVE2D projects.

Developed by Keagan Gilmore and Jayden Vawdrey, Nebula Pack allows developers 
to easily integrate libraries into their projects without worrying about complex 
dependencies or configurations.

Key features:
- Streamlined installation of Lua and LÖVE2D libraries
- Simplified package management for Lua projects
- No additional dependencies required

Nebula Pack is ideal for developers who want to focus on creating Lua applications 
without the hassle of manual library setup.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Show additional information")
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Set project path")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
