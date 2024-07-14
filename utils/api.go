// utils/api.go

package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"nep/configs"
	"net/http"
	"path/filepath"
	"strings"
)

// Response represents the structure of the API response.
type Response struct {
	Data              Data              `json:"data"`
	Key               string            `json:"key"`
	TemporalSemantics TemporalSemantics `json:"temporal_semantics"`
}

// Data represents the data part of the API response.
type Data struct {
	GithubURL    string `json:"github_url"`
	HasRockspec  bool   `json:"hasRockspec"`
	IsLua        bool   `json:"isLua"`
	ScanResponse struct {
		Lua string `json:"lua"`
	} `json:"scanResponse"`
	Version string `json:"version"`
}

// TemporalSemantics represents the temporal semantics part of the API response.
type TemporalSemantics struct {
	LatestGetRequest string `json:"latest-get-request"`
}

// FetchPackageData fetches data from the API for a given package.
func FetchPackageData(packageName string) (*Response, error) {
	var apiUrl string
	parts := strings.Split(packageName, "::")
	if len(parts) == 2 {
		apiUrl = fmt.Sprintf("%s/api/%s:v%s", configs.APIBaseURL, parts[0], parts[1])
	} else {
		apiUrl = fmt.Sprintf("%s/api/%s", configs.APIBaseURL, packageName)
	}

	resp, err := http.Get(apiUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch data from API for %s: %s", packageName, err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read API response body for %s: %s", packageName, err)
	}

	var responseData Response
	if err := json.Unmarshal(body, &responseData); err != nil {
		return nil, fmt.Errorf("failed to parse JSON for %s: %s", packageName, err)
	}

	return &responseData, nil
}

// SaveResponseToFile saves API response JSON to a file with proper indentation.
func SaveResponseToFile(responseData *Response, clonePath string) error {
	responseJSON, err := json.MarshalIndent(responseData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to serialize API response: %s", err)
	}

	responseFilePath := filepath.Join(clonePath, configs.ResponseFileName+".json")
	if err := ioutil.WriteFile(responseFilePath, responseJSON, 0644); err != nil {
		return fmt.Errorf("failed to save API response JSON: %s", err)
	}

	return nil
}
