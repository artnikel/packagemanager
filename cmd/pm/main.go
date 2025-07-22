// Package main is an entry point to application
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/artnikel/packagemanager/internal/config"
	"github.com/artnikel/packagemanager/internal/manager"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	const minArgs = 3
	if len(os.Args) < minArgs {
		fmt.Println("Usage:")
		fmt.Println("  pm create ./packet.json")
		fmt.Println("  pm update ./packages.json")
		os.Exit(1)
	}

	command := os.Args[1]
	configPath := os.Args[2]

	switch command {
	case "create":
		if err := manager.CreatePacket(configPath, cfg); err != nil {
			log.Fatalf("Error creating a packet: %v", err)
		}
	case "update":
		if err := manager.UpdatePackages(configPath, cfg); err != nil {
			log.Fatalf("Package update error: %v", err)
		}
	default:
		fmt.Printf("Unknown command: %s\n", command)
		os.Exit(1)
	}
}
