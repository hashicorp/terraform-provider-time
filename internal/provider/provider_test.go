// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"time": providerserver.NewProtocol5WithError(New()),
	}
}

func providerVersion080() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"time": {
			VersionConstraint: "0.8.0",
			Source:            "hashicorp/time",
		},
	}
}
