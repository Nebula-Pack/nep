package cmd

import (
	"fmt"
	"nep/configs"
	"nep/utils"
	"os"

	"github.com/spf13/cobra"
)

var (
	useCurrentDir bool
	uninteractive bool
)

var initCmd = &cobra.Command{
	Use:   "init [name]",
	Short: "Initialize a new project",
	Args:  cobra.MaximumNArgs(1),
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

		var projectName, des, author, license string
		if !uninteractive {
			prompts := []string{"Project Name", "Description", "Author"}
			responses, err := utils.GroupedTextInput(prompts)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = responses[0]
			}
			des, author = responses[1], responses[2]
			if verbose {
				fmt.Printf("Project Name: %s\nDescription: %s\nAuthor: %s\n", projectName, des, author)
			}
			licenses := []utils.Item{
				{TitleText: "MIT License", Desc: "A permissive license that is short and to the point."},
				{TitleText: "Apache License 2.0", Desc: "A license that allows you to do almost anything with the project."},
				{TitleText: "GNU GPLv3", Desc: "A copyleft license that requires derivative works to be open source."},
				{TitleText: "BSD 3-Clause License", Desc: "A permissive license similar to the MIT License."},
			}
			license, err = utils.SelectFromList("Choose a License", licenses, 3)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				os.Exit(1)
			}
			if verbose {
				fmt.Printf("Selected License: %s\n", license)
			}
		} else {
			if useCurrentDir {
				projectName = configs.DefaultName
			} else if len(args) > 0 {
				projectName = args[0]
			} else {
				projectName = configs.DefaultName

			}
		}
		dir, err := utils.CreateProject(projectName, useCurrentDir)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		if verbose {
			fmt.Printf("Created project directory: %s\n", dir)
		}
		if !uninteractive {
			updates := []utils.UpdatePath{
				{Path: []string{"author"}, Value: author},
				{Path: []string{"name"}, Value: projectName},
				{Path: []string{"license"}, Value: license},
				{Path: []string{"description"}, Value: des},
			}
			err = utils.UpdateConfig(dir, updates)
			if err != nil {
				fmt.Println(err)
			}
		}
		if path != "" {
			fmt.Printf("Initialized empty project in %s\n", path)
			fmt.Printf("To start working on your project, run:\n\n\tcd %s/%s\n\n", path, projectName)
		} else if useCurrentDir {
			fmt.Printf("Initialized empty project in current directory\n")
		} else {
			fmt.Printf("Initialized empty project in %s\n", dir)
			fmt.Printf("To start working on your project, run:\n\n\tcd %s\n\n", projectName)
		}
	},
}

func init() {
	initCmd.Flags().BoolVarP(&useCurrentDir, "current", "c", false, "Use current directory as project directory")
	initCmd.Flags().BoolVarP(&uninteractive, "uninteractive", "i", false, "Use uninteractive mode")
	initCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	initCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(initCmd)
}
