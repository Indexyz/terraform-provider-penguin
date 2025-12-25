// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package penguin

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"testing"
)

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(req *http.Request) (*http.Response, error) {
	return f(req)
}

func TestBuildAuthHeader(t *testing.T) {
	t.Run("none", func(t *testing.T) {
		if got := buildAuthHeader("", ""); got != "" {
			t.Fatalf("expected empty header, got %q", got)
		}
	})

	t.Run("legacy", func(t *testing.T) {
		got := buildAuthHeader(" legacy ", "")
		if got != "Bearer legacy" {
			t.Fatalf("unexpected header: %q", got)
		}
	})

	t.Run("jwt", func(t *testing.T) {
		got := buildAuthHeader("", " jwt ")
		if got != "Bearer jwt" {
			t.Fatalf("unexpected header: %q", got)
		}
	})

	t.Run("both", func(t *testing.T) {
		got := buildAuthHeader("legacy", "jwt")
		if got != "Bearer legacy, Bearer jwt" {
			t.Fatalf("unexpected header: %q", got)
		}
	})
}

func TestClient_ListZones(t *testing.T) {
	t.Parallel()

	transport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/tencentcloud/zones" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.Header.Get("Authorization"); got != "Bearer legacy, Bearer jwt" {
			t.Fatalf("unexpected auth header: %q", got)
		}
		if got := r.Header.Get("User-Agent"); got != "test-agent" {
			t.Fatalf("unexpected user agent: %q", got)
		}

		payload, _ := json.Marshal(ZonesResponse{
			Zones: []Zone{
				{
					Region:     "ap-guangzhou",
					RegionName: "Guangzhou",
					Zone:       "ap-guangzhou-1",
					ZoneName:   "Guangzhou 1",
					State:      "AVAILABLE",
				},
			},
		})

		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(payload)),
		}, nil
	})

	client, err := NewClient("http://example.com", "legacy", "jwt", ClientOptions{
		HTTPClient: &http.Client{Transport: transport},
		UserAgent:  "test-agent",
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	zones, err := client.ListZones(context.Background())
	if err != nil {
		t.Fatalf("ListZones error: %v", err)
	}
	if len(zones) != 1 || zones[0].Zone != "ap-guangzhou-1" {
		t.Fatalf("unexpected zones: %#v", zones)
	}
}

func TestClient_ErrorResponse(t *testing.T) {
	t.Parallel()

	transport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		return &http.Response{
			StatusCode: http.StatusUnauthorized,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewBufferString(`{"status":401,"message":"unauthorized"}`)),
		}, nil
	})

	client, err := NewClient("http://example.com", "", "", ClientOptions{
		HTTPClient: &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	_, err = client.ListZones(context.Background())
	if err == nil {
		t.Fatalf("expected error")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected APIError, got %T (%v)", err, err)
	}
	if apiErr.Status != 401 || apiErr.Message != "unauthorized" {
		t.Fatalf("unexpected api error: %#v", apiErr)
	}
}

func TestClient_SelectBandwidthPackage(t *testing.T) {
	t.Parallel()

	transport := roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		if r.URL.Path != "/tencentcloud/bandwidth-packages" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		if got := r.URL.Query().Get("region"); got != "ap-guangzhou" {
			t.Fatalf("unexpected region: %q", got)
		}
		if got := r.URL.Query().Get("networkType"); got != "BGP" {
			t.Fatalf("unexpected networkType: %q", got)
		}

		payload, _ := json.Marshal(BandwidthPackageSelectionResponse{
			ID:             "bwp-123",
			AvailableCount: 190,
		})
		return &http.Response{
			StatusCode: http.StatusOK,
			Header:     http.Header{"Content-Type": []string{"application/json"}},
			Body:       io.NopCloser(bytes.NewReader(payload)),
		}, nil
	})

	client, err := NewClient("http://example.com", "", "", ClientOptions{
		HTTPClient: &http.Client{Transport: transport},
	})
	if err != nil {
		t.Fatalf("NewClient error: %v", err)
	}

	out, err := client.SelectBandwidthPackage(context.Background(), "ap-guangzhou", "BGP")
	if err != nil {
		t.Fatalf("SelectBandwidthPackage error: %v", err)
	}
	if out.ID != "bwp-123" || out.AvailableCount != 190 {
		t.Fatalf("unexpected response: %#v", out)
	}
}
