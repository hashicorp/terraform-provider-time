// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"log"
	"os"

	timeprovider "github.com/hashicorp/terraform-provider-time/internal/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	var debug, metadataJSON bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.BoolVar(&metadataJSON, "metadata-json", false, "set to true to output provider metadata JSON information instead of starting the provider server")
	flag.Parse()

	if metadataJSON {
		metadata := provider.BuildMetadata(context.TODO(), timeprovider.New())

		// Ensure warnings are also written to stderr
		/* if len(diags) > 0 {
		//	fmt.Fprint(os.Stderr /* diags textual form */ //)
		/* }

		if diags.HasError() {
			os.Exit(1)
		}
		*/

		fmt.Fprint(os.Stdout, metadata)
		os.Exit(0)
	}

	err := providerserver.Serve(context.Background(), timeprovider.New, providerserver.ServeOpts{
		Address:         "registry.terraform.io/hashicorp/time",
		Debug:           debug,
		ProtocolVersion: 5,
	})
	if err != nil {
		log.Fatal(err)
	}
}
