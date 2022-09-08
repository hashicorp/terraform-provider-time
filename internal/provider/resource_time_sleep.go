package provider

import (
	"context"
	"regexp"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.ResourceType = (*timeSleepResourceType)(nil)

type timeSleepResourceType struct{}

func (t timeSleepResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a resource that delays creation and/or destruction, typically for further resources. " +
			"This prevents cross-platform compatibility and destroy-time issues with using " +
			"the [`local-exec` provisioner](https://www.terraform.io/docs/provisioners/local-exec.html).",
		Attributes: map[string]tfsdk.Attribute{
			"create_duration": {
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource creation. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay.",
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("destroy_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"destroy_duration": {
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource destroy. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay. " +
					"This value or any updates to it must be successfully applied into the Terraform state before destroying this resource to take effect.",
				Type:     types.StringType,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("create_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"triggers": {
				Description: "(Optional) Arbitrary map of values that, when changed, will run any creation or destroy delays again. " +
					"See [the main provider documentation](../index.md) for more information.",
				Type:     types.MapType{ElemType: types.StringType},
				Optional: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"id": {
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (t timeSleepResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	//TODO implement me
	panic("implement me")
}

var (
	_ tfsdk.Resource                = (*timeSleepResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeSleepResource)(nil)
)

type timeSleepResource struct {
}

func (t timeSleepResource) ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Create(ctx context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeSleepResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	//TODO implement me
	panic("implement me")
}
