package testutils

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
)

func EnsureFileContent(t *testing.T, fs afero.Fs, path, content string) {
	t.Helper()
	fd, err := fs.Create(path)
	require.NoError(t, err)
	defer fd.Close()
	_, err = fd.WriteString(content)
	require.NoError(t, err)
}

func EnsureYAMLFileContent(t *testing.T, fs afero.Fs, path string, content interface{}) {
	t.Helper()
	fd, err := fs.Create(path)
	require.NoError(t, err)
	defer fd.Close()
	require.NoError(t, yaml.NewEncoder(fd).Encode(content))
}

func AssertFileExists(t assert.TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	_, err := fs.Stat(path)
	if os.IsNotExist(err) {
		return assert.Fail(t, fmt.Sprintf("Expect file path %s to exist in filesystem %v", path, fs.Name()), msgAndArgs...)
	}
	return true
}

func RequireFileExists(t require.TestingT, fs afero.Fs, path string, msgAndArgs ...interface{}) {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	if AssertFileExists(t, fs, path, msgAndArgs...) {
		return
	}
	t.FailNow()
}

func AssertFileContents[C string | []byte](t assert.TestingT, fs afero.Fs, path string, expected C, msgAndArgs ...interface{}) bool {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	fd, err := fs.Open(path)
	if err != nil {
		return assert.NoError(t, err, msgAndArgs...)
	}
	actual, err := io.ReadAll(fd)
	if err != nil {
		return assert.NoError(t, err, msgAndArgs...)
	}
	switch interface{}(expected).(type) {
	case string:
		if !assert.Equal(t, expected, string(actual), msgAndArgs...) {
			return false
		}
	case []byte:
		if !assert.Equal(t, expected, actual, msgAndArgs...) {
			return false
		}
	}
	return true
}

func RequireFileContents[C string | []byte](t require.TestingT, fs afero.Fs, path string, expected C, msgAndArgs ...interface{}) {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	if AssertFileContents(t, fs, path, expected, msgAndArgs...) {
		return
	}
	t.FailNow()
}

// AssertFsFileEquivalent asserts that a file is equivalent on 2 different file systems
//
//	AssertFsFileEquivalent(t, referenceFs, testedFs, path,  "error message %s", "formatted")
//
// Equivalence considers the availability of the file (no error returned when opening it),
// the file permissions, size and its content.
// Creation, Modification and Access time stamps are excluded from the comparison
func AssertFsFileEquivalent(t assert.TestingT, expected, actual afero.Fs, path string, msgAndArgs ...interface{}) bool {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	statE, errE := expected.Stat(path)
	statA, errA := actual.Stat(path)
	if errE != nil {
		if errA == nil {
			assert.Fail(t, fmt.Sprintf("Expecting an error %+v when opening file path %s but got nil", errE, path))
			return false
		} else {
			return true
		}
	} else {
		if errA != nil {
			assert.Fail(t, fmt.Sprintf("Expecting no error when opening file path %s but got:\n%+v", path, errA))
			return false
		}
	}
	if !assert.Equal(t, statE.Mode(), statA.Mode(), msgAndArgs...) {
		return false
	}
	if !assert.Equal(t, statE.Size(), statA.Size(), msgAndArgs...) {
		return false
	}
	fdE, err := expected.Open(path)
	if !assert.NoError(t, err, msgAndArgs...) {
		return false
	}
	fdA, err := actual.Open(path)
	if !assert.NoError(t, err, msgAndArgs...) {
		return false
	}
	bytesE, err := io.ReadAll(fdE)
	if !assert.NoError(t, err, msgAndArgs...) {
		return false
	}
	bytesA, err := io.ReadAll(fdA)
	if !assert.NoError(t, err, msgAndArgs...) {
		return false
	}
	if !assert.Equal(t, bytesE, bytesA, msgAndArgs...) {
		return false
	}
	return true
}

// RequireFsFileEquivalent asserts that a file is equivalent on 2 different file systems
//
//	RequireFsFileEquivalent(t, referenceFs, testedFs, path,  "error message %s", "formatted")
//
// Equivalence considers the availability of the file (no error returned when opening it),
// the file permissions, size and its content.
// Creation, Modification and Access time stamps are excluded from the comparison
func RequireFsFileEquivalent(t require.TestingT, expected, actual afero.Fs, path string, msgAndArgs ...interface{}) {

	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	if AssertFsFileEquivalent(t, expected, actual, path, msgAndArgs...) {
		return
	}
	t.FailNow()
}
