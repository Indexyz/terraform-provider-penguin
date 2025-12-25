// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import "github.com/indexyz/terraform-provider-penguin/internal/penguin"

func apiErrorStatus(err error) (int, bool) {
	apiErr, ok := err.(*penguin.APIError)
	if !ok || apiErr == nil {
		return 0, false
	}
	return apiErr.Status, true
}

func isNotFound(err error) bool {
	status, ok := apiErrorStatus(err)
	return ok && status == 404
}
