# testutils

A helper to write more concise and self-explanatory tests.

## Examples

````
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
```
