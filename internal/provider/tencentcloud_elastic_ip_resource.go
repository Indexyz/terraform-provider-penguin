// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

var (
	_ resource.Resource                = &TencentCloudElasticIPResource{}
	_ resource.ResourceWithImportState = &TencentCloudElasticIPResource{}
)

func NewTencentCloudElasticIPResource() resource.Resource {
	return &TencentCloudElasticIPResource{}
}

type TencentCloudElasticIPResource struct {
	client *penguin.Client
}

type TencentCloudElasticIPResourceModel struct {
	ID                       types.String `tfsdk:"id"`
	Region                   types.String `tfsdk:"region"`
	BandwidthLimitMbps       types.Int64  `tfsdk:"bandwidth_limit_mbps"`
	AddressName              types.String `tfsdk:"address_name"`
	SharedBandwidthPackageID types.String `tfsdk:"shared_bandwidth_package_id"`
	Address                  types.String `tfsdk:"address"`
}

func (r *TencentCloudElasticIPResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_elastic_ip"
}

func (r *TencentCloudElasticIPResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manage Tencent Cloud elastic IPs via the Penguin service. Note: the Penguin API currently has no read endpoint for EIPs, so refresh cannot detect out-of-band changes.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"region": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"bandwidth_limit_mbps": schema.Int64Attribute{
				Required: true,
				PlanModifiers: []planmodifier.Int64{
					int64planmodifier.RequiresReplace(),
				},
			},
			"address_name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"shared_bandwidth_package_id": schema.StringAttribute{
				Optional: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"address": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (r *TencentCloudElasticIPResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	configureResourceClient(req, resp, &r.client)
}

func (r *TencentCloudElasticIPResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var plan TencentCloudElasticIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Region.IsUnknown() || plan.BandwidthLimitMbps.IsUnknown() || plan.AddressName.IsUnknown() || plan.SharedBandwidthPackageID.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown elastic IP configuration",
			"All input values must be known during planning to create an elastic IP.",
		)
		return
	}

	request := penguin.CreateElasticIPRequest{
		Region:             plan.Region.ValueString(),
		BandwidthLimitMbps: plan.BandwidthLimitMbps.ValueInt64(),
		AddressName:        plan.AddressName.ValueString(),
	}
	if !plan.SharedBandwidthPackageID.IsNull() {
		v := plan.SharedBandwidthPackageID.ValueString()
		request.SharedBandwidthPackageID = &v
	}

	out, err := r.client.CreateElasticIP(ctx, request)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create elastic IP", err.Error())
		return
	}

	state := plan
	state.ID = types.StringValue(out.ID)
	state.Address = types.StringValue(out.Address)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudElasticIPResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TencentCloudElasticIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudElasticIPResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TencentCloudElasticIPResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TencentCloudElasticIPResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var state TencentCloudElasticIPResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.Region.IsUnknown() || state.ID.IsUnknown() {
		resp.Diagnostics.AddError("Unknown state", "Cannot delete elastic IP with unknown state.")
		return
	}

	if err := r.client.DeleteElasticIP(ctx, state.Region.ValueString(), state.ID.ValueString()); err != nil {
		if isNotFound(err) {
			return
		}
		resp.Diagnostics.AddError("Failed to delete elastic IP", err.Error())
		return
	}
}

func (r *TencentCloudElasticIPResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	parts := strings.SplitN(req.ID, ":", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Import Identifier",
			"Expected import identifier with format `region:id` (e.g. `ap-guangzhou:eip-12345678`).",
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("region"), parts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), parts[1])...)
}
