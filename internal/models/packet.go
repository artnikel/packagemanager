// Package models defines data structures used throughout the package manager system
package models

import "encoding/json"

// PacketFile package description structure
type PacketFile struct {
	Name    string   `json:"name"`
	Version string   `json:"ver"`
	Targets []Target `json:"targets"`
	Packets []Packet `json:"packets,omitempty"`
}

// Target structure for describing target files
type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

// UnmarshalJSON to support both string and object format in targets
func (t *Target) UnmarshalJSON(data []byte) error {
	var path string
	if err := json.Unmarshal(data, &path); err == nil {
		t.Path = path
		return nil
	}

	type target Target
	var temp target
	if err := json.Unmarshal(data, &temp); err != nil {
		return err
	}
	*t = Target(temp)
	return nil
}

// Packet structure for package dependencies
type Packet struct {
	Name    string `json:"name"`
	Version string `json:"ver,omitempty"`
}

// PackagesFile structure for the packages.json file
type PackagesFile struct {
	Packages []Packet `json:"packages"`
}
