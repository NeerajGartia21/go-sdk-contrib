package integration

import (
	"context"
	"fmt"
	"github.com/cucumber/godog"
	"github.com/open-feature/go-sdk/openfeature"
	"time"
)

var readyHandlerRunCount = 0
var errorHandlerRunCount = 0

// InitializeFlagdReconnectTestSuite register provider supplier and register test steps
func InitializeFlagdReconnectTestSuite(providerSupplier func() openfeature.FeatureProvider) func(*godog.TestSuiteContext) {
	test_provider_supplier = providerSupplier

	return func(suiteContext *godog.TestSuiteContext) {}
}

// InitializeFlagdReconnectScenario initializes the flagd reconnect test scenario
func InitializeFlagdReconnectScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a flagd provider is set$`, aFlagdProviderIsSet)
	ctx.Step(`^a PROVIDER_READY handler and a PROVIDER_ERROR handler are added$`, aProviderReadyAndProviderErrorHandlerAdded)
	ctx.Step(`^the PROVIDER_READY handler must run when the provider connects$`, aProviderReadyHandlerMustRunWhenConnectionEstablished)
	ctx.Step(`^the PROVIDER_ERROR handler must run when the provider's connection is lost$`, aProviderErrorHandlerMustRunWhenConnectionLost)
	ctx.Step(`^when the connection is reestablished the PROVIDER_READY handler must run again$`, aProviderReadyHandlerMustRunAgainWhenConnectionReestablished)
}

func waitUntilTimeOut(condition func() bool, interval time.Duration) bool {
	ctx, cancel := context.WithTimeout(context.Background(), interval)
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			return false
		default:
			if condition() {
				return true
			}
		}
	}
}

func aProviderReadyAndProviderErrorHandlerAdded(ctx context.Context) (context.Context, error) {
	client := ctx.Value(ctxClientKey{}).(*openfeature.Client)
	providerReadyCallback := func(details openfeature.EventDetails) {
		readyHandlerRunCount++
	}
	providerErrorCallback := func(details openfeature.EventDetails) {
		errorHandlerRunCount++
	}
	client.AddHandler(openfeature.ProviderError, &providerErrorCallback)
	client.AddHandler(openfeature.ProviderReady, &providerReadyCallback)
	return ctx, nil
}

func aProviderReadyHandlerMustRunWhenConnectionEstablished(ctx context.Context) (context.Context, error) {
	if errorHandlerRunCount != 0 {
		return ctx, fmt.Errorf("expected errorHandlerRunCount to be 0 but found %d", errorHandlerRunCount)
	}
	readyHandlerAssert := func() bool {
		return readyHandlerRunCount == 1
	}
	if !waitUntilTimeOut(readyHandlerAssert, time.Second*15) {
		return ctx, fmt.Errorf("expected readyHandlerRunCount to be 1 but found %d", readyHandlerRunCount)
	}
	return ctx, nil
}

func aProviderErrorHandlerMustRunWhenConnectionLost(ctx context.Context) (context.Context, error) {
	errorHandlerAssert := func() bool {
		return errorHandlerRunCount > 0
	}
	if !waitUntilTimeOut(errorHandlerAssert, time.Second*15) {
		return ctx, fmt.Errorf("expected errorHandlerRunCount to be greater than 0 but found %d", errorHandlerRunCount)
	}
	return ctx, nil
}

func aProviderReadyHandlerMustRunAgainWhenConnectionReestablished(ctx context.Context) (context.Context, error) {
	readyHandlerAssert := func() bool {
		return readyHandlerRunCount > 1
	}
	if !waitUntilTimeOut(readyHandlerAssert, time.Second*15) {
		return ctx, fmt.Errorf("expected readyHandlerRunCount to be greater than 1 but found %d", readyHandlerRunCount)
	}
	return ctx, nil
}
