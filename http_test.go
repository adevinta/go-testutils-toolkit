package testutils_test

import (
	"io/ioutil"
	"net/http"
	"testing"

	testutils "github.com/adevinta/go-testutils-toolkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHTTPRresponseBuilder(t *testing.T) {
	assert.Equal(t, http.StatusBadGateway, testutils.NewHTTPResponseBuilder().WithTB(t).WithStatusCode(http.StatusBadGateway).Build().StatusCode)
	resp := testutils.NewHTTPResponseBuilder().WithTB(t).WithJsonBody(map[string]string{"hello": "world"}).Build()
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)
	require.NotNil(t, body)
	assert.JSONEq(t, `{"hello": "world"}`, string(body))
	assert.Equal(t, "application/json", resp.Header.Get("Content-Type"))
}
