// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	resourceschema "github.com/hashicorp/terraform-plugin-framework/resource/schema"
	schemavalidator "github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestTencentCloudBandwidthPackageSelectionResourceSchemaValidator(t *testing.T) {
	t.Parallel()

	r := &TencentCloudBandwidthPackageSelectionResource{}
	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	attr, ok := resp.Schema.Attributes["shared_bandwidth_package_id"]
	if !ok {
		t.Fatal("shared_bandwidth_package_id attribute not found")
	}

	stringAttr, ok := attr.(resourceschema.StringAttribute)
	if !ok {
		t.Fatalf("unexpected attribute type: %T", attr)
	}
	if len(stringAttr.Validators) == 0 {
		t.Fatal("expected shared_bandwidth_package_id to have at least one validator")
	}

	stringValidator := stringAttr.Validators[0]

	t.Run("accepts canonical value", func(t *testing.T) {
		var validateResp schemavalidator.StringResponse
		stringValidator.ValidateString(context.Background(), schemavalidator.StringRequest{
			Path:        path.Root("shared_bandwidth_package_id"),
			ConfigValue: types.StringValue("@blue"),
		}, &validateResp)

		if validateResp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %#v", validateResp.Diagnostics)
		}
	})

	t.Run("rejects surrounding whitespace", func(t *testing.T) {
		var validateResp schemavalidator.StringResponse
		stringValidator.ValidateString(context.Background(), schemavalidator.StringRequest{
			Path:        path.Root("shared_bandwidth_package_id"),
			ConfigValue: types.StringValue(" @blue "),
		}, &validateResp)

		if !validateResp.Diagnostics.HasError() {
			t.Fatal("expected diagnostics for whitespace-padded selector")
		}
	})

	t.Run("accepts empty string as omitted", func(t *testing.T) {
		var validateResp schemavalidator.StringResponse
		stringValidator.ValidateString(context.Background(), schemavalidator.StringRequest{
			Path:        path.Root("shared_bandwidth_package_id"),
			ConfigValue: types.StringValue(""),
		}, &validateResp)

		if validateResp.Diagnostics.HasError() {
			t.Fatalf("unexpected diagnostics: %#v", validateResp.Diagnostics)
		}
	})
}

func TestCanonicalSharedBandwidthPackageSelector(t *testing.T) {
	t.Parallel()

	t.Run("null stays unset", func(t *testing.T) {
		got, err := canonicalSharedBandwidthPackageSelector(types.StringNull())
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "" {
			t.Fatalf("expected empty selector, got %q", got)
		}
	})

	t.Run("canonical value accepted", func(t *testing.T) {
		got, err := canonicalSharedBandwidthPackageSelector(types.StringValue("@blue"))
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if got != "@blue" {
			t.Fatalf("unexpected selector: %q", got)
		}
	})

	t.Run("empty string treated as omitted", func(t *testing.T) {
		got, err := canonicalSharedBandwidthPackageSelector(types.StringValue(""))
		if err == nil {
			if got != "" {
				t.Fatalf("expected empty selector, got %q", got)
			}
			return
		}
		t.Fatalf("unexpected error: %v", err)
	})

	t.Run("surrounding whitespace rejected", func(t *testing.T) {
		_, err := canonicalSharedBandwidthPackageSelector(types.StringValue(" @blue "))
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("whitespace only rejected", func(t *testing.T) {
		_, err := canonicalSharedBandwidthPackageSelector(types.StringValue("   "))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
