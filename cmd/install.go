package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"nep/configs"
	"nep/utils"

	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

var (
	asynchronous bool
	mu           sync.Mutex
)

var installCmd = &cobra.Command{
	Use:     "install [packages]",
	Aliases: []string{"i"},
	Short:   "Install packages",
	Run: func(cmd *cobra.Command, args []string) {
		if err := changeDirectory(); err != nil {
			exitWithError(err)
		}

		projectPath, err := utils.FindProjectDir()

		if projectPath == "" {
			fmt.Println("Not a nep project")
			utils.CreateProject(configs.DefaultName, true)
			fmt.Println("Project initialized in current directory")
			projectPath, err = os.Getwd()
			if err != nil {
				fmt.Println("Error getting current directory:", err)
				os.Exit(1)
			}
		} else if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Ensure FolderName exists within the project directory
		folderPath := filepath.Join(projectPath, configs.FolderName)
		if _, err := os.Stat(folderPath); os.IsNotExist(err) {
			if err := os.Mkdir(folderPath, 0755); err != nil {
				fmt.Printf("Failed to create directory %s: %s\n", folderPath, err)
				os.Exit(1)
			}
		} else if err != nil {
			fmt.Printf("Error checking directory %s: %s\n", folderPath, err)
			os.Exit(1)
		}

		if len(args) == 0 {
			if verbose {
				fmt.Println("No arguments provided, checking dependencies...")
			}

			// Use ReadConfig to get the dependencies
			results, err := utils.ReadConfig(projectPath, [][]string{{"dependencies"}})
			if err != nil {
				fmt.Printf("Error: Failed to read dependencies from config: %v\n", err)
				os.Exit(1)
			}

			if len(results) == 0 || results[0] == nil {
				if verbose {
					fmt.Println("No dependencies found in config")
				}
				return
			}

			// Assert and convert the result to map[string]interface{}
			dependenciesInterface, ok := results[0].(map[string]interface{})
			if !ok {
				fmt.Println("Error: Dependencies in config are not in the expected format")
				os.Exit(1)
			}

			if verbose {
				fmt.Printf("Dependencies found in config: %v\n", dependenciesInterface)
			}

			// Convert to args format
			for pkg, version := range dependenciesInterface {
				versionStr, ok := version.(string)
				if !ok {
					fmt.Printf("Warning: Invalid version format for dependency %s: %v\n", pkg, version)
					continue
				}
				// Remove 'v' prefix if it exists
				versionStr = strings.TrimPrefix(versionStr, "v")
				newArg := fmt.Sprintf("%s::%s", pkg, versionStr)
				args = append(args, newArg)
				if verbose {
					fmt.Printf("Added dependency to install: %s\n", newArg)
				}
			}

			if verbose {
				fmt.Printf("Dependencies to install: %v\n", args)
			}

			if len(args) == 0 {
				fmt.Println("No valid dependencies found to install")
				return
			}
		}

		if asynchronous {
			var wg sync.WaitGroup
			for _, pkg := range args {
				wg.Add(1)
				go func(pkg string) {
					defer wg.Done()
					installPackage(pkg, projectPath, folderPath)
				}(pkg)
			}
			wg.Wait()
		} else {
			for _, pkg := range args {
				installPackage(pkg, projectPath, folderPath)
			}
		}
	},
}

func installPackage(pkg, projectPath, folderPath string) {
	// Fetch data from API
	responseData, err := utils.FetchPackageData(pkg)
	if err != nil {
		fmt.Printf("Failed to fetch data from API for %s: %s\n", pkg, err)
		return
	}

	// Extract package name and version
	parts := strings.Split(pkg, "::")
	packageName := parts[0]

	// Create a new directory for the package inside FolderName
	packageDir := filepath.Join(folderPath, packageName)

	// Ensure the directory exists or create it
	if err := os.MkdirAll(packageDir, 0755); err != nil {
		fmt.Printf("Failed to create directory %s: %s\n", packageDir, err)
		return
	}

	// Clone GitHub repository into the package directory
	options := &git.CloneOptions{
		URL:      responseData.Data.GithubURL,
		Progress: os.Stdout,
	}
	if !verbose {
		options.Progress = nil
	}

	_, err = git.PlainClone(packageDir, false, options)
	if err != nil {
		fmt.Printf("Failed to clone %s: %s\n", pkg, err)
		return
	}

	// Save API response to file inside the package directory
	if err := utils.SaveResponseToFile(responseData, packageDir); err != nil {
		fmt.Printf("Failed to save API response JSON for %s: %s\n", pkg, err)
	}

	fmt.Printf("Successfully cloned %s into %s\n", pkg, packageDir)

	name := ""

	if responseData.Key == "" {
		name = packageName
	} else {

		name = responseData.Key
	}

	updates := []utils.UpdatePath{
		{Path: []string{"dependencies", name}, Value: responseData.Data.Version},
	}

	// Use mutex to lock access to the config file during updates
	mu.Lock()
	defer mu.Unlock()

	err = utils.UpdateConfig(projectPath, updates)
	if err != nil {
		fmt.Println("Error updating config:", err)
	}
}

func init() {
	installCmd.Flags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	installCmd.Flags().BoolVarP(&asynchronous, "asynchronous", "a", false, "Install packages in parallel")
	installCmd.Flags().StringVarP(&path, "path", "p", "", "Set project path")
	rootCmd.AddCommand(installCmd)
}
