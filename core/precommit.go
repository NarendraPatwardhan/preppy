package core

import (
	"embed"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

//go:embed known-precommits.json
var p embed.FS

type preCommit struct {
	Repos []struct {
		Repo  string `yaml:"repo"`
		Rev   string `yaml:"rev"`
		Hooks []struct {
			ID string `yaml:"id"`
		} `yaml:"hooks"`
	} `yaml:"repos"`
}

func onPyPi(arg string) bool {
	url := "https://pypi.python.org/pypi/" + arg
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	if resp.StatusCode == 200 {
		return true
	}
	return false
}

func getKnownPrecommits() (map[string]bool, error) {
	knownPrecommits := make(map[string]bool)
	preCommits, err := p.ReadFile("known-precommits.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(preCommits, &knownPrecommits)
	if err != nil {
		return nil, err
	}
	return knownPrecommits, nil
}

func GetDevDependencies(path string) ([]string, error) {
	precommitPath := path + "/.pre-commit-config.yaml"
	// Check if pre-commit file exists
	if _, err := os.Stat(precommitPath); os.IsNotExist(err) {
		return nil, nil
	}
	source, err := ioutil.ReadFile(precommitPath)
	if err != nil {
		return nil, err
	}
	// Parse the yaml file.
	var data preCommit
	err = yaml.Unmarshal(source, &data)
	if err != nil {
		return nil, err
	}
	knownPrecommits, err := getKnownPrecommits()
	if err != nil {
		return nil, err
	}
	parsedPrecommits := []string{}
	// Loop through the repos.
	for _, repo := range data.Repos {
		// Ignore if repo is original precommit repo
		if repo.Repo == "https://github.com/pre-commit/pre-commit-hooks" {
			continue
		}
		// Loop through the hooks.
		for _, hook := range repo.Hooks {
			// Check if the hook is known.
			if knownPrecommits[hook.ID] {
				// Append to parsedPrecommits
				parsedPrecommits = append(parsedPrecommits, hook.ID)
			} else {
				// Check if the hook is on PyPi
				if onPyPi(hook.ID) {
					// Append to parsedPrecommits
					parsedPrecommits = append(parsedPrecommits, hook.ID)
				}
			}
		}
	}
	return parsedPrecommits, nil
}
