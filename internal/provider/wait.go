// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"
)

func waitUntil(ctx context.Context, interval time.Duration, check func(context.Context) (done bool, err error)) error {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		done, err := check(ctx)
		if err != nil {
			return err
		}
		if done {
			return nil
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
