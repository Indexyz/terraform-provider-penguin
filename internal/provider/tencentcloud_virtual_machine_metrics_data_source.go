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

var _ datasource.DataSource = &TencentCloudVirtualMachineMetricsDataSource{}

func NewTencentCloudVirtualMachineMetricsDataSource() datasource.DataSource {
	return &TencentCloudVirtualMachineMetricsDataSource{}
}

type TencentCloudVirtualMachineMetricsDataSource struct {
	client *penguin.Client
}

type TencentCloudVirtualMachineMetricsDataSourceModel struct {
	ID                   types.String  `tfsdk:"id"`
	Range                types.String  `tfsdk:"range"`
	CPUAveragePercent    types.Float64 `tfsdk:"cpu_average_percent"`
	MemoryAveragePercent types.Float64 `tfsdk:"memory_average_percent"`
	NetworkOutKB         types.Int64   `tfsdk:"network_out_kb"`
	NetworkInKB          types.Int64   `tfsdk:"network_in_kb"`
	Start                types.String  `tfsdk:"start"`
	End                  types.String  `tfsdk:"end"`
}

func (d *TencentCloudVirtualMachineMetricsDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_virtual_machine_metrics"
}

func (d *TencentCloudVirtualMachineMetricsDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read VM metrics from `GET /tencentcloud/vms/:id/metrics`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"range": schema.StringAttribute{
				Optional: true,
			},
			"cpu_average_percent": schema.Float64Attribute{
				Computed: true,
			},
			"memory_average_percent": schema.Float64Attribute{
				Computed: true,
			},
			"network_out_kb": schema.Int64Attribute{
				Computed: true,
			},
			"network_in_kb": schema.Int64Attribute{
				Computed: true,
			},
			"start": schema.StringAttribute{
				Computed: true,
			},
			"end": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *TencentCloudVirtualMachineMetricsDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *TencentCloudVirtualMachineMetricsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var config TencentCloudVirtualMachineMetricsDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() || config.Range.IsUnknown() {
		resp.Diagnostics.AddError("Unknown configuration", "`id` and `range` must be known during planning.")
		return
	}

	metrics, err := d.client.GetVirtualMachineMetrics(ctx, config.ID.ValueString(), config.Range.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read virtual machine metrics", err.Error())
		return
	}

	state := TencentCloudVirtualMachineMetricsDataSourceModel{
		ID:                   config.ID,
		Range:                types.StringValue(metrics.Range),
		CPUAveragePercent:    types.Float64Value(metrics.CPUAveragePercent),
		MemoryAveragePercent: types.Float64Value(metrics.MemoryAveragePercent),
		NetworkOutKB:         types.Int64Value(metrics.NetworkOutKB),
		NetworkInKB:          types.Int64Value(metrics.NetworkInKB),
		Start:                types.StringValue(metrics.Start),
		End:                  types.StringValue(metrics.End),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
