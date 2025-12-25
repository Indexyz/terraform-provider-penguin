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

var _ datasource.DataSource = &TencentCloudVirtualMachineVNCDataSource{}

func NewTencentCloudVirtualMachineVNCDataSource() datasource.DataSource {
	return &TencentCloudVirtualMachineVNCDataSource{}
}

type TencentCloudVirtualMachineVNCDataSource struct {
	client *penguin.Client
}

type TencentCloudVirtualMachineVNCDataSourceModel struct {
	ID  types.String `tfsdk:"id"`
	URL types.String `tfsdk:"url"`
}

func (d *TencentCloudVirtualMachineVNCDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_virtual_machine_vnc"
}

func (d *TencentCloudVirtualMachineVNCDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Get the Tencent Cloud VNC websocket URL from `GET /tencentcloud/vms/:id/vnc`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"url": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *TencentCloudVirtualMachineVNCDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *TencentCloudVirtualMachineVNCDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var config TencentCloudVirtualMachineVNCDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() {
		resp.Diagnostics.AddError("Unknown virtual machine id", "`id` must be known during planning.")
		return
	}

	out, err := d.client.GetVirtualMachineVNC(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read virtual machine VNC URL", err.Error())
		return
	}

	state := TencentCloudVirtualMachineVNCDataSourceModel{
		ID:  config.ID,
		URL: types.StringValue(out.URL),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
