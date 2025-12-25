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

var _ datasource.DataSource = &InternalHealthDataSource{}

func NewInternalHealthDataSource() datasource.DataSource {
	return &InternalHealthDataSource{}
}

type InternalHealthDataSource struct {
	client *penguin.Client
}

type InternalHealthDataSourceModel struct {
	ID       types.String `tfsdk:"id"`
	Status   types.String `tfsdk:"status"`
	Database types.String `tfsdk:"database"`
}

func (d *InternalHealthDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_internal_health"
}

func (d *InternalHealthDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read Penguin internal health information from `/_internal/health`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"database": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *InternalHealthDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *InternalHealthDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	health, err := d.client.InternalHealth(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read internal health", err.Error())
		return
	}

	state := InternalHealthDataSourceModel{
		ID:       types.StringValue("internal_health"),
		Status:   types.StringValue(health.Status),
		Database: types.StringValue(health.Database),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
