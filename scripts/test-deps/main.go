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
// from it and prints it to stdout as a JSON array that can be used in GitHub Actions matrix.
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

		vs, ok := versions.([]any)
		if !ok {
			exitWithErr(fmt.Errorf("dependency %s.%s is of unsupported type: %T", testLevel, dependency, versions))
		}

		// If only the latest, use only the latest version...
		if latest {
			vs = []any{vs[len(vs)-1]}
		}

		// Print the versions as a JSON array.
		err = json.NewEncoder(os.Stdout).Encode(vs)
		exitIfErr(err)
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
