package testutils_test

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	testutils "github.com/adevinta/go-testutils-toolkit"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func ExampleNewSelfSignedCertificate() {
	// use the real *testing.T from the test
	t := &testutils.FakeTest{Name: "TestICanCreateCertificates"}
	fs := afero.NewMemMapFs()

	testutils.NewSelfSignedCertificate(t, fs, "/my/certificates", "localhost")

	fd, err := fs.Open("/my/certificates")
	if err != nil {
		return
	}
	names, err := fd.Readdirnames(-1)
	if err != nil {
		return
	}

	sort.Strings(names)

	fmt.Println(strings.Join(names, " "))

	// Usually, this is done by the go framework
	fmt.Println(t)
	// Output:
	// tls.crt tls.key
	// --- PASS: TestICanCreateCertificates
}

// start ReadMe examples

func TestICanCreateCertificates(t *testing.T) {
	fs := afero.NewMemMapFs()

	testutils.NewSelfSignedCertificate(t, fs, "/my/certificates", "localhost")

	fd, err := fs.Open("/my/certificates")
	if err != nil {
		return
	}
	names, err := fd.Readdirnames(-1)
	if err != nil {
		return
	}

	assert.Contains(t, names, "tls.crt")
	assert.Contains(t, names, "tls.key")
}
