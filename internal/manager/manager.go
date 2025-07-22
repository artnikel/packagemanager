// Package manager handles creation and installation of packages from JSON configs
package manager

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/artnikel/packagemanager/internal/archive"
	"github.com/artnikel/packagemanager/internal/config"
	"github.com/artnikel/packagemanager/internal/models"
	"github.com/artnikel/packagemanager/internal/remote"
)

// CreatePacket creates an archive based on packet.json
func CreatePacket(configPath string, cfg *config.Config) error {
	fmt.Printf("Creating a packet of %s...\n", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var packet models.PacketFile
	if err := json.Unmarshal(data, &packet); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	archiveName := fmt.Sprintf("%s-%s.tar.gz", packet.Name, packet.Version)
	if err := archive.CreateArchive(archiveName, packet.Targets); err != nil {
		return fmt.Errorf("error archive creation: %v", err)
	}

	fmt.Printf("Archive %s successfully created\n", archiveName)

	if err := remote.UploadToServer(archiveName, cfg); err != nil {
		fmt.Printf("Warning: failed to upload to server: %v\n", err)
	}

	return nil
}

// UpdatePackages downloads and installs packages from packages.json
func UpdatePackages(configPath string, cfg *config.Config) error {
	fmt.Printf("Updating packages from %s...\n", configPath)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var packages models.PackagesFile
	if err := json.Unmarshal(data, &packages); err != nil {
		return fmt.Errorf("error parsing JSON: %v", err)
	}

	for _, pkg := range packages.Packages {
		fmt.Printf("Package update: %s %s\n", pkg.Name, pkg.Version)

		version := "1.0"
		if pkg.Version != "" {
			version = strings.TrimLeft(pkg.Version, "><=")
		}

		archiveName := fmt.Sprintf("%s-%s.tar.gz", pkg.Name, version)

		if err := remote.DownloadFromServer(archiveName, cfg); err != nil {
			fmt.Printf("Warning: failed to download %s: %v\n", archiveName, err)
			continue
		}

		if err := archive.ExtractArchive(archiveName, "./packages/"+pkg.Name); err != nil {
			fmt.Printf("Unpacking error %s: %v\n", archiveName, err)
			continue
		}

		fmt.Printf("Package %s successfully installed\n", pkg.Name)
	}

	return nil
}
