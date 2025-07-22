package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type SSHConfig struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password,omitempty"`
	KeyPath  string `yaml:"key_path,omitempty"`
}

type PathConfig struct {
	RemotePackageDir string `yaml:"remote_package_dir"` 
	LocalPackageDir  string `yaml:"local_package_dir"`  
}

type Config struct {
	SSH  SSHConfig  `yaml:"ssh"`
	Path PathConfig `yaml:"path"`
}

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
