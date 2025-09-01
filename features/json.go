package features

import (
	"encoding/json"
	"errors"
	"os"
	"strings"
)

// JSONConfig represents a configuration structure
type JSONConfig struct {
	Name     string   `json:"name"`
	Version  string   `json:"version"`
	Port     int      `json:"port"`
	Enabled  bool     `json:"enabled"`
	Features []string `json:"features"`
}

func ReadConfigFromFile(filename string) (*JSONConfig, error) {
	jsonText, err := os.ReadFile(filename)
	if err != nil {
		return nil, errors.New("no file was found")
	}
	config := JSONConfig{}
	err = json.Unmarshal(jsonText, &config)
	return &config, err
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !errors.Is(err, os.ErrNotExist)
}

func WriteToConfig(configPath string, param string, value any) {
	if !fileExists(configPath) {
		createDefaultConfig(configPath)
	}
	config, err := ReadConfigFromFile(configPath)
	if err != nil {
		panic(err)
	}
	switch strings.ToLower(param) {
	case "name":
		if v, ok := value.(string); ok {
			config.Name = v
		}
	case "version":
		if v, ok := value.(string); ok {
			config.Version = v
		}
	case "port":
		if v, ok := value.(int); ok {
			config.Port = v
		}
	case "enabled":
		if v, ok := value.(bool); ok {
			config.Enabled = v
		}
	case "features":
		if v, ok := value.([]string); ok {
			config.Features = v
		}
	default:
		panic("Did not expect parameter: " + param)
	}

	jsonEncoding, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, jsonEncoding, 0644)
	if err != nil {
		panic(err)
	}
}

func createDefaultConfig(configPath string) {
	defaultConfig := JSONConfig{
		Name:     "defaultName",
		Version:  "0.0",
		Port:     3000,
		Enabled:  false,
		Features: []string{""},
	}
	jsonEncoding, err := json.Marshal(defaultConfig)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(configPath, jsonEncoding, 0644)
	if err != nil {
		panic(err)
	}
}
