package manager

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/artnikel/packagemanager/internal/config"
	"github.com/artnikel/packagemanager/internal/models"
)

func TestCreatePacket(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "manager_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testFile := filepath.Join(tmpDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0o644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	packet := models.PacketFile{
		Name:    "test-package",
		Version: "1.0.0",
		Targets: []models.Target{
			{Path: testFile},
		},
	}

	packetData, err := json.Marshal(packet)
	if err != nil {
		t.Fatalf("Failed to marshal packet: %v", err)
	}

	packetPath := filepath.Join(tmpDir, "packet.json")
	if err := os.WriteFile(packetPath, packetData, 0o644); err != nil {
		t.Fatalf("Failed to write packet file: %v", err)
	}

	cfg := &config.Config{
		SSH: config.SSHConfig{
			Host: "localhost",
			Port: "22",
			User: "test",
		},
		Path: config.PathConfig{
			RemotePackageDir: "/tmp/remote",
			LocalPackageDir:  "/tmp/local",
		},
	}

	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	err = CreatePacket(packetPath, cfg)
	if err != nil {
		t.Logf("Expected upload error: %v", err)
	}

	expectedArchive := "test-package-1.0.0.tar.gz"
	if _, err := os.Stat(expectedArchive); os.IsNotExist(err) {
		t.Error("Archive was not created")
	}
}

func TestUpdatePackages(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "update_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	packages := models.PackagesFile{
		Packages: []models.Packet{
			{Name: "test-package", Version: "1.0.0"},
		},
	}

	packagesData, err := json.Marshal(packages)
	if err != nil {
		t.Fatalf("Failed to marshal packages: %v", err)
	}

	packagesPath := filepath.Join(tmpDir, "packages.json")
	if err := os.WriteFile(packagesPath, packagesData, 0o644); err != nil {
		t.Fatalf("Failed to write packages file: %v", err)
	}

	cfg := &config.Config{
		SSH: config.SSHConfig{
			Host: "localhost",
			Port: "22",
			User: "test",
		},
		Path: config.PathConfig{
			RemotePackageDir: "/tmp/remote",
			LocalPackageDir:  "/tmp/local",
		},
	}

	originalDir, _ := os.Getwd()
	os.Chdir(tmpDir)
	defer os.Chdir(originalDir)

	err = UpdatePackages(packagesPath, cfg)
	if err != nil {
		t.Logf("Expected download error: %v", err)
	}
}
