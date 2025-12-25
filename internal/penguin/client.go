// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package penguin

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strings"
	"time"
)

const (
	defaultTimeout   = 60 * time.Second
	maxResponseBytes = 2 << 20 // 2 MiB
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
	authHeader string
	userAgent  string
}

type ClientOptions struct {
	HTTPClient *http.Client
	UserAgent  string
}

func NewClient(endpoint string, legacyToken string, jwt string, opts ClientOptions) (*Client, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, errors.New("endpoint is required")
	}

	parsed, err := url.Parse(endpoint)
	if err != nil {
		return nil, fmt.Errorf("parse endpoint: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return nil, fmt.Errorf("endpoint must include scheme and host, got %q", endpoint)
	}
	parsed.Path = strings.TrimRight(parsed.Path, "/")

	httpClient := opts.HTTPClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: defaultTimeout}
	}

	userAgent := strings.TrimSpace(opts.UserAgent)
	if userAgent == "" {
		userAgent = "terraform-provider-penguin"
	}

	return &Client{
		baseURL:    parsed,
		httpClient: httpClient,
		authHeader: buildAuthHeader(legacyToken, jwt),
		userAgent:  userAgent,
	}, nil
}

func buildAuthHeader(legacyToken string, jwt string) string {
	legacyToken = strings.TrimSpace(legacyToken)
	jwt = strings.TrimSpace(jwt)
	if legacyToken == "" && jwt == "" {
		return ""
	}
	if legacyToken != "" && jwt != "" {
		return fmt.Sprintf("Bearer %s, Bearer %s", legacyToken, jwt)
	}
	if legacyToken != "" {
		return "Bearer " + legacyToken
	}
	return "Bearer " + jwt
}

func (c *Client) urlFor(p string) string {
	clone := *c.baseURL
	clone.Path = path.Join(c.baseURL.Path, p)
	return clone.String()
}

func (c *Client) doJSON(ctx context.Context, method string, p string, query url.Values, in any, out any, okStatuses ...int) error {
	var body io.Reader
	if in != nil {
		payload, err := json.Marshal(in)
		if err != nil {
			return fmt.Errorf("marshal request: %w", err)
		}
		body = bytes.NewReader(payload)
	}

	req, err := http.NewRequestWithContext(ctx, method, c.urlFor(p), body)
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	if len(query) > 0 {
		req.URL.RawQuery = query.Encode()
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", c.userAgent)
	if in != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	if c.authHeader != "" {
		req.Header.Set("Authorization", c.authHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	reader := io.LimitReader(resp.Body, maxResponseBytes)
	respBytes, readErr := io.ReadAll(reader)
	if readErr != nil {
		return fmt.Errorf("read response: %w", readErr)
	}

	if len(okStatuses) == 0 {
		okStatuses = []int{http.StatusOK}
	}
	if !statusIn(resp.StatusCode, okStatuses) {
		return parseAPIError(resp.StatusCode, respBytes)
	}

	if out == nil {
		return nil
	}
	if len(respBytes) == 0 {
		return nil
	}

	if err := json.Unmarshal(respBytes, out); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func statusIn(got int, allowed []int) bool {
	for _, s := range allowed {
		if got == s {
			return true
		}
	}
	return false
}

func parseAPIError(status int, respBytes []byte) error {
	var apiErr APIError
	if len(respBytes) > 0 && json.Unmarshal(respBytes, &apiErr) == nil && apiErr.Message != "" {
		if apiErr.Status == 0 {
			apiErr.Status = status
		}
		return &apiErr
	}

	msg := strings.TrimSpace(string(respBytes))
	if msg == "" {
		msg = http.StatusText(status)
	}
	return &APIError{
		Status:  status,
		Message: msg,
	}
}
