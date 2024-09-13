package testutils

import (
	"os"
	"testing"

	system "github.com/adevinta/go-system-toolkit"
	"github.com/stretchr/testify/assert"
)

type fakeSkippableTest struct {
	helperCalls int
	skipCalls   int
	skipFunc    func(...any)
	helperFunc  func()
}

func (t *fakeSkippableTest) Helper() {
	t.helperCalls++
	if t.helperFunc != nil {
		t.helperFunc()
	}
}

func (t *fakeSkippableTest) Skip(args ...any) {
	t.skipCalls++
	if t.skipFunc != nil {
		t.skipFunc(args...)
	}
}

func TestIntegrationTestHelper(t *testing.T) {
	t.Run("When the RUN_INTEGRATION_TESTS environment variable is not set", func(t *testing.T) {
		t.Cleanup(system.Reset)
		os.Unsetenv("RUN_INTEGRATION_TESTS")
		fakeT := &fakeSkippableTest{
			skipFunc: func(a ...any) {
				assert.Equal(t, []any{"RUN_INTEGRATION_TESTS environment variable is not set, skipping integration test"}, a)
			},
		}
		IntegrationTest(fakeT)
		assert.Equal(t, 1, fakeT.helperCalls)
		assert.Equal(t, 1, fakeT.skipCalls)
	})
	t.Run("When the RUN_INTEGRATION_TESTS environment variable is set", func(t *testing.T) {
		t.Cleanup(system.Reset)
		t.Setenv("RUN_INTEGRATION_TESTS", "true")
		fakeT := &fakeSkippableTest{
			skipFunc: func(a ...any) {
				t.Errorf("with the RUN_INTEGRATION_TESTS environment variable, the test should not be skept")
			},
		}
		IntegrationTest(fakeT)
		assert.Equal(t, 1, fakeT.helperCalls)
		assert.Equal(t, 0, fakeT.skipCalls)
	})
}
