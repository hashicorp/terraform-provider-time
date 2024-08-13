// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	testingresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"time": providerserver.NewProtocol5WithError(New()),
	}
}

func protoV5ProviderFactoriesTestProvider(testClock clock.Clock) map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"time": providerserver.NewProtocol5WithError(NewTestProvider(testClock)),
	}
}

func providerVersion080() map[string]testingresource.ExternalProvider {
	return map[string]testingresource.ExternalProvider{
		"time": {
			VersionConstraint: "0.8.0",
			Source:            "hashicorp/time",
		},
	}
}
