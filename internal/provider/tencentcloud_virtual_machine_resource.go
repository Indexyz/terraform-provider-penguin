// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

var (
	_ resource.Resource = &TencentCloudVirtualMachineResource{}
)

func NewTencentCloudVirtualMachineResource() resource.Resource {
	return &TencentCloudVirtualMachineResource{}
}

type TencentCloudVirtualMachineResource struct {
	client *penguin.Client
}

type TencentCloudVirtualMachineResourceModel struct {
	ID types.String `tfsdk:"id"`

	Name           types.String `tfsdk:"name"`
	Zone           types.String `tfsdk:"zone"`
	InstanceType   types.String `tfsdk:"instance_type"`
	SecurityGroup  types.String `tfsdk:"security_group"`
	SystemImage    types.String `tfsdk:"system_image"`
	VPCID          types.String `tfsdk:"vpc_id"`
	SubnetID       types.String `tfsdk:"subnet_id"`
	PrivateIP      types.String `tfsdk:"private_ip_address"`
	SystemDiskGiB  types.Int64  `tfsdk:"system_disk_size_gib"`
	SharedBWPKGID  types.String `tfsdk:"shared_bandwidth_package_id"`
	ElasticIPID    types.String `tfsdk:"elastic_ip_id"`
	BandwidthLimit types.Int64  `tfsdk:"bandwidth_limit_mbps"`
	ChargeType     types.String `tfsdk:"charge_type"`
	RootPassword   types.String `tfsdk:"root_login_password"`
	TotalTransfer  types.Int64  `tfsdk:"total_transfer_kb"`
	ProjectID      types.Int64  `tfsdk:"project_id"`
	PeriodMonths   types.Int64  `tfsdk:"period_months"`
	CloudInitData  types.String `tfsdk:"cloud_init_data"`
	AutoRenew      types.Bool   `tfsdk:"auto_renew"`

	InstanceID        types.String `tfsdk:"instance_id"`
	InstanceState     types.String `tfsdk:"instance_state"`
	CPU               types.Int64  `tfsdk:"cpu"`
	MemoryGiB         types.Int64  `tfsdk:"memory_gib"`
	PrivateIPs        types.List   `tfsdk:"private_ips"`
	PublicIPs         types.List   `tfsdk:"public_ips"`
	ImageID           types.String `tfsdk:"image_id"`
	OSName            types.String `tfsdk:"os_name"`
	CreatedAt         types.String `tfsdk:"created_at"`
	ExpiredAt         types.String `tfsdk:"expired_at"`
	UsedTransferKB    types.Int64  `tfsdk:"used_transfer_kb"`
	RemainingTransfer types.Int64  `tfsdk:"remaining_transfer_kb"`
	Password          types.String `tfsdk:"password"`
	DefaultLoginUser  types.String `tfsdk:"default_login_user"`
}

func (r *TencentCloudVirtualMachineResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_virtual_machine"
}

func (r *TencentCloudVirtualMachineResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	replaceStrings := []planmodifier.String{
		stringplanmodifier.RequiresReplace(),
	}
	replaceInt64 := []planmodifier.Int64{
		int64planmodifier.RequiresReplace(),
	}

	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Tencent Cloud CVM instances via the Penguin service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},

			"name": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"zone": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"instance_type": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"security_group": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"system_image": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"vpc_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"subnet_id": schema.StringAttribute{
				Required:      true,
				PlanModifiers: replaceStrings,
			},
			"private_ip_address": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: replaceStrings,
			},
			"system_disk_size_gib": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: replaceInt64,
			},
			"shared_bandwidth_package_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: replaceStrings,
			},
			"elastic_ip_id": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: replaceStrings,
			},
			"bandwidth_limit_mbps": schema.Int64Attribute{
				Optional: true,
			},
			"charge_type": schema.StringAttribute{
				Optional:      true,
				PlanModifiers: replaceStrings,
			},
			"root_login_password": schema.StringAttribute{
				Optional:      true,
				Sensitive:     true,
				PlanModifiers: replaceStrings,
			},
			"total_transfer_kb": schema.Int64Attribute{
				Required:      true,
				PlanModifiers: replaceInt64,
			},
			"project_id": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: replaceInt64,
			},
			"period_months": schema.Int64Attribute{
				Optional:      true,
				PlanModifiers: replaceInt64,
			},
			"cloud_init_data": schema.StringAttribute{
				Optional:      true,
				Sensitive:     true,
				PlanModifiers: replaceStrings,
			},
			"auto_renew": schema.BoolAttribute{
				Optional: true,
			},

			"instance_id": schema.StringAttribute{
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
			"used_transfer_kb": schema.Int64Attribute{
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

func (r *TencentCloudVirtualMachineResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	configureResourceClient(req, resp, &r.client)
}

func (r *TencentCloudVirtualMachineResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var plan TencentCloudVirtualMachineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.IsUnknown() ||
		plan.Zone.IsUnknown() ||
		plan.InstanceType.IsUnknown() ||
		plan.SecurityGroup.IsUnknown() ||
		plan.SystemImage.IsUnknown() ||
		plan.VPCID.IsUnknown() ||
		plan.SubnetID.IsUnknown() ||
		plan.PrivateIP.IsUnknown() ||
		plan.SystemDiskGiB.IsUnknown() ||
		plan.SharedBWPKGID.IsUnknown() ||
		plan.ElasticIPID.IsUnknown() ||
		plan.BandwidthLimit.IsUnknown() ||
		plan.ChargeType.IsUnknown() ||
		plan.RootPassword.IsUnknown() ||
		plan.TotalTransfer.IsUnknown() ||
		plan.ProjectID.IsUnknown() ||
		plan.PeriodMonths.IsUnknown() ||
		plan.CloudInitData.IsUnknown() ||
		plan.AutoRenew.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown virtual machine configuration",
			"All input values must be known during planning to create a virtual machine.",
		)
		return
	}

	if !plan.ElasticIPID.IsNull() && (!plan.SharedBWPKGID.IsNull() || !plan.BandwidthLimit.IsNull()) {
		resp.Diagnostics.AddError(
			"Invalid network configuration",
			"When `elastic_ip_id` is set, omit `shared_bandwidth_package_id` and `bandwidth_limit_mbps`.",
		)
		return
	}
	if plan.ElasticIPID.IsNull() && plan.BandwidthLimit.IsNull() {
		resp.Diagnostics.AddError(
			"Missing bandwidth limit",
			"Set `bandwidth_limit_mbps` when not providing `elastic_ip_id`.",
		)
		return
	}

	request := penguin.CreateVirtualMachineRequest{
		Name:              plan.Name.ValueString(),
		Zone:              plan.Zone.ValueString(),
		InstanceType:      plan.InstanceType.ValueString(),
		SecurityGroup:     plan.SecurityGroup.ValueString(),
		SystemImage:       plan.SystemImage.ValueString(),
		VPCID:             plan.VPCID.ValueString(),
		SubnetID:          plan.SubnetID.ValueString(),
		SystemDiskSizeGiB: plan.SystemDiskGiB.ValueInt64(),
		TotalTransferKB:   plan.TotalTransfer.ValueInt64(),
	}

	if !plan.PrivateIP.IsNull() {
		v := plan.PrivateIP.ValueString()
		request.PrivateIPAddress = &v
	}
	if !plan.SharedBWPKGID.IsNull() {
		v := plan.SharedBWPKGID.ValueString()
		request.SharedBandwidthPackageID = &v
	}
	if !plan.ElasticIPID.IsNull() {
		v := plan.ElasticIPID.ValueString()
		request.ElasticIPID = &v
	}
	if !plan.BandwidthLimit.IsNull() {
		v := plan.BandwidthLimit.ValueInt64()
		request.BandwidthLimitMbps = &v
	}
	if !plan.ChargeType.IsNull() {
		v := plan.ChargeType.ValueString()
		request.ChargeType = &v
	}
	if !plan.RootPassword.IsNull() {
		v := plan.RootPassword.ValueString()
		request.RootLoginPassword = &v
	}
	if !plan.ProjectID.IsNull() {
		v := plan.ProjectID.ValueInt64()
		request.ProjectID = &v
	}
	if !plan.PeriodMonths.IsNull() {
		v := plan.PeriodMonths.ValueInt64()
		request.PeriodMonths = &v
	}
	if !plan.CloudInitData.IsNull() {
		v := plan.CloudInitData.ValueString()
		request.CloudInitData = &v
	}
	if !plan.AutoRenew.IsNull() {
		v := plan.AutoRenew.ValueBool()
		request.AutoRenew = &v
	}

	out, err := r.client.CreateVirtualMachine(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create virtual machine", err.Error())
		return
	}

	status, err := r.waitForVirtualMachineStatus(ctx, out.ID, 10*time.Second)
	if err != nil {
		resp.Diagnostics.AddError("Failed to read virtual machine status after creation", err.Error())
		return
	}

	state := plan
	state.ID = types.StringValue(out.ID)
	resp.Diagnostics.Append(r.applyStatusToState(ctx, status, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudVirtualMachineResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var state TencentCloudVirtualMachineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	status, err := r.client.GetVirtualMachineStatus(ctx, state.ID.ValueString())
	if err != nil {
		if isNotFound(err) {
			resp.State.RemoveResource(ctx)
			return
		}
		resp.Diagnostics.AddError("Failed to read virtual machine status", err.Error())
		return
	}

	resp.Diagnostics.Append(r.applyStatusToState(ctx, status, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudVirtualMachineResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var plan TencentCloudVirtualMachineResourceModel
	var state TencentCloudVirtualMachineResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.BandwidthLimit.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown bandwidth limit",
			"`bandwidth_limit_mbps` must be known during planning to update bandwidth.",
		)
		return
	}

	if !plan.BandwidthLimit.IsNull() && (state.BandwidthLimit.IsNull() || plan.BandwidthLimit.ValueInt64() != state.BandwidthLimit.ValueInt64()) {
		if err := r.client.AdjustVirtualMachineBandwidth(ctx, state.ID.ValueString(), plan.BandwidthLimit.ValueInt64()); err != nil {
			resp.Diagnostics.AddError("Failed to adjust virtual machine bandwidth", err.Error())
			return
		}
	}

	status, err := r.client.GetVirtualMachineStatus(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Failed to read virtual machine status", err.Error())
		return
	}

	newState := plan
	newState.ID = state.ID
	resp.Diagnostics.Append(r.applyStatusToState(ctx, status, &newState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &newState)...)
}

func (r *TencentCloudVirtualMachineResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var state TencentCloudVirtualMachineResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if err := r.client.DeleteVirtualMachine(ctx, state.ID.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Failed to delete virtual machine", err.Error())
		return
	}

	err := waitUntil(ctx, 10*time.Second, func(ctx context.Context) (bool, error) {
		_, err := r.client.GetVirtualMachineStatus(ctx, state.ID.ValueString())
		if err == nil {
			return false, nil
		}
		if isNotFound(err) {
			return true, nil
		}
		return false, err
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed waiting for virtual machine deletion", err.Error())
		return
	}
}

func (r *TencentCloudVirtualMachineResource) waitForVirtualMachineStatus(ctx context.Context, id string, interval time.Duration) (*penguin.VirtualMachineStatus, error) {
	var lastStatus *penguin.VirtualMachineStatus
	err := waitUntil(ctx, interval, func(ctx context.Context) (bool, error) {
		status, err := r.client.GetVirtualMachineStatus(ctx, id)
		if err != nil {
			if isNotFound(err) {
				return false, nil
			}
			return false, err
		}
		lastStatus = status
		return true, nil
	})
	if err != nil {
		return nil, err
	}
	return lastStatus, nil
}

func (r *TencentCloudVirtualMachineResource) applyStatusToState(ctx context.Context, status *penguin.VirtualMachineStatus, state *TencentCloudVirtualMachineResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	privateIPs, listDiags := types.ListValueFrom(ctx, types.StringType, status.PrivateIPs)
	diags.Append(listDiags...)
	publicIPs, listDiags := types.ListValueFrom(ctx, types.StringType, status.PublicIPs)
	diags.Append(listDiags...)

	state.InstanceID = types.StringValue(status.InstanceID)
	state.Zone = types.StringValue(status.Zone)
	state.InstanceType = types.StringValue(status.InstanceType)
	state.InstanceState = types.StringValue(status.InstanceState)
	state.CPU = types.Int64Value(status.CPU)
	state.MemoryGiB = types.Int64Value(status.MemoryGiB)
	state.PrivateIPs = privateIPs
	state.PublicIPs = publicIPs
	state.ImageID = types.StringPointerValue(status.ImageID)
	state.OSName = types.StringPointerValue(status.OSName)
	state.CreatedAt = types.StringPointerValue(status.CreatedAt)
	state.ExpiredAt = types.StringPointerValue(status.ExpiredAt)
	state.UsedTransferKB = types.Int64Value(status.UsedTransfer)
	state.RemainingTransfer = types.Int64PointerValue(status.RemainingTransfer)
	state.Password = types.StringPointerValue(status.Password)
	state.DefaultLoginUser = types.StringPointerValue(status.DefaultLoginUser)

	return diags
}
