package cmd

import (
	"fmt"
	"nep/utils"
	"os"

	"github.com/spf13/cobra"
	lua "github.com/yuin/gopher-lua"
)

type Scripts map[string]string

var (
	path    string
	scripts Scripts
)

func loadScripts() error {
	projectPath := utils.Prepare(true, path)

	// Read the configuration section for "scripts"
	results, err := utils.ReadConfig(projectPath, [][]string{{"scripts"}})
	if err != nil {
		return err
	}

	// Ensure results is in the expected format
	rawScripts, ok := results[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("scripts configuration is not in the expected format")
	}

	// Initialize the scripts map
	scripts = make(map[string]string)

	// Convert rawScripts to map[string]string manually
	for key, value := range rawScripts {
		strValue, ok := value.(string)
		if !ok {
			return fmt.Errorf("script %s is not a string", key)
		}
		scripts[key] = strValue
	}

	return nil
}

var rootCmd = &cobra.Command{
	Use:     "nep [script]",
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
	Args: cobra.MaximumNArgs(1), // Ensures only one argument is allowed
	Run: func(cmd *cobra.Command, args []string) {
		if err := loadScripts(); err != nil {
			fmt.Println("Error loading scripts:", err)
			os.Exit(1)
		}

		if len(args) > 0 {
			scriptName := args[0]
			luaScript, exists := scripts[scriptName]
			if !exists {
				fmt.Printf("script %s not found\n", scriptName)
				return
			}

			L := lua.NewState()
			defer L.Close()

			if err := L.DoString(luaScript); err != nil {
				fmt.Printf("Error running Lua script %s: %v\n", scriptName, err)
			}
		} else {
			fmt.Println("No script name provided.")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&path, "path", "p", "", "Set project path")
}
