// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"code.cloudfoundry.org/clock"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"

	testingresource "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var (
	_ provider.ProviderWithFunctions = (*testTimeProvider)(nil)
)

func NewTestProvider(clock clock.Clock) provider.Provider {
	return &testTimeProvider{
		clock: clock,
	}
}

type testTimeProvider struct {
	clock clock.Clock
}

func (p *testTimeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "time"
}

func (p *testTimeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
}

func (p *testTimeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *testTimeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTimeOffsetResource,
		p.NewTestTimeRotatingResource,
		NewTimeSleepResource,
		NewTimeStaticResource,
	}
}

func (p *testTimeProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}

func (p *testTimeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewRFC3339ParseFunction,
	}
}

func (p *testTimeProvider) NewTestTimeRotatingResource() resource.Resource {
	return &timeRotatingResource{
		clock: p.clock,
	}
}

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
