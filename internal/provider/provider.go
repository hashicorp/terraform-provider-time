// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

func New() provider.Provider {
	return &timeProvider{
		clock.NewClock(),
	}
}

func NewTestProvider(clock clock.Clock) provider.Provider {
	return &timeProvider{
		clock: clock,
	}
}

var (
	_ provider.ProviderWithFunctions = (*timeProvider)(nil)
)

type timeProvider struct {
	clock clock.Clock
}

func (p *timeProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "time"
}

func (p *timeProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	resp.ResourceData = p.clock
}

func (p *timeProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *timeProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTimeOffsetResource,
		NewTimeRotatingResource,
		NewTimeSleepResource,
		NewTimeStaticResource,
	}
}

func (p *timeProvider) Schema(context.Context, provider.SchemaRequest, *provider.SchemaResponse) {
}

func (p *timeProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{
		NewDurationParseFunction,
		NewRFC3339ParseFunction,
		NewUnixTimestampParseFunction,
	}
}
