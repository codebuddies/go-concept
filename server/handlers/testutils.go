package handlers

import (
	"io/ioutil"
	"net/http"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"
)

func requireResponse(t *testing.T, res *http.Response, expectedStatusCode int, expectedBody string) {
	t.Helper()

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, expectedStatusCode, res.StatusCode)
	require.Equal(t, expectedBody, string(body))
}

func requireResponseRegex(t *testing.T, res *http.Response, expectedStatusCode int, expectedRegex *regexp.Regexp) {
	t.Helper()

	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		require.NoError(t, err)
	}

	require.Equal(t, expectedStatusCode, res.StatusCode)
	require.Regexp(t, expectedRegex, string(body))
}
