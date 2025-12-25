// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/indexyz/terraform-provider-penguin/internal/penguin"
)

func clientFromProviderData(data any) (*penguin.Client, error) {
	if data == nil {
		return nil, fmt.Errorf("provider not configured")
	}

	client, ok := data.(*penguin.Client)
	if !ok {
		return nil, fmt.Errorf("unexpected provider data type: %T", data)
	}

	return client, nil
}

func configureDataSourceClient(req datasource.ConfigureRequest, resp *datasource.ConfigureResponse, target **penguin.Client) {
	if req.ProviderData == nil {
		return
	}

	client, err := clientFromProviderData(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Data Source Configure Type", err.Error())
		return
	}

	*target = client
}

func configureResourceClient(req resource.ConfigureRequest, resp *resource.ConfigureResponse, target **penguin.Client) {
	if req.ProviderData == nil {
		return
	}

	client, err := clientFromProviderData(req.ProviderData)
	if err != nil {
		resp.Diagnostics.AddError("Unexpected Resource Configure Type", err.Error())
		return
	}

	*target = client
}
