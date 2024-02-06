package e2e

import (
	"flag"
	flagd "github.com/open-feature/go-sdk-contrib/providers/flagd/pkg"
	"testing"

	"github.com/cucumber/godog"
	"github.com/open-feature/go-sdk-contrib/tests/flagd/pkg/integration"
	"github.com/open-feature/go-sdk/openfeature"
)

func TestFlagdReconnectInRPC(t *testing.T) {
	if testing.Short() {
		// skip e2e if testing -short
		t.Skip()
	}

	flag.Parse()

	name := "flagd-reconnect.feature"

	testSuite := godog.TestSuite{
		Name: name,
		TestSuiteInitializer: integration.InitializeFlagdReconnectTestSuite(func() openfeature.FeatureProvider {
			return flagd.NewProvider(flagd.WithPort(8014))
		}),
		ScenarioInitializer: integration.InitializeFlagdReconnectScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../../flagd-testbed/gherkin/flagd-reconnect.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run reconnect tests")
	}
}

func TestFlagdReconnectInProcess(t *testing.T) {
	if testing.Short() {
		// skip e2e if testing -short
		t.Skip()
	}

	flag.Parse()

	name := "flagd-reconnect.feature"

	testSuite := godog.TestSuite{
		Name: name,
		TestSuiteInitializer: integration.InitializeFlagdReconnectTestSuite(func() openfeature.FeatureProvider {
			return flagd.NewProvider(flagd.WithInProcessResolver(), flagd.WithPort(9091))
		}),
		ScenarioInitializer: integration.InitializeFlagdReconnectScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"../../../flagd-testbed/gherkin/flagd-reconnect.feature"},
			TestingT: t, // Testing instance that will run subtests.
			Strict:   true,
		},
	}

	if testSuite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run reconnect tests")
	}
}
