// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package penguin

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

func (c *Client) Health(ctx context.Context) error {
	return c.doJSON(ctx, http.MethodGet, "/health", nil, nil, nil, http.StatusOK)
}

func (c *Client) InternalHealth(ctx context.Context) (*InternalHealthResponse, error) {
	var out InternalHealthResponse
	if err := c.doJSON(ctx, http.MethodGet, "/_internal/health", nil, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ListZones(ctx context.Context) ([]Zone, error) {
	var out ZonesResponse
	if err := c.doJSON(ctx, http.MethodGet, "/tencentcloud/zones", nil, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return out.Zones, nil
}

func (c *Client) SelectBandwidthPackage(ctx context.Context, region string, networkType string) (*BandwidthPackageSelectionResponse, error) {
	query := url.Values{}
	query.Set("region", region)
	if networkType != "" {
		query.Set("networkType", networkType)
	}

	var out BandwidthPackageSelectionResponse
	if err := c.doJSON(ctx, http.MethodGet, "/tencentcloud/bandwidth-packages", query, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) CreateVirtualMachine(ctx context.Context, req CreateVirtualMachineRequest) (*CreateVirtualMachineResponse, error) {
	var out CreateVirtualMachineResponse
	if err := c.doJSON(ctx, http.MethodPost, "/tencentcloud/vms", nil, req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteVirtualMachine(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/tencentcloud/vms/%s", url.PathEscape(id)), nil, nil, nil, http.StatusAccepted)
}

func (c *Client) GetVirtualMachineStatus(ctx context.Context, id string) (*VirtualMachineStatus, error) {
	var out VirtualMachineStatus
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/tencentcloud/vms/%s/status", url.PathEscape(id)), nil, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetVirtualMachineMetrics(ctx context.Context, id string, r string) (*VirtualMachineMetricsResponse, error) {
	query := url.Values{}
	if r != "" {
		query.Set("range", r)
	}
	var out VirtualMachineMetricsResponse
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/tencentcloud/vms/%s/metrics", url.PathEscape(id)), query, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) GetVirtualMachineVNC(ctx context.Context, id string) (*VirtualMachineVNCResponse, error) {
	var out VirtualMachineVNCResponse
	if err := c.doJSON(ctx, http.MethodGet, fmt.Sprintf("/tencentcloud/vms/%s/vnc", url.PathEscape(id)), nil, nil, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) AdjustVirtualMachineBandwidth(ctx context.Context, id string, bandwidthLimitMbps int64) error {
	req := AdjustBandwidthRequest{BandwidthLimitMbps: bandwidthLimitMbps}
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/tencentcloud/vms/%s/bandwidth", url.PathEscape(id)), nil, req, nil, http.StatusAccepted)
}

func (c *Client) RenewVirtualMachine(ctx context.Context, id string, req RenewVirtualMachineRequest) (*RenewVirtualMachineResponse, error) {
	var out RenewVirtualMachineResponse
	if err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/tencentcloud/vms/%s/renew", url.PathEscape(id)), nil, req, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ReinstallVirtualMachine(ctx context.Context, id string, req ReinstallVirtualMachineRequest) error {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/tencentcloud/vms/%s/reinstall", url.PathEscape(id)), nil, req, nil, http.StatusAccepted)
}

func (c *Client) ResetVirtualMachinePassword(ctx context.Context, id string, req ResetVirtualMachinePasswordRequest) (*ResetVirtualMachinePasswordResponse, error) {
	var out ResetVirtualMachinePasswordResponse
	if err := c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/tencentcloud/vms/%s/reset-password", url.PathEscape(id)), nil, req, &out, http.StatusOK); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) ResetVirtualMachineTransfer(ctx context.Context, id string) error {
	return c.doJSON(ctx, http.MethodPost, fmt.Sprintf("/tencentcloud/vms/%s/reset-transfer", url.PathEscape(id)), nil, nil, nil, http.StatusNoContent)
}

func (c *Client) CreateElasticIP(ctx context.Context, req CreateElasticIPRequest) (*CreateElasticIPResponse, error) {
	var out CreateElasticIPResponse
	if err := c.doJSON(ctx, http.MethodPost, "/tencentcloud/eips", nil, req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *Client) DeleteElasticIP(ctx context.Context, region string, id string) error {
	query := url.Values{}
	query.Set("region", region)
	return c.doJSON(ctx, http.MethodDelete, fmt.Sprintf("/tencentcloud/eips/%s", url.PathEscape(id)), query, nil, nil, http.StatusNoContent)
}

func (c *Client) IssueJWT(ctx context.Context, req IssueJWTRequest) (*IssueJWTResponse, error) {
	var out IssueJWTResponse
	if err := c.doJSON(ctx, http.MethodPost, "/auth/jwt", nil, req, &out, http.StatusCreated); err != nil {
		return nil, err
	}
	return &out, nil
}
