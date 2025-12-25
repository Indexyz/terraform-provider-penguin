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

var _ datasource.DataSource = &TencentCloudVirtualMachineStatusDataSource{}

func NewTencentCloudVirtualMachineStatusDataSource() datasource.DataSource {
	return &TencentCloudVirtualMachineStatusDataSource{}
}

type TencentCloudVirtualMachineStatusDataSource struct {
	client *penguin.Client
}

type TencentCloudVirtualMachineStatusDataSourceModel struct {
	ID                types.String `tfsdk:"id"`
	InstanceID        types.String `tfsdk:"instance_id"`
	Zone              types.String `tfsdk:"zone"`
	InstanceType      types.String `tfsdk:"instance_type"`
	InstanceState     types.String `tfsdk:"instance_state"`
	CPU               types.Int64  `tfsdk:"cpu"`
	MemoryGiB         types.Int64  `tfsdk:"memory_gib"`
	SystemDiskSizeGiB types.Int64  `tfsdk:"system_disk_size_gib"`
	PrivateIPs        types.List   `tfsdk:"private_ips"`
	PublicIPs         types.List   `tfsdk:"public_ips"`
	ImageID           types.String `tfsdk:"image_id"`
	OSName            types.String `tfsdk:"os_name"`
	CreatedAt         types.String `tfsdk:"created_at"`
	ExpiredAt         types.String `tfsdk:"expired_at"`
	TotalTransferKB   types.Int64  `tfsdk:"total_transfer_kb"`
	UsedTransferKB    types.Int64  `tfsdk:"used_transfer_kb"`
	TxTransferKB      types.Int64  `tfsdk:"tx_transfer_kb"`
	RxTransferKB      types.Int64  `tfsdk:"rx_transfer_kb"`
	RemainingTransfer types.Int64  `tfsdk:"remaining_transfer_kb"`
	Password          types.String `tfsdk:"password"`
	DefaultLoginUser  types.String `tfsdk:"default_login_user"`
}

func (d *TencentCloudVirtualMachineStatusDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_virtual_machine_status"
}

func (d *TencentCloudVirtualMachineStatusDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Read the latest VM status from `GET /tencentcloud/vms/:id/status`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Required: true,
			},
			"instance_id": schema.StringAttribute{
				Computed: true,
			},
			"zone": schema.StringAttribute{
				Computed: true,
			},
			"instance_type": schema.StringAttribute{
				Computed: true,
			},
			"instance_state": schema.StringAttribute{
				Computed: true,
			},
			"cpu": schema.Int64Attribute{
				Computed: true,
			},
			"memory_gib": schema.Int64Attribute{
				Computed: true,
			},
			"system_disk_size_gib": schema.Int64Attribute{
				Computed: true,
			},
			"private_ips": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"public_ips": schema.ListAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"image_id": schema.StringAttribute{
				Computed: true,
			},
			"os_name": schema.StringAttribute{
				Computed: true,
			},
			"created_at": schema.StringAttribute{
				Computed: true,
			},
			"expired_at": schema.StringAttribute{
				Computed: true,
			},
			"total_transfer_kb": schema.Int64Attribute{
				Computed: true,
			},
			"used_transfer_kb": schema.Int64Attribute{
				Computed: true,
			},
			"tx_transfer_kb": schema.Int64Attribute{
				Computed: true,
			},
			"rx_transfer_kb": schema.Int64Attribute{
				Computed: true,
			},
			"remaining_transfer_kb": schema.Int64Attribute{
				Computed: true,
			},
			"password": schema.StringAttribute{
				Computed:  true,
				Sensitive: true,
			},
			"default_login_user": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *TencentCloudVirtualMachineStatusDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *TencentCloudVirtualMachineStatusDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var config TencentCloudVirtualMachineStatusDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.ID.IsUnknown() {
		resp.Diagnostics.AddError("Unknown virtual machine id", "`id` must be known during planning.")
		return
	}

	status, err := d.client.GetVirtualMachineStatus(ctx, config.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read virtual machine status", err.Error())
		return
	}

	privateIPs, diags := types.ListValueFrom(ctx, types.StringType, status.PrivateIPs)
	resp.Diagnostics.Append(diags...)
	publicIPs, diags := types.ListValueFrom(ctx, types.StringType, status.PublicIPs)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	state := TencentCloudVirtualMachineStatusDataSourceModel{
		ID:                types.StringValue(status.ID),
		InstanceID:        types.StringValue(status.InstanceID),
		Zone:              types.StringValue(status.Zone),
		InstanceType:      types.StringValue(status.InstanceType),
		InstanceState:     types.StringValue(status.InstanceState),
		CPU:               types.Int64Value(status.CPU),
		MemoryGiB:         types.Int64Value(status.MemoryGiB),
		SystemDiskSizeGiB: types.Int64Value(status.SystemDiskSizeGiB),
		PrivateIPs:        privateIPs,
		PublicIPs:         publicIPs,
		TotalTransferKB:   types.Int64Value(status.TotalTransfer),
		UsedTransferKB:    types.Int64Value(status.UsedTransfer),
		ImageID:           types.StringPointerValue(status.ImageID),
		OSName:            types.StringPointerValue(status.OSName),
		CreatedAt:         types.StringPointerValue(status.CreatedAt),
		ExpiredAt:         types.StringPointerValue(status.ExpiredAt),
		Password:          types.StringPointerValue(status.Password),
		DefaultLoginUser:  types.StringPointerValue(status.DefaultLoginUser),
		TxTransferKB:      types.Int64PointerValue(status.TxTransfer),
		RxTransferKB:      types.Int64PointerValue(status.RxTransfer),
		RemainingTransfer: types.Int64PointerValue(status.RemainingTransfer),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
