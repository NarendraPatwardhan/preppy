package core

import (
	"io/ioutil"
	"strings"
)

type preIdentified struct {
	Operator string
	Version  string
}

func ReadExistingRequirements(path string) (map[string]preIdentified, string, error) {
	result := make(map[string]preIdentified)
	// Read requirements.txt present in the path
	requirements, err := ioutil.ReadFile(path + "/requirements.txt")
	if err != nil {
		return result, "", err
	}
	// Split the requirements.txt into lines
	lines := strings.Split(string(requirements), "\n")
	// Iterate over the lines
	for _, line := range lines {
		// Skip empty lines
		if len(line) == 0 {
			continue
		}
		// Skip comments
		if strings.HasPrefix(line, "#") {
			continue
		}
		// Split the line again using # to omit inline comments
		pkgStr := strings.Split(line, "#")[0]
		appended := false
		// Try to split the line into operator and version
		for _, operator := range []string{"==", ">=", ">"} {
			pkgInfo := strings.Split(pkgStr, operator)
			if len(pkgInfo) == 2 {
				result[pkgInfo[0]] = preIdentified{operator, pkgInfo[1]}
				appended = true
				break
			}
		}
		if !appended {
			// Strip whitespace
			pkgStr = strings.TrimSpace(pkgStr)
			result[pkgStr] = preIdentified{"none", "none"}
		}
	}
	return result, string(requirements), nil
}
