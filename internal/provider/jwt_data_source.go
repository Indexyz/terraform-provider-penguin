// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

var _ datasource.DataSource = &JWTDataSource{}

func NewJWTDataSource() datasource.DataSource {
	return &JWTDataSource{}
}

type JWTDataSource struct {
	client *penguin.Client
}

type JWTDataSourceModel struct {
	ID                   types.String `tfsdk:"id"`
	MaxTransferKB        types.Int64  `tfsdk:"max_transfer_kb"`
	AllowedInstanceTypes types.List   `tfsdk:"allowed_instance_types"`
	AllowedZones         types.List   `tfsdk:"allowed_zones"`
	MaxBandwidthMbps     types.Int64  `tfsdk:"max_bandwidth_mbps"`
	ProjectID            types.Int64  `tfsdk:"project_id"`
	TTLMinutes           types.Int64  `tfsdk:"ttl_minutes"`

	Token     types.String `tfsdk:"token"`
	ExpiresAt types.String `tfsdk:"expires_at"`
}

func (d *JWTDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_jwt"
}

func (d *JWTDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Issue a JWT via `POST /auth/jwt`. Requires the provider `auth_token` (legacy bearer token).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"max_transfer_kb": schema.Int64Attribute{
				Optional: true,
			},
			"allowed_instance_types": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"allowed_zones": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"max_bandwidth_mbps": schema.Int64Attribute{
				Optional: true,
			},
			"project_id": schema.Int64Attribute{
				Optional: true,
			},
			"ttl_minutes": schema.Int64Attribute{
				Required: true,
			},
			"token": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"expires_at": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *JWTDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *JWTDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var config JWTDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.MaxTransferKB.IsUnknown() ||
		config.AllowedInstanceTypes.IsUnknown() ||
		config.AllowedZones.IsUnknown() ||
		config.MaxBandwidthMbps.IsUnknown() ||
		config.ProjectID.IsUnknown() ||
		config.TTLMinutes.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown JWT data source configuration",
			"All input values must be known during planning to issue a JWT.",
		)
		return
	}

	request := penguin.IssueJWTRequest{
		TTLMinutes: config.TTLMinutes.ValueInt64(),
	}

	if !config.MaxTransferKB.IsNull() {
		v := config.MaxTransferKB.ValueInt64()
		request.MaxTransferKB = &v
	}
	if !config.MaxBandwidthMbps.IsNull() {
		v := config.MaxBandwidthMbps.ValueInt64()
		request.MaxBandwidthMbps = &v
	}
	if !config.ProjectID.IsNull() {
		v := config.ProjectID.ValueInt64()
		request.ProjectID = &v
	}

	if !config.AllowedInstanceTypes.IsNull() {
		var values []string
		resp.Diagnostics.Append(config.AllowedInstanceTypes.ElementsAs(ctx, &values, false)...)
		request.AllowedInstanceTypes = values
	}
	if !config.AllowedZones.IsNull() {
		var values []string
		resp.Diagnostics.Append(config.AllowedZones.ElementsAs(ctx, &values, false)...)
		request.AllowedZones = values
	}
	if resp.Diagnostics.HasError() {
		return
	}

	out, err := d.client.IssueJWT(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to issue JWT", err.Error())
		return
	}

	state := JWTDataSourceModel{
		ID:                   types.StringValue("jwt"),
		MaxTransferKB:        config.MaxTransferKB,
		AllowedInstanceTypes: config.AllowedInstanceTypes,
		AllowedZones:         config.AllowedZones,
		MaxBandwidthMbps:     config.MaxBandwidthMbps,
		ProjectID:            config.ProjectID,
		TTLMinutes:           config.TTLMinutes,
		Token:                types.StringValue(out.Token),
		ExpiresAt:            types.StringValue(out.ExpiresAt),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
