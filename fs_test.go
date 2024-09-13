package testutils_test

import (
	"testing"

	testutils "github.com/adevinta/go-testutils-toolkit"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAssertFileContents(t *testing.T) {
	t.Run("When the file does not exist on the filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fakeT := &testutils.FakeTest{}
		assert.False(t, testutils.AssertFileContents(fakeT, fs, "/hello/world", "my-content"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "open /hello/world: file does not exist")
	})
	t.Run("When the file content differs", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, fs, "/hello/world", "this is wrong")
		fakeT := &testutils.FakeTest{}
		assert.False(t, testutils.AssertFileContents(fakeT, fs, "/hello/world", "hello world"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], `expected: "hello world"`)
		assert.Contains(t, fakeT.ErrorMessages[0], `actual  : "this is wrong"`)
	})
	t.Run("When the file contents matches on the file system", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, fs, "/hello/world", "hello world")
		fakeT := &testutils.FakeTest{}
		assert.True(t, testutils.AssertFileContents(fakeT, fs, "/hello/world", "hello world"))
		assert.False(t, fakeT.Failed)
		assert.Empty(t, fakeT.ErrorMessages)
	})
}

func TestAssertFileExists(t *testing.T) {
	t.Run("When the file does not exist on the filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fakeT := &testutils.FakeTest{}
		assert.False(t, testutils.AssertFileExists(fakeT, fs, "/hello/world"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expect file path /hello/world to exist in filesystem MemMapFS")
	})
	t.Run("When the file exists on the filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, fs, "/hello/world", "hello world")
		fakeT := &testutils.FakeTest{}
		assert.True(t, testutils.AssertFileExists(fakeT, fs, "/hello/world"))
		assert.False(t, fakeT.Failed)
		assert.Empty(t, fakeT.ErrorMessages)
	})
}

func TestRequireFileExists(t *testing.T) {
	t.Run("When the file does not exist on the filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		fakeT := &testutils.FakeTest{}
		testutils.RequireFileExists(fakeT, fs, "/hello/world")
		assert.True(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expect file path /hello/world to exist in filesystem MemMapFS")
	})
	t.Run("When the file exists on the filesystem", func(t *testing.T) {
		fs := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, fs, "/hello/world", "hello world")
		fakeT := &testutils.FakeTest{}
		testutils.RequireFileExists(fakeT, fs, "/hello/world")
		assert.False(t, fakeT.Failed)
		assert.Empty(t, fakeT.ErrorMessages)
	})
}

func TestAssertFileEquivalent(t *testing.T) {
	t.Run("When the file does not exist neither on the reference file system nor on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()

		fakeT := &testutils.FakeTest{}

		assert.True(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 0)
	})
	t.Run("When the file exists on the reference file system but not on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		assert.False(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expecting no error when opening file path /my/path but got:\n")
		assert.Contains(t, fakeT.ErrorMessages[0], "open /my/path: file does not exist")
	})
	t.Run("When the file does not exist on the reference file system but does on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		assert.False(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expecting an error open /my/path: file does not exist when opening file path /my/path but got nil")
	})
	t.Run("When the file stat differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello")

		fakeT := &testutils.FakeTest{}

		assert.False(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
	})
	t.Run("When the file has the same stats (size) but contents differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world-1")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world-2")

		fakeT := &testutils.FakeTest{}

		assert.False(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
	})
	t.Run("When the file has the same stats (size) and contents differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		assert.True(t, testutils.AssertFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path"))
		assert.False(t, fakeT.Failed)
	})
}

func TestRequireFileEquivalent(t *testing.T) {
	t.Run("When the file does not exist neither on the reference file system nor on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.False(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 0)
	})
	t.Run("When the file exists on the reference file system but not on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.True(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expecting no error when opening file path /my/path but got:\n")
		assert.Contains(t, fakeT.ErrorMessages[0], "open /my/path: file does not exist")
	})
	t.Run("When the file does not exist on the reference file system but does on the actual one", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.True(t, fakeT.Failed)
		require.Len(t, fakeT.ErrorMessages, 1)
		assert.Contains(t, fakeT.ErrorMessages[0], "Expecting an error open /my/path: file does not exist when opening file path /my/path but got nil")
	})
	t.Run("When the file stat differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello")

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.True(t, fakeT.Failed)
	})
	t.Run("When the file has the same stats (size) but contents differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world-1")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world-2")

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.True(t, fakeT.Failed)
	})
	t.Run("When the file has the same stats (size) and contents differs on the expected and the actual File Systems", func(t *testing.T) {
		expectedFS := afero.NewMemMapFs()
		actualFS := afero.NewMemMapFs()
		testutils.EnsureFileContent(t, expectedFS, "/my/path", "hello-world")
		testutils.EnsureFileContent(t, actualFS, "/my/path", "hello-world")

		fakeT := &testutils.FakeTest{}

		testutils.RequireFsFileEquivalent(fakeT, expectedFS, actualFS, "/my/path")
		assert.False(t, fakeT.Failed)
	})
}
