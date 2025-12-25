// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/indexyz/terraform-provider-penguin/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
	commit string = "none"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	addr := strings.TrimSpace(os.Getenv("PENGUIN_PROVIDER_ADDRESS"))
	if addr == "" {
		addr = "registry.terraform.io/indexyz/penguin"
	}

	opts := providerserver.ServeOpts{
		Address: addr,
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), provider.New(version, commit), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
