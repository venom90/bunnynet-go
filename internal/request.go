// Package internal provides internal utilities for the Bunny.net API client
package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/google/go-querystring/query"
	"github.com/venom90/bunnynet-go-client/common"
)

// DoRequest sends an HTTP request and returns the response
func DoRequest(client *http.Client, req *http.Request) (*http.Response, error) {
	resp, err := client.Do(req)
	if err != nil {
		return nil, common.NewClientError("failed to send request", err)
	}

	if resp.StatusCode >= 400 {
		err := common.ParseErrorResponse(resp)
		return resp, err
	}

	return resp, nil
}

// NewRequest creates a new HTTP request with the given method, URL, and body
func NewRequest(method, baseURL, path string, body interface{}, apiKey, userAgent string) (*http.Request, error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return nil, common.NewClientError("failed to parse base URL", err)
	}

	// Resolve the path relative to the base URL
	u = u.ResolveReference(&url.URL{Path: path})

	var buf io.ReadWriter
	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return nil, common.NewClientError("failed to encode request body", err)
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, common.NewClientError("failed to create request", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")

	if body != nil && (method == http.MethodPost || method == http.MethodPut || method == http.MethodPatch) {
		req.Header.Set("Content-Type", "application/json")
	}

	if apiKey != "" {
		req.Header.Set("AccessKey", apiKey)
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	return req, nil
}

// AddQueryParams adds query parameters to the request URL
func AddQueryParams(req *http.Request, params interface{}) error {
	if params == nil {
		return nil
	}

	// For types that implement RequestParams
	if rp, ok := params.(common.RequestParams); ok {
		q := req.URL.Query()
		for k, v := range rp.ToQueryParams() {
			q.Add(k, v)
		}
		req.URL.RawQuery = q.Encode()
		return nil
	}

	// For struct types with tags
	v, err := query.Values(params)
	if err != nil {
		return common.NewClientError("failed to encode query parameters", err)
	}

	// If the URL already has a query string, merge the values
	q := req.URL.Query()
	for k, vs := range v {
		for _, v := range vs {
			q.Add(k, v)
		}
	}
	req.URL.RawQuery = q.Encode()

	return nil
}
