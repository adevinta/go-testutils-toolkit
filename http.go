package testutils

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type RoundTripperFunc func(r *http.Request) (*http.Response, error)

func (r RoundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	if r != nil {
		return r(req)
	}
	return nil, errors.New("nil round-tripper func")
}

// RoundTripperFunc should implement the http.RoundTripper interface
var _ http.RoundTripper = RoundTripperFunc(nil)

type HTTPResponseBuilder struct {
	resp *http.Response
	tb   testing.TB
}

func StringBody(body string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(body))
}

func NewHTTPResponseBuilder() *HTTPResponseBuilder {
	return &HTTPResponseBuilder{resp: &http.Response{}}
}

func (b *HTTPResponseBuilder) Build() *http.Response {
	return b.resp
}

func (b *HTTPResponseBuilder) WithTB(tb testing.TB) *HTTPResponseBuilder {
	b.tb = tb
	return b
}

func (b *HTTPResponseBuilder) WithStatusCode(code int) *HTTPResponseBuilder {
	b.resp.StatusCode = code
	return b
}
func (b *HTTPResponseBuilder) WithBody(body io.ReadCloser) *HTTPResponseBuilder {
	b.resp.Body = body
	return b
}

func (b *HTTPResponseBuilder) WithJsonBody(body interface{}) *HTTPResponseBuilder {
	data := bytes.Buffer{}
	err := json.NewEncoder(&data).Encode(body)
	if b.tb != nil {
		b.tb.Helper()
		assert.NoError(b.tb, err)
	}
	b.WithBody(io.NopCloser(&data))
	if b.resp.Header == nil {
		b.resp.Header = http.Header{}
	}
	b.resp.Header.Set("Content-Type", "application/json")
	return b
}
