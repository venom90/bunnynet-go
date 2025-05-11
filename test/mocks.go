package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// MockServer creates a new test server for mocking API responses
func MockServer(t *testing.T, statusCode int, body string, validateRequest func(r *http.Request)) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if validateRequest != nil {
			validateRequest(r)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, body)
	}))
}

// AssertRequestHasHeader asserts that the request has the expected header
func AssertRequestHasHeader(t *testing.T, r *http.Request, key, value string) {
	assert.Equal(t, value, r.Header.Get(key), "Request should have the correct %s header", key)
}

// AssertRequestMethod asserts that the request has the expected method
func AssertRequestMethod(t *testing.T, r *http.Request, method string) {
	assert.Equal(t, method, r.Method, "Request should have method %s but got %s", method, r.Method)
}

// AssertRequestPath asserts that the request has the expected path
func AssertRequestPath(t *testing.T, r *http.Request, path string) {
	assert.Equal(t, path, r.URL.Path, "Request should have path %s but got %s", path, r.URL.Path)
}
