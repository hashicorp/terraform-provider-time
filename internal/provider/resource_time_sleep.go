// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
)

var (
	_ resource.Resource                = (*timeSleepResource)(nil)
	_ resource.ResourceWithImportState = (*timeSleepResource)(nil)
)

func NewTimeSleepResource() resource.Resource {
	return &timeSleepResource{}
}

type timeSleepResource struct{}

func (t timeSleepResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sleep"
}

func (t timeSleepResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a resource that delays creation and/or destruction, typically for further resources. " +
			"This prevents cross-platform compatibility and destroy-time issues with using " +
			"the [`local-exec` provisioner](https://www.terraform.io/docs/provisioners/local-exec.html).",
		Attributes: map[string]schema.Attribute{
			"create_duration": schema.StringAttribute{
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource creation. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("destroy_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"destroy_duration": schema.StringAttribute{
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource destroy. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay. " +
					"This value or any updates to it must be successfully applied into the Terraform state before destroying this resource to take effect.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("create_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"triggers": schema.MapAttribute{
				Description: "(Optional) Arbitrary map of values that, when changed, will run any creation or destroy delays again. " +
					"See [the main provider documentation](../index.md) for more information.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"id": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (t timeSleepResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID

	idParts := strings.Split(id, ",")

	if len(idParts) != 2 || (idParts[0] == "" && idParts[1] == "") {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected CREATEDURATION,DESTROYDURATION where at least one value is non-empty", id))

		return
	}

	state := timeSleepModelV0{
		CreateDuration:  types.StringNull(),
		DestroyDuration: types.StringNull(),
		ID:              timetypes.NewRFC3339Value(time.Now().UTC().Format(time.RFC3339)),
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
		state.CreateDuration = types.StringValue(idParts[0])
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
		state.DestroyDuration = types.StringValue(idParts[1])
	}

	state.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeSleepResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeSleepModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.CreateDuration.ValueString() != "" {
		duration, err := time.ParseDuration(plan.CreateDuration.ValueString())
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
		ID:              timetypes.NewRFC3339Value(time.Now().UTC().Format(time.RFC3339)),
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeSleepResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (t timeSleepResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data timeSleepModelV0

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t timeSleepResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state timeSleepModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.DestroyDuration.ValueString() != "" {
		duration, err := time.ParseDuration(state.DestroyDuration.ValueString())
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

type timeSleepModelV0 struct {
	CreateDuration  types.String      `tfsdk:"create_duration"`
	DestroyDuration types.String      `tfsdk:"destroy_duration"`
	Triggers        types.Map         `tfsdk:"triggers"`
	ID              timetypes.RFC3339 `tfsdk:"id"`
}
