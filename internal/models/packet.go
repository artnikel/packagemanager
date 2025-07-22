package models

import "encoding/json"

type PacketFile struct {
	Name    string   `json:"name"`
	Version string   `json:"ver"`
	Targets []Target `json:"targets"`
	Packets []Packet `json:"packets,omitempty"`
}

type Target struct {
	Path    string `json:"path"`
	Exclude string `json:"exclude,omitempty"`
}

type Packet struct {
	Name    string `json:"name"`
	Version string `json:"ver,omitempty"`
}

type PackagesFile struct {
	Packages []Packet `json:"packages"`
}

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
