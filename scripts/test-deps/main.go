package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"sigs.k8s.io/yaml"
)

// DependenciesFile is a map of test levels and their dependencies.
type DependenciesFile map[string]Dependencies

// Dependencies is a map of dependencies and their versions OR
// a slice of matrix entries that can be used in GitHub Actions matrix directly.
type Dependencies any

// This program reads `.github/test_dependencies.yaml` file, extracts a requested dependency's versions
// from it and prints it to stdout as a JSON value (as an array or a single value, depending on YAML definition).
//
// Usage:
// go run scripts/test-deps/main.go <test-level> <dependency> [--latest]
func main() {
	b, err := os.ReadFile(".github/test_dependencies.yaml")
	exitIfErr(err)

	deps := DependenciesFile{}
	err = yaml.Unmarshal(b, &deps)
	exitIfErr(err)

	testLevel := ""
	flag.StringVar(&testLevel, "test-level", "", "Test level to extract dependencies for.")

	dependency := ""
	flag.StringVar(&dependency, "dependency", "", "Dependency to extract versions for.")

	latest := false
	flag.BoolVar(&latest, "latest", false, "If set, only the latest version will be used.")

	flag.Parse()

	// Extract the requested dependency.
	testLevelDeps, ok := deps[testLevel]
	if !ok {
		exitWithErr(fmt.Errorf("test level %s not found", testLevel))
	}

	switch deps := testLevelDeps.(type) {
	case []string:
		// If the test level is a slice of matrix entries, print it as a JSON array and exit.
		err = json.NewEncoder(os.Stdout).Encode(deps)
		exitIfErr(err)
		os.Exit(0)
	case map[string]interface{}:
		// If the test level is a map of dependencies and their versions, continue.
		versions, ok := deps[dependency]
		if !ok {
			exitWithErr(fmt.Errorf("dependency %s.%s not found", testLevel, dependency))
		}

		switch v := versions.(type) {
		case string:
			// If the dependency is a single version, print it as a JSON value and exit.
			err = json.NewEncoder(os.Stdout).Encode(v)
			exitIfErr(err)
			os.Exit(0)
		case []any:
			// If --latest, use only the latest version...
			if latest {
				v = []any{v[0]}
			}
			// Print the versions as a JSON array and exit.
			err = json.NewEncoder(os.Stdout).Encode(v)
			exitIfErr(err)
			os.Exit(0)
		default:
			exitWithErr(fmt.Errorf("dependency %s.%s is of unsupported type: %T", testLevel, dependency, versions))
		}

	default:
		exitWithErr(fmt.Errorf("test level %q is of unsupported type: %T", testLevel, testLevelDeps))
	}

}

func exitIfErr(err error) {
	if err != nil {
		exitWithErr(err)
	}
}

func exitWithErr(err error) {
	fmt.Printf("ERROR: %s\n", err)
	os.Exit(1)
}
