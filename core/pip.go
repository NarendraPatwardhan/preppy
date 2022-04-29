package core

import (
	"encoding/json"
	"os/exec"
)

type pkg struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

func PipList() (map[string]string, error) {
	// Execute pip list
	cmd := exec.Command("pip3", "list", "--format", "json")
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}
	// Parse pip list output
	pkgList := []pkg{}
	json.Unmarshal(out, &pkgList)
	// Create map of packages
	pkgs := make(map[string]string)
	for _, pkg := range pkgList {
		pkgs[pkg.Name] = pkg.Version
	}
	return pkgs, nil
}
