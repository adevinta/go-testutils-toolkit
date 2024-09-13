package testutils

import (
	"crypto/x509"
	"encoding/pem"
	"io"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewSelfSignedCertificate(t *testing.T) {
	fs := afero.NewMemMapFs()

	NewSelfSignedCertificate(t, fs, "/certs", "my.domain.tld")

	AssertFileExists(t, fs, "/certs/tls.crt")
	AssertFileExists(t, fs, "/certs/tls.key")

	fd, err := fs.Open("/certs/tls.crt")
	require.NoError(t, err)
	defer fd.Close()
	bytes, err := io.ReadAll(fd)
	require.NoError(t, err)
	p, rest := pem.Decode(bytes)
	assert.Empty(t, rest)
	assert.Equal(t, "CERTIFICATE", p.Type)
	cert, err := x509.ParseCertificate(p.Bytes)
	require.NoError(t, err)
	assert.Equal(t, []string{"my.domain.tld"}, cert.DNSNames)

	fd, err = fs.Open("/certs/tls.key")
	require.NoError(t, err)
	defer fd.Close()
	bytes, err = io.ReadAll(fd)
	require.NoError(t, err)
	p, rest = pem.Decode(bytes)
	assert.Empty(t, rest)
	assert.Equal(t, "EC PRIVATE KEY", p.Type)
}
