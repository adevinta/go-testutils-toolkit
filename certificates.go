package testutils

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func publicKey(priv interface{}) interface{} {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	default:
		return nil
	}
}

func pemBlockForKey(priv interface{}) *pem.Block {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(k)}
	case *ecdsa.PrivateKey:
		b, err := x509.MarshalECPrivateKey(k)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Unable to marshal ECDSA private key: %v", err)
			os.Exit(2)
		}
		return &pem.Block{Type: "EC PRIVATE KEY", Bytes: b}
	default:
		return nil
	}
}

// NewSelfSignedCertificate generates a new self-signed public key and certificate in the destination folder.
//
// The certificate is stored in destinationFolder/tls.crt
// The private key in destinationFolder/tls.key
func NewSelfSignedCertificate(t require.TestingT, fs afero.Fs, destinationFolder string, hosts ...string) {
	if h, ok := t.(TestHelper); ok {
		h.Helper()
	}
	require.NoError(t, fs.MkdirAll(destinationFolder, 0755))
	priv, err := ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	require.NoError(t, err)
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName:   hosts[0],
			Organization: []string{"adevinta-toolkit-integration-tests"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(time.Hour * 24 * 180),
		DNSNames:              hosts,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	require.NoError(t, err)
	certFD, err := fs.Create(filepath.Join(destinationFolder, "tls.crt"))
	require.NoError(t, err)
	defer certFD.Close()
	require.NoError(t, pem.Encode(certFD, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}))

	keyFD, err := fs.Create(filepath.Join(destinationFolder, "tls.key"))
	require.NoError(t, err)
	defer keyFD.Close()
	require.NoError(t, pem.Encode(keyFD, pemBlockForKey(priv)))
}
