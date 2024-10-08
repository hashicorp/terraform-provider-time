// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

var (
	_ resource.Resource                = (*timeStaticResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*timeStaticResource)(nil)
	_ resource.ResourceWithImportState = (*timeStaticResource)(nil)
	_ resource.ResourceWithConfigure   = (*timeStaticResource)(nil)
)

func NewTimeStaticResource() resource.Resource {
	return &timeStaticResource{}
}

type timeStaticResource struct {
	clock clock.Clock
}

func (t *timeStaticResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	// Always perform a nil check when handling ProviderData because Terraform
	// sets that data after it calls the ConfigureProvider RPC.
	if req.ProviderData == nil {
		return
	}

	pClock, ok := req.ProviderData.(clock.Clock)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected clock.Clock, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}

	t.clock = pClock
}

func (t *timeStaticResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Skip plan modification unless it's a create operation
	if req.Plan.Raw.IsNull() || !req.State.Raw.IsNull() {
		return
	}

	var plan timeStaticModelV0

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Currently, it is only possible to enhance the plan when the rfc3339 value is defined in configuration (i.e. value is not null and known in plan).
	//
	// Terraform calls the PlanResourceChange RPC twice (initial planned state and final planned state) and currently has no mechanism for sharing information between
	// the initial plan call and final plan call. This means that we can't create a final plan that matches the initial plan using something like time.Now()
	// which will differ between the two calls and result in a "Provider produced inconsistent final plan" error from Terraform.
	//
	// If functionality is introduced in the future that allows us to create consistent final and initial plans, we'd likely want to introduce a new managed resource that
	// always determines its results at plan time. Changing this resource to adopt that behavior would be a breaking change for practitioners who are relying on the time being
	// determined at apply time.
	//
	// There is no time provider feature request currently for this behavior, but a similar long-standing issue exists on the random provider:
	// - https://github.com/hashicorp/terraform-provider-random/issues/121
	if plan.RFC3339.IsNull() || plan.RFC3339.IsUnknown() {
		return
	}

	rfc3339, diags := plan.RFC3339.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.Year = types.Int64Value(int64(rfc3339.Year()))
	plan.Month = types.Int64Value(int64(rfc3339.Month()))
	plan.Day = types.Int64Value(int64(rfc3339.Day()))
	plan.Hour = types.Int64Value(int64(rfc3339.Hour()))
	plan.Minute = types.Int64Value(int64(rfc3339.Minute()))
	plan.Second = types.Int64Value(int64(rfc3339.Second()))
	plan.Unix = types.Int64Value(rfc3339.Unix())
	plan.ID = plan.RFC3339

	resp.Diagnostics.Append(resp.Plan.Set(ctx, &plan)...)
}

func (t *timeStaticResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_static"
}

func (t *timeStaticResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a static time resource, which keeps a locally sourced UTC timestamp stored in the Terraform state. " +
			"This prevents perpetual differences caused by using " +
			"the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).",
		Attributes: map[string]schema.Attribute{
			"day": schema.Int64Attribute{
				Description: "Number day of timestamp.",
				Computed:    true,
			},
			"hour": schema.Int64Attribute{
				Description: "Number hour of timestamp.",
				Computed:    true,
			},
			"triggers": schema.MapAttribute{
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. " +
					"See [the main provider documentation](../index.md) for more information.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"minute": schema.Int64Attribute{
				Description: "Number minute of timestamp.",
				Computed:    true,
			},
			"month": schema.Int64Attribute{
				Description: "Number month of timestamp.",
				Computed:    true,
			},
			"rfc3339": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"second": schema.Int64Attribute{
				Description: "Number second of timestamp.",
				Computed:    true,
			},
			"unix": schema.Int64Attribute{
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Computed:    true,
			},
			"year": schema.Int64Attribute{
				Description: "Number year of timestamp.",
				Computed:    true,
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

func (t *timeStaticResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	timestamp, err := time.Parse(time.RFC3339, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time static error",
			"The id that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	state := timeStaticModelV0{
		Year:    types.Int64Value(int64(timestamp.Year())),
		Month:   types.Int64Value(int64(timestamp.Month())),
		Day:     types.Int64Value(int64(timestamp.Day())),
		Hour:    types.Int64Value(int64(timestamp.Hour())),
		Minute:  types.Int64Value(int64(timestamp.Minute())),
		Second:  types.Int64Value(int64(timestamp.Second())),
		RFC3339: timetypes.NewRFC3339TimeValue(timestamp),
		Unix:    types.Int64Value(timestamp.Unix()),
		ID:      timetypes.NewRFC3339TimeValue(timestamp),
	}
	state.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t *timeStaticResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeStaticModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := t.clock.Now().UTC()

	if !plan.RFC3339.IsNull() && !plan.RFC3339.IsUnknown() {
		rfc3339, diags := plan.RFC3339.ValueRFC3339Time()

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		timestamp = rfc3339
	}

	state := timeStaticModelV0{
		Triggers: plan.Triggers,
		Year:     types.Int64Value(int64(timestamp.Year())),
		Month:    types.Int64Value(int64(timestamp.Month())),
		Day:      types.Int64Value(int64(timestamp.Day())),
		Hour:     types.Int64Value(int64(timestamp.Hour())),
		Minute:   types.Int64Value(int64(timestamp.Minute())),
		Second:   types.Int64Value(int64(timestamp.Second())),
		RFC3339:  timetypes.NewRFC3339TimeValue(timestamp),
		Unix:     types.Int64Value(timestamp.Unix()),
		ID:       timetypes.NewRFC3339TimeValue(timestamp),
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t *timeStaticResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (t *timeStaticResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data timeStaticModelV0

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t *timeStaticResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

type timeStaticModelV0 struct {
	Day      types.Int64       `tfsdk:"day"`
	Hour     types.Int64       `tfsdk:"hour"`
	Triggers types.Map         `tfsdk:"triggers"`
	Minute   types.Int64       `tfsdk:"minute"`
	Month    types.Int64       `tfsdk:"month"`
	RFC3339  timetypes.RFC3339 `tfsdk:"rfc3339"`
	Second   types.Int64       `tfsdk:"second"`
	Unix     types.Int64       `tfsdk:"unix"`
	Year     types.Int64       `tfsdk:"year"`
	ID       timetypes.RFC3339 `tfsdk:"id"`
}
