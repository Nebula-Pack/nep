package cmd

import (
	"fmt"
	"nep/configs"
	"nep/utils"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var uninstallCmd = &cobra.Command{
	Use:     "uninstall [packages]",
	Aliases: []string{"un"},
	Short:   "Uninstall packages",
	Long:    "Uninstall one or more packages, or use 'all' to uninstall all packages",
	Run:     runUninstall,
}

func runUninstall(cmd *cobra.Command, args []string) {
	if err := changeDirectory(); err != nil {
		exitWithError(err)
	}

	if len(args) == 0 {
		exitWithError(fmt.Errorf("no packages specified"))
	}

	projectPath, err := utils.FindProjectDir()
	if err != nil {
		exitWithError(err)
	}

	updates, err := getUpdates(projectPath, args)
	if err != nil {
		exitWithError(err)
	}

	if err := utils.UpdateConfig(projectPath, updates); err != nil {
		exitWithError(fmt.Errorf("error updating config: %v", err))
	}

	fmt.Println("Packages uninstalled successfully.")
}

func changeDirectory() error {
	if path != "" {
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}
		if verbose {
			fmt.Printf("Changed working directory to: %s\n", path)
		}
	}
	return nil
}

func getUpdates(projectPath string, args []string) ([]utils.UpdatePath, error) {
	var updates []utils.UpdatePath

	if args[0] == configs.All {
		return getAllUpdates(projectPath)
	}

	for _, pkg := range args {
		pkgPath := filepath.Join(projectPath, configs.FolderName, pkg)
		if err := os.RemoveAll(pkgPath); err != nil {
			fmt.Printf("Warning: Error removing package %s: %v\n", pkg, err)
		}
		updates = append(updates, utils.UpdatePath{Path: []string{"dependencies", pkg}, Value: configs.RemoveMarker})
	}

	return updates, nil
}

func getAllUpdates(projectPath string) ([]utils.UpdatePath, error) {
	keys := [][]string{{"dependencies"}}
	results, err := utils.ReadConfig(projectPath, keys)
	if err != nil {
		return nil, fmt.Errorf("error reading config: %v", err)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no dependencies found")
	}

	depMap, ok := results[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("dependencies are not in the expected format")
	}

	var updates []utils.UpdatePath
	for pkg := range depMap {
		pkgPath := filepath.Join(projectPath, configs.FolderName, pkg)
		if err := os.RemoveAll(pkgPath); err != nil {
			fmt.Printf("Warning: Error removing package %s: %v\n", pkg, err)
		}
		updates = append(updates, utils.UpdatePath{Path: []string{"dependencies", pkg}, Value: configs.RemoveMarker})
	}

	return updates, nil
}

func exitWithError(err error) {
	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
	os.Exit(1)
}

func init() {
	uninstallCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	uninstallCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(uninstallCmd)
}
