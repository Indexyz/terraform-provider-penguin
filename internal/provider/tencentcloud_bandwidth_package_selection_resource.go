// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

var _ resource.Resource = &TencentCloudBandwidthPackageSelectionResource{}

func NewTencentCloudBandwidthPackageSelectionResource() resource.Resource {
	return &TencentCloudBandwidthPackageSelectionResource{}
}

// TencentCloudBandwidthPackageSelectionResource snapshots a schedulable shared bandwidth package selection.
// The Penguin API endpoint returns the "best" current candidate, which can change over time; this resource
// stores the selected ID in Terraform state so dependents stay stable unless the resource is replaced.
type TencentCloudBandwidthPackageSelectionResource struct {
	client *penguin.Client
}

type TencentCloudBandwidthPackageSelectionResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Region             types.String `tfsdk:"region"`
	NetworkType        types.String `tfsdk:"network_type"`
	BandwidthPackageID types.String `tfsdk:"bandwidth_package_id"`
	AvailableCount     types.Int64  `tfsdk:"available_count"`
}

func (r *TencentCloudBandwidthPackageSelectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_bandwidth_package_selection"
}

func (r *TencentCloudBandwidthPackageSelectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Select (once) the schedulable shared bandwidth package in a region (via `GET /tencentcloud/bandwidth-packages`) and persist the choice in Terraform state. Use `terraform taint` to force reselection.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.UseStateForUnknown()},
			},
			"region": schema.StringAttribute{
				Required:      true,
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
			},
			"network_type": schema.StringAttribute{
				Optional:      true,
				Computed:      true,
				Default:       stringdefault.StaticString("BGP"),
				PlanModifiers: []planmodifier.String{stringplanmodifier.RequiresReplace()},
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

func (r *TencentCloudBandwidthPackageSelectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	configureResourceClient(req, resp, &r.client)
}

func (r *TencentCloudBandwidthPackageSelectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	if r.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	var plan TencentCloudBandwidthPackageSelectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Region.IsUnknown() || plan.NetworkType.IsUnknown() {
		resp.Diagnostics.AddError("Unknown configuration", "`region` and `network_type` must be known during planning.")
		return
	}

	region := strings.TrimSpace(plan.Region.ValueString())
	if region == "" {
		resp.Diagnostics.AddError("Invalid region", "`region` must not be empty.")
		return
	}

	networkType := strings.TrimSpace(plan.NetworkType.ValueString())
	if networkType == "" {
		networkType = "BGP"
	}

	out, err := r.client.SelectBandwidthPackage(ctx, region, networkType)
	if err != nil {
		resp.Diagnostics.AddError("Failed to select bandwidth package", err.Error())
		return
	}

	state := TencentCloudBandwidthPackageSelectionResourceModel{
		ID:                 types.StringValue(out.ID),
		Region:             types.StringValue(region),
		NetworkType:        types.StringValue(networkType),
		BandwidthPackageID: types.StringValue(out.ID),
		AvailableCount:     types.Int64Value(out.AvailableCount),
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudBandwidthPackageSelectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state TencentCloudBandwidthPackageSelectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Snapshot semantics: keep state stable; no API call on refresh.
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *TencentCloudBandwidthPackageSelectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan TencentCloudBandwidthPackageSelectionResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (r *TencentCloudBandwidthPackageSelectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No remote object to delete; selection is only stored in state.
}
