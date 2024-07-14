package cmd

import (
	"fmt"
	"nep/configs"
	"nep/utils"
	"os"
	"path/filepath"

	"strings"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update [packages]",
	Aliases: []string{"up"},
	Short:   "update packages",
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

		cachePath := filepath.Join(projectPath, configs.CacheFolderName)
		packagePath := filepath.Join(projectPath, configs.FolderName)

		if len(args) > 0 && args[0] == configs.All {
			// Update all packages
			updateAllPackages(projectPath, cachePath, packagePath)
		} else {
			// Update specific packages
			for _, pkg := range args {
				installPackage(pkg, projectPath, cachePath)
			}
			updateSpecificPackages(args, cachePath, packagePath)
		}
	},
}

func updateAllPackages(projectPath string, cachePath string, packagePath string) {
	keys := [][]string{
		{"dependencies"},
	}
	packs := []string{}
	results, err := utils.ReadConfig(projectPath, keys)
	if err != nil {
		fmt.Printf("Error reading config: %v\n", err)
		os.Exit(1)
	}
	if len(results) > 0 {
		depMap, ok := results[0].(map[string]interface{})
		if !ok {
			fmt.Println("Error: dependencies are not in the expected format")
			os.Exit(1)
		}

		for pkg := range depMap {
			packs = append(packs, pkg)
			installPackage(pkg, projectPath, cachePath)
		}
		updateSpecificPackages(packs, cachePath, packagePath)
	}

}

func updateSpecificPackages(packages []string, cachePath, packagePath string) {
	for _, pkg := range packages {
		pkg = strings.Split(pkg, "::")[0]
		sourcePath := filepath.Join(cachePath, pkg)
		destPath := filepath.Join(packagePath, pkg)

		if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
			fmt.Printf("Package %s not found in cache\n", pkg)
			continue
		}

		if err := os.RemoveAll(destPath); err != nil {
			fmt.Printf("Error removing package %s: %v\n", pkg, err)
			continue
		}

		if err := os.Rename(sourcePath, destPath); err != nil {
			fmt.Printf("Error updating package %s: %v\n", pkg, err)
		} else if verbose {
			fmt.Printf("Updated package: %s\n", pkg)
		}
	}
}

func init() {
	updateCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	updateCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(updateCmd)
}
