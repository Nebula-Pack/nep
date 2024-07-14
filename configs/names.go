package configs

import (
	_ "embed"
)

const (
	JSONName         string = "nebula-config"
	FolderName       string = "nebpack"
	DefaultName      string = "Nebula-Pack-Project"
	ResponseFileName string = "nebula-config"
	RemoveMarker     string = "__REMOVE__"
	All              string = "*"
)

//go:embed project_config.json
var ConfigFileBytes []byte
