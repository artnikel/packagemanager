package models

import (
	"encoding/json"
	"testing"
)

func TestPacketFileUnmarshal(t *testing.T) {
	jsonData := `{
		"name": "test-package",
		"ver": "1.0.0",
		"targets": [
			"./src/*.go",
			{
				"path": "./docs/*",
				"exclude": "*.tmp"
			}
		],
		"packets": [
			{
				"name": "dependency1",
				"ver": ">=2.0.0"
			}
		]
	}`

	var packet PacketFile
	err := json.Unmarshal([]byte(jsonData), &packet)
	if err != nil {
		t.Fatalf("Failed to unmarshal PacketFile: %v", err)
	}

	if packet.Name != "test-package" {
		t.Errorf("Expected name 'test-package', got '%s'", packet.Name)
	}
	if packet.Version != "1.0.0" {
		t.Errorf("Expected version '1.0.0', got '%s'", packet.Version)
	}

	if len(packet.Targets) != 2 {
		t.Errorf("Expected 2 targets, got %d", len(packet.Targets))
	}

	if packet.Targets[0].Path != "./src/*.go" {
		t.Errorf("Expected path './src/*.go', got '%s'", packet.Targets[0].Path)
	}
	if packet.Targets[0].Exclude != "" {
		t.Errorf("Expected empty exclude, got '%s'", packet.Targets[0].Exclude)
	}

	if packet.Targets[1].Path != "./docs/*" {
		t.Errorf("Expected path './docs/*', got '%s'", packet.Targets[1].Path)
	}
	if packet.Targets[1].Exclude != "*.tmp" {
		t.Errorf("Expected exclude '*.tmp', got '%s'", packet.Targets[1].Exclude)
	}

	if len(packet.Packets) != 1 {
		t.Errorf("Expected 1 packet dependency, got %d", len(packet.Packets))
	}
}

func TestPackagesFileUnmarshal(t *testing.T) {
	jsonData := `{
		"packages": [
			{
				"name": "package1",
				"ver": "1.0.0"
			},
			{
				"name": "package2",
				"ver": ">=2.0.0"
			}
		]
	}`

	var packages PackagesFile
	err := json.Unmarshal([]byte(jsonData), &packages)
	if err != nil {
		t.Fatalf("Failed to unmarshal PackagesFile: %v", err)
	}

	if len(packages.Packages) != 2 {
		t.Errorf("Expected 2 packages, got %d", len(packages.Packages))
	}

	if packages.Packages[0].Name != "package1" {
		t.Errorf("Expected name 'package1', got '%s'", packages.Packages[0].Name)
	}
	if packages.Packages[1].Version != ">=2.0.0" {
		t.Errorf("Expected version '>=2.0.0', got '%s'", packages.Packages[1].Version)
	}
}
