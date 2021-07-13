package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	SeverityOverrides map[string]string `json:"severity_overrides,omitempty" yaml:"severity_overrides,omitempty"`
	ExcludedChecks    []string          `json:"exclude,omitempty" yaml:"exclude,omitempty"`
}

func LoadConfig(configFilePath string) (*Config, error) {
	var config = &Config{}

	if _, err := os.Stat(configFilePath); err != nil {
		return nil, fmt.Errorf("failed to access config file '%s': %s", configFilePath, err)
	}

	configFileContent, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file '%s': %s", configFilePath, err)
	}

	ext := filepath.Ext(configFilePath)
	switch strings.ToLower(ext) {
	case ".json":
		err = json.Unmarshal(configFileContent, config)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file '%s': %s", configFilePath, err)
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal(configFileContent, config)
		if err != nil {
			return nil, fmt.Errorf("failed to load config file '%s': %s", configFilePath, err)
		}
	default:
		return nil, fmt.Errorf("couldn't process the file %s", configFilePath)
	}

	rewriteSeverityOverrides(config)

	return config, nil
}

func rewriteSeverityOverrides(config *Config) error {

	for k, s := range config.SeverityOverrides {
		switch strings.ToUpper(s) {
		case "CRITICAL", "HIGH", "MEDIUM", "LOW":
			continue
		case "ERROR":
			config.SeverityOverrides[k] = "HIGH"
		case "WARNING":
			config.SeverityOverrides[k] = "MEDIUM"
		case "INFO":
			config.SeverityOverrides[k] = "LOW"
		default:
			return fmt.Errorf("could not rewrite the severity code [%s]", s)
		}
	}

	return nil
}
