package testutils

import (
	"os"
	"testing"
)

type SkippableTest interface {
	Helper()
	Skip(...any)
}

var _ SkippableTest = &testing.T{}

func IntegrationTest(t SkippableTest) {
	t.Helper()
	if os.Getenv("RUN_INTEGRATION_TESTS") != "true" {
		t.Skip("RUN_INTEGRATION_TESTS environment variable is not set, skipping integration test")
		return
	}
}
