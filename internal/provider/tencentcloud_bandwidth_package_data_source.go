// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

var _ datasource.DataSource = &TencentCloudBandwidthPackageDataSource{}

func NewTencentCloudBandwidthPackageDataSource() datasource.DataSource {
	return &TencentCloudBandwidthPackageDataSource{}
}

type TencentCloudBandwidthPackageDataSource struct {
	client *penguin.Client
}

type TencentCloudBandwidthPackageDataSourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Region             types.String `tfsdk:"region"`
	NetworkType        types.String `tfsdk:"network_type"`
	BandwidthPackageID types.String `tfsdk:"bandwidth_package_id"`
	AvailableCount     types.Int64  `tfsdk:"available_count"`
}

func (d *TencentCloudBandwidthPackageDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_bandwidth_package"
}

func (d *TencentCloudBandwidthPackageDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Select the schedulable shared bandwidth package with the highest available capacity in a region (via `GET /tencentcloud/bandwidth-packages`).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"region": schema.StringAttribute{
				Required: true,
			},
			"network_type": schema.StringAttribute{
				Optional: true,
			},
			"bandwidth_package_id": schema.StringAttribute{
				Computed: true,
			},
			"available_count": schema.Int64Attribute{
				Computed: true,
			},
		},
	}
}

func (d *TencentCloudBandwidthPackageDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *TencentCloudBandwidthPackageDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var config TencentCloudBandwidthPackageDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Region.IsUnknown() || config.NetworkType.IsUnknown() {
		resp.Diagnostics.AddError("Unknown configuration", "`region` and `network_type` must be known during planning.")
		return
	}

	region := strings.TrimSpace(config.Region.ValueString())
	if region == "" {
		resp.Diagnostics.AddError("Invalid region", "`region` must not be empty.")
		return
	}

	out, err := d.client.SelectBandwidthPackage(ctx, region, strings.TrimSpace(config.NetworkType.ValueString()))
	if err != nil {
		resp.Diagnostics.AddError("Failed to select bandwidth package", err.Error())
		return
	}

	state := TencentCloudBandwidthPackageDataSourceModel{
		ID:                 types.StringValue(region + ":" + config.NetworkType.ValueString()),
		Region:             types.StringValue(region),
		NetworkType:        config.NetworkType,
		BandwidthPackageID: types.StringValue(out.ID),
		AvailableCount:     types.Int64Value(out.AvailableCount),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
