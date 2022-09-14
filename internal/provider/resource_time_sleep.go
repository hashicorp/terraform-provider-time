package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

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
	return &timeSleepResource{}, nil
}

var (
	_ tfsdk.Resource                = (*timeSleepResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeSleepResource)(nil)
)

type timeSleepResource struct {
}

type timeSleepModelV0 struct {
	CreateDuration  types.String `tfsdk:"create_duration"`
	DestroyDuration types.String `tfsdk:"destroy_duration"`
	Triggers        types.Map    `tfsdk:"triggers"`
	ID              types.String `tfsdk:"id"`
}

func (t timeSleepResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	id := req.ID

	idParts := strings.Split(id, ",")

	if len(idParts) != 2 || (idParts[0] == "" && idParts[1] == "") {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected CREATEDURATION,DESTROYDURATION where at least one value is non-empty", id))

		return
	}

	state := timeSleepModelV0{
		CreateDuration:  types.String{Null: true},
		DestroyDuration: types.String{Null: true},
		ID:              types.String{Value: time.Now().UTC().Format(time.RFC3339)},
	}

	if idParts[0] != "" {
		_, err := time.ParseDuration(idParts[0])
		if err != nil {
			resp.Diagnostics.AddError(
				"Import time sleep error",
				"The create_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
		state.CreateDuration.Null = false
		state.CreateDuration.Value = idParts[0]
	}

	if idParts[1] != "" {
		_, err := time.ParseDuration(idParts[1])
		if err != nil {
			resp.Diagnostics.AddError(
				"Import time sleep error",
				"The create_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
		state.DestroyDuration.Null = false
		state.DestroyDuration.Value = idParts[1]
	}

	state.Triggers.ElemType = types.StringType
	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

func (t timeSleepResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan timeSleepModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CreateDuration.Value != "" {
		duration, err := time.ParseDuration(plan.CreateDuration.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Create time sleep error",
				"The create_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Create time sleep error",
				fmt.Sprintf("Original Error: %s", ctx.Err()),
			)
			return
		case <-time.After(duration):
		}
	}

	state := timeSleepModelV0{
		CreateDuration:  plan.CreateDuration,
		DestroyDuration: plan.DestroyDuration,
		Triggers:        plan.Triggers,
		ID:              types.String{Value: time.Now().UTC().Format(time.RFC3339)},
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeSleepResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

}

func (t timeSleepResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

}

func (t timeSleepResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state timeSleepModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.DestroyDuration.Value != "" {
		duration, err := time.ParseDuration(state.DestroyDuration.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Delete time sleep error",
				"The create_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Delete time sleep error",
				fmt.Sprintf("Original Error: %s", ctx.Err()),
			)
			return
		case <-time.After(duration):
		}
	}
}
