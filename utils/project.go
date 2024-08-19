package utils

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"nep/configs"
	"os"
	"path/filepath"
)

// CreateProject creates a new project directory with nebpack folder and nebula-config.json file
func CreateProject(projectName string, useCurrentDir bool) (string, error) {
	var projectDir string
	var err error

	if useCurrentDir {
		projectDir, err = os.Getwd()
		if err != nil {
			return "", fmt.Errorf("failed to get current directory: %v", err)
		}
	} else {
		if projectName == "" {
			return "", fmt.Errorf("project name cannot be empty")
		}
		projectDir = filepath.Join(".", projectName)
		err = os.Mkdir(projectDir, 0755)
		if err != nil {
			return "", fmt.Errorf("failed to create project directory: %v", err)
		}
	}

	// Create nebpack directory
	nebpackDir := filepath.Join(projectDir, configs.FolderName)
	err = os.Mkdir(nebpackDir, 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create %s directory: %v", configs.JSONName, err)
	}

	// Create nebula-config.json file
	configFilePath := filepath.Join(projectDir, configs.JSONName+".json")
	err = os.WriteFile(configFilePath, configs.ConfigFileBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("failed to write to nebula-config.json: %v", err)
	}

	return projectDir, nil
}

// UpdatePath represents a path in the JSON structure and the value to set
type UpdatePath struct {
	Path  []string
	Value interface{}
}

// UpdateConfig updates the nebula-config.json file with the given key-value pairs.
func UpdateConfig(projectDir string, updates []UpdatePath) error {
	configFilePath := filepath.Join(projectDir, configs.JSONName+".json")

	// Read the existing config file
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the JSON content
	var config map[string]interface{}
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return fmt.Errorf("failed to parse config file: %v", err)
	}

	// Update the config with the new values
	for _, update := range updates {
		nestedUpdate(config, update.Path, update.Value)
	}

	// Convert the updated config back to JSON
	updatedConfigBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal updated config: %v", err)
	}

	// Write the updated config back to the file
	err = os.WriteFile(configFilePath, updatedConfigBytes, 0644)
	if err != nil {
		return fmt.Errorf("failed to write updated config file: %v", err)
	}

	return nil
}

func nestedUpdate(config map[string]interface{}, keys []string, value interface{}) {
	lastKey := keys[len(keys)-1]
	m := config

	for _, k := range keys[:len(keys)-1] {
		if _, ok := m[k].(map[string]interface{}); !ok {
			m[k] = make(map[string]interface{})
		}
		m = m[k].(map[string]interface{})
	}

	if value == configs.RemoveMarker {
		delete(m, lastKey)
	} else {
		m[lastKey] = value
	}
}

// FindProjectDir checks if the current directory or any parent directory is a project directory
func FindProjectDir() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("failed to get current directory: %v", err)
	}

	for {
		configFilePath := filepath.Join(dir, configs.JSONName+".json")
		if _, err := os.Stat(configFilePath); err == nil {
			return dir, nil
		}

		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			return "", fmt.Errorf("no project directory found")
		}
		dir = parentDir
	}
}

// ReadConfig reads values from the nebula-config.json file based on the given paths.
func ReadConfig(projectDir string, paths [][]string) ([]interface{}, error) {
	configFilePath := filepath.Join(projectDir, configs.JSONName+".json")

	// Read the existing config file
	configFileBytes, err := os.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	// Parse the JSON content
	var config map[string]interface{}
	err = json.Unmarshal(configFileBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Read the values based on the provided paths
	results := make([]interface{}, len(paths))
	for i, path := range paths {
		value, err := nestedRead(config, path)
		if err != nil {
			return nil, fmt.Errorf("failed to read path %v: %v", path, err)
		}
		results[i] = value
	}

	return results, nil
}

// nestedRead reads a nested value from a map based on the given keys
func nestedRead(config map[string]interface{}, keys []string) (interface{}, error) {
	var current interface{} = config

	for _, key := range keys {
		switch v := current.(type) {
		case map[string]interface{}:
			var ok bool
			current, ok = v[key]
			if !ok {
				return nil, fmt.Errorf("key not found: %s", key)
			}
		default:
			return nil, fmt.Errorf("invalid path: %v", keys)
		}
	}

	return current, nil
}

func ChangeDirectory(path string) error {
	if path != "" {
		if err := os.Chdir(path); err != nil {
			return fmt.Errorf("error changing directory: %v", err)
		}
	}
	return nil
}

func Prepare(create bool, path string) (projectPath string) {

	// Change to the specified directory
	if err := ChangeDirectory(path); err != nil {
		log.Fatal(err)
	}

	// Find the project directory
	var err error
	projectPath, err = FindProjectDir()
	if err != nil {
		log.Fatal(err)
	}

	if projectPath == "" {
		if !create {
			log.Fatal("Not a nep project")
		}
		fmt.Println("Not a nep project\nInitializing project...")
		CreateProject(configs.DefaultName, true)
		fmt.Println("Project initialized in current directory")

		// Get the current directory after creating the project
		projectPath, err = os.Getwd()
		if err != nil {
			log.Fatal("Error getting current directory:", err)
		}
	}

	return
}

func GetFolder(projectPath string) string {

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

	return folderPath
}
