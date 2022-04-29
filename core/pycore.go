package core

import (
	"embed"
	"encoding/json"
)

//go:embed stdlib.json
var f embed.FS

func GetSystemPackages() (map[string]bool, error) {
	// Read stdlib.json
	stdlib, err := f.ReadFile("stdlib.json")
	if err != nil {
		return nil, err
	}
	// Parse stdlib.json into a map
	stdlibMap := make(map[string]bool)
	json.Unmarshal(stdlib, &stdlibMap)
	return stdlibMap, nil
}
