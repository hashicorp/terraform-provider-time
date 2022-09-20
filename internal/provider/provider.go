package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	p "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func New() p.Provider {
	return &provider{}
}

var (
	_ p.Provider             = (*provider)(nil)
	_ p.ProviderWithMetadata = (*provider)(nil)
)

type provider struct{}

func (p *provider) Metadata(ctx context.Context, req p.MetadataRequest, resp *p.MetadataResponse) {
	resp.TypeName = "time"
}

func (p *provider) Configure(ctx context.Context, req p.ConfigureRequest, resp *p.ConfigureResponse) {

}

func (p *provider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return nil
}

func (p *provider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTimeOffsetResource,
		NewTimeRotatingResource,
		NewTimeSleepResource,
		NewTimeStaticResource,
	}
}

func (p *provider) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{}, nil
}
