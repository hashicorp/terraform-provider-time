package provider

import (
	"context"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

var _ tfsdk.ResourceType = (*timeStaticResourceType)(nil)

type timeStaticResourceType struct{}

func (t timeStaticResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	//TODO implement me
	panic("implement me")
}

func (t timeStaticResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	//TODO implement me
	panic("implement me")
}

var (
	_ tfsdk.Resource = (*timeStaticResource)(nil)
)

type timeStaticResource struct {
}

func (t timeStaticResource) Create(ctx context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeStaticResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeStaticResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeStaticResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	//TODO implement me
	panic("implement me")
}
