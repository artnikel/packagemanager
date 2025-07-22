// Package config provides configuration loading from YAML files
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// SSHConfig holds SSH connection settings
type SSHConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password,omitempty"`
	KeyPath  string `yaml:"key_path,omitempty"`
}

// PathConfig holds file path configuration
type PathConfig struct {
	RemotePackageDir string `yaml:"remote_package_dir"`
	LocalPackageDir  string `yaml:"local_package_dir"`
}

// Config aggregates all application configurations
type Config struct {
	SSH  SSHConfig  `yaml:"ssh"`
	Path PathConfig `yaml:"path"`
}

// LoadConfig loads the configuration from the given YAML file path
func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
