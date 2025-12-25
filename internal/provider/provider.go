// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

// Ensure PenguinProvider satisfies various provider interfaces.
var _ provider.Provider = &PenguinProvider{}

// PenguinProvider defines the provider implementation.
type PenguinProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
	commit  string
}

// PenguinProviderModel describes the provider configuration model.
type PenguinProviderModel struct {
	Endpoint  types.String `tfsdk:"endpoint"`
	AuthToken types.String `tfsdk:"auth_token"`
	JWT       types.String `tfsdk:"jwt"`
}

func (p *PenguinProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "penguin"
	resp.Version = p.version
}

func (p *PenguinProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Penguin service base URL, e.g. `http://127.0.0.1:8080`. Can also be set via `PENGUIN_ENDPOINT`.",
				Optional:            true,
			},
			"auth_token": schema.StringAttribute{
				MarkdownDescription: "Legacy bearer token used to authenticate to Penguin. Can also be set via `PENGUIN_AUTH_TOKEN`.",
				Optional:            true,
				Sensitive:           true,
			},
			"jwt": schema.StringAttribute{
				MarkdownDescription: "Optional JWT to enforce provisioning limits. Can also be set via `PENGUIN_JWT`.",
				Optional:            true,
				Sensitive:           true,
			},
		},
	}
}

func (p *PenguinProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data PenguinProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Endpoint.IsUnknown() || data.AuthToken.IsUnknown() || data.JWT.IsUnknown() {
		resp.Diagnostics.AddError(
			"Unknown Penguin Provider Configuration",
			"Provider configuration values must be known during planning. Check for unknown values in `endpoint`, `auth_token`, or `jwt`.",
		)
		return
	}

	endpoint := strings.TrimSpace(data.Endpoint.ValueString())
	if endpoint == "" {
		endpoint = strings.TrimSpace(os.Getenv("PENGUIN_ENDPOINT"))
	}
	legacyToken := strings.TrimSpace(data.AuthToken.ValueString())
	if legacyToken == "" {
		legacyToken = strings.TrimSpace(os.Getenv("PENGUIN_AUTH_TOKEN"))
	}
	jwt := strings.TrimSpace(data.JWT.ValueString())
	if jwt == "" {
		jwt = strings.TrimSpace(os.Getenv("PENGUIN_JWT"))
	}

	if endpoint == "" {
		resp.Diagnostics.AddError(
			"Missing Penguin Endpoint",
			"Set provider attribute `endpoint` or environment variable `PENGUIN_ENDPOINT`.",
		)
		return
	}

	client, err := penguin.NewClient(endpoint, legacyToken, jwt, penguin.ClientOptions{
		UserAgent: fmt.Sprintf("terraform-provider-penguin/%s (%s)", p.version, p.commit),
	})
	if err != nil {
		resp.Diagnostics.AddError("Failed to create Penguin client", err.Error())
		return
	}

	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *PenguinProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTencentCloudVirtualMachineResource,
		NewTencentCloudElasticIPResource,
		NewTencentCloudBandwidthPackageSelectionResource,
	}
}

func (p *PenguinProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewTencentCloudZonesDataSource,
		NewTencentCloudBandwidthPackageDataSource,
		NewTencentCloudVirtualMachineStatusDataSource,
		NewTencentCloudVirtualMachineMetricsDataSource,
		NewTencentCloudVirtualMachineVNCDataSource,
		NewInternalHealthDataSource,
		NewJWTDataSource,
	}
}

func New(version string, commit string) func() provider.Provider {
	return func() provider.Provider {
		return &PenguinProvider{
			version: version,
			commit:  commit,
		}
	}
}
