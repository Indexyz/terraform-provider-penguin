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

var _ datasource.DataSource = &TencentCloudZonesDataSource{}

func NewTencentCloudZonesDataSource() datasource.DataSource {
	return &TencentCloudZonesDataSource{}
}

type TencentCloudZonesDataSource struct {
	client *penguin.Client
}

type TencentCloudZonesDataSourceModel struct {
	ID    types.String                 `tfsdk:"id"`
	Zones []TencentCloudZonesZoneModel `tfsdk:"zones"`
}

type TencentCloudZonesZoneModel struct {
	Region     types.String `tfsdk:"region"`
	RegionName types.String `tfsdk:"region_name"`
	Zone       types.String `tfsdk:"zone"`
	ZoneName   types.String `tfsdk:"zone_name"`
	ZoneID     types.String `tfsdk:"zone_id"`
	State      types.String `tfsdk:"state"`
}

func (d *TencentCloudZonesDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tencentcloud_zones"
}

func (d *TencentCloudZonesDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "List available Tencent Cloud zones from the Penguin service.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"zones": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"region": schema.StringAttribute{
							Computed: true,
						},
						"region_name": schema.StringAttribute{
							Computed: true,
						},
						"zone": schema.StringAttribute{
							Computed: true,
						},
						"zone_name": schema.StringAttribute{
							Computed: true,
						},
						"zone_id": schema.StringAttribute{
							Computed: true,
						},
						"state": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

func (d *TencentCloudZonesDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	configureDataSourceClient(req, resp, &d.client)
}

func (d *TencentCloudZonesDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	if d.client == nil {
		resp.Diagnostics.AddError("Unconfigured provider", "The provider has not been configured.")
		return
	}

	zones, err := d.client.ListZones(ctx)
	if err != nil {
		resp.Diagnostics.AddError("Failed to list zones", err.Error())
		return
	}

	state := TencentCloudZonesDataSourceModel{
		ID:    types.StringValue("zones"),
		Zones: make([]TencentCloudZonesZoneModel, 0, len(zones)),
	}

	for _, zone := range zones {
		state.Zones = append(state.Zones, TencentCloudZonesZoneModel{
			Region:     types.StringValue(zone.Region),
			RegionName: types.StringValue(zone.RegionName),
			Zone:       types.StringValue(zone.Zone),
			ZoneName:   types.StringValue(zone.ZoneName),
			ZoneID:     types.StringValue(zone.ZoneID),
			State:      types.StringValue(zone.State),
		})
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}
