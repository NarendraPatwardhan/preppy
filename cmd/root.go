package cmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"preppy/core"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "preppy",
	Short: "A command line tool to manage prerequisites for your python project",
	Run: func(cmd *cobra.Command, args []string) {
		root := cmd.Flag("root").Value.String()
		if root == "" {
			root, _ = os.Getwd()
		}
		dryRun := cmd.Flag("dry-run").Value.String() == "true"
		paths, folders, err := core.GetRecursivePaths(root)
		if err != nil {
			fmt.Println("Error walking the directory tree")
			return
		}
		if len(paths) == 0 {
			fmt.Println("No python files found")
			return
		}
		var wg sync.WaitGroup
		// Create a channel to receive the results.
		importCh := make(chan []string, len(paths))
		// Create a channel to receive the errors.
		errCh := make(chan error, len(paths))
		for _, path := range paths {
			wg.Add(1)
			go func(path string) {
				defer wg.Done()
				// Read file
				source, err := ioutil.ReadFile(path)
				if err != nil {
					errCh <- err
					return
				}
				// Get imports
				imports, err := core.GetImports(source)
				if err != nil {
					errCh <- err
					return
				}
				// Send result to channel
				importCh <- imports
			}(path)
		}
		// Create a channel to receive the pip packages.
		preInstalledCh := make(chan map[string]string, 1)
		// Create a channel to receive the errors.
		preInstalledErrCh := make(chan error, 1)
		wg.Add(1)
		go func() {
			defer wg.Done()
			preInstalled, err := core.PipList()
			if err != nil {
				preInstalledErrCh <- err
				return
			}
			preInstalledCh <- preInstalled
		}()
		// Wait for all goroutines to finish.
		wg.Wait()
		// Close the channels.
		close(importCh)
		close(errCh)
		close(preInstalledCh)
		close(preInstalledErrCh)
		// Collect the results.
		stdLib, err := core.GetSystemPackages()
		if err != nil {
			fmt.Println("Error loading system packages")
			return
		}
		err = <-errCh
		if err != nil {
			fmt.Println("Error parsing file(s)")
			return
		}
		err = <-preInstalledErrCh
		if err != nil {
			fmt.Println("Error loading pip packages")
			return
		}
		preInstalled := <-preInstalledCh
		uniqueImports := make(map[string]bool)
		for imports := range importCh {
			for _, import_ := range imports {
				uniqueImports[import_] = true
			}
		}

		var requirements string
		var versioned []string
		var unversioned []string
		for import_ := range uniqueImports {
			// Rule 1: If import is in stdlib, skip
			if _, ok := stdLib[import_]; ok {
				continue
			}
			// Rule 2: If import is in folders and not in pip, skip
			version, installed := preInstalled[import_]
			_, ok := folders[import_]
			if ok && !installed {
				continue
			}
			// Rule 3: If import is in pip, add to requirements with version
			if installed {
				// Split version by '.'
				versionSplit := strings.Split(version, ".")
				// If length is less than 3 or last element is not a number, pin exact version
				last := versionSplit[len(versionSplit)-1]
				// Try parsing last element as int
				_, err = strconv.Atoi(last)
				if len(versionSplit) != 3 || err != nil {
					versioned = append(versioned, fmt.Sprintf("%s==%s", import_, version))
				} else {
					versioned = append(versioned, fmt.Sprintf("%s>=%s", import_, version))
				}
			} else {
				// If not in pip, add to requirements without version
				unversioned = append(unversioned, import_)
			}
		}
		// Sort versioned
		sort.Strings(versioned)
		// Sort unversioned
		sort.Strings(unversioned)
		// Join versioned and unversioned
		requirements = strings.Join(versioned, "\n")
		requirements += "\n"
		requirements += strings.Join(unversioned, "\n")
		if dryRun {
			fmt.Println(requirements)
		} else {
			err = ioutil.WriteFile("./requirements.txt", []byte(requirements), 0644)
			if err != nil {
				fmt.Println("Error writing requirements.txt")
				return
			}
			fmt.Println("Successfully wrote requirements.txt")
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("root", "r", "", "root directory")
	rootCmd.Flags().BoolP("dry-run", "d", false, "preview without writing to disk")
}
