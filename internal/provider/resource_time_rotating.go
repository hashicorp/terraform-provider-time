// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
	"github.com/hashicorp/terraform-provider-time/internal/modifiers/timemodifier"
)

var (
	_ resource.Resource                     = (*timeRotatingResource)(nil)
	_ resource.ResourceWithImportState      = (*timeRotatingResource)(nil)
	_ resource.ResourceWithModifyPlan       = (*timeRotatingResource)(nil)
	_ resource.ResourceWithConfigValidators = (*timeRotatingResource)(nil)
	_ resource.ResourceWithConfigure        = (*timeRotatingResource)(nil)
)

func NewTimeRotatingResource() resource.Resource {
	return &timeRotatingResource{}
}

type timeRotatingResource struct {
	clock clock.Clock
}

func (t *timeRotatingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (t *timeRotatingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rotating"
}

func (t *timeRotatingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a rotating time resource, which keeps a rotating UTC timestamp stored in the Terraform " +
			"state and proposes resource recreation when the locally sourced current time is beyond the rotation time. " +
			"This rotation only occurs when Terraform is executed, meaning there will be drift between the rotation " +
			"timestamp and actual rotation. The new rotation timestamp offset includes this drift. " +
			"This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html) " +
			"by only forcing a new value on the set cadence.",
		Attributes: map[string]schema.Attribute{
			"day": schema.Int64Attribute{
				Description: "Number day of timestamp.",
				Computed:    true,
			},
			"rotation_days": schema.Int64Attribute{
				Description: "Number of days to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_hours": schema.Int64Attribute{
				Description: "Number of hours to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_minutes": schema.Int64Attribute{
				Description: "Number of minutes to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_months": schema.Int64Attribute{
				Description: "Number of months to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_rfc3339": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Description: "Configure the rotation timestamp with an " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format of the offset timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplaceIf(
						timemodifier.ReplaceIfOutdated,
						"resource will be replaced if current time is past the saved time",
						"resource will be replaced if current time is past the saved time"),
				},
			},
			"rotation_years": schema.Int64Attribute{
				Description: "Number of years to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"hour": schema.Int64Attribute{
				Description: "Number hour of timestamp.",
				Computed:    true,
			},
			"triggers": schema.MapAttribute{
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved." +
					" These conditions recreate the resource in addition to other rotation arguments. " +
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
				Description: "RFC3339 format of the timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Computed:    true,
			},
		},
	}
}

func (t *timeRotatingResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("rotation_minutes"),
			path.MatchRoot("rotation_hours"),
			path.MatchRoot("rotation_days"),
			path.MatchRoot("rotation_months"),
			path.MatchRoot("rotation_years"),
			path.MatchRoot("rotation_rfc3339"),
		),
	}
}

func (t *timeRotatingResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Plan does not need to be modified when the resource is being destroyed.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Plan only needs modifying if the resource already exists as the purpose of
	// the plan modifier is to show updated attribute values on CLI.
	if req.State.Raw.IsNull() {
		return
	}

	var state, plan timeRotatingModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.RotationYears == plan.RotationYears &&
		state.RotationMonths == plan.RotationMonths &&
		state.RotationDays == plan.RotationDays &&
		state.RotationHours == plan.RotationHours &&
		state.RotationMinutes == plan.RotationMinutes &&
		state.RotationRFC3339 == plan.RotationRFC3339 {
		return
	}

	var RFC3339, rotationRFC3339 timetypes.RFC3339

	diags = req.Plan.GetAttribute(ctx, path.Root("rfc3339"), &RFC3339)
	resp.Diagnostics = append(resp.Diagnostics, diags...)

	diags = req.Plan.GetAttribute(ctx, path.Root("rotation_rfc3339"), &rotationRFC3339)
	resp.Diagnostics = append(resp.Diagnostics, diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// RFC3339 and rotationRFC3339 could be unknown if there is no value set in the config as the attribute is
	// optional and computed. If base_rfc3339 is not set in config then the previous value from
	// state is used and propagated to the update function.
	if RFC3339.IsUnknown() {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("rfc3339"), &RFC3339)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	if rotationRFC3339.IsUnknown() {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("rotation_rfc3339"), &rotationRFC3339)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	timestamp, diags := RFC3339.ValueRFC3339Time()

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(setRotationValues(&plan, timestamp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.Plan.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	var state timeRotatingModelV0
	var err error

	idParts := strings.Split(id, ",")

	if len(idParts) != 2 && len(idParts) != 6 {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES or BASETIMESTAMP,ROTATIONTIMESTAMP", id))

		return
	}

	if len(idParts) == 2 {
		if idParts[0] == "" || idParts[1] == "" {
			resp.Diagnostics.AddError(
				"Unexpected Format of ID",
				fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,ROTATIONTIMESTAMP", id))
			return
		}

		state, err = parseTwoPartId(idParts)
		if err != nil {
			resp.Diagnostics.AddError(
				"Import time rotating error",
				"The timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

	} else {
		if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "") {
			resp.Diagnostics.AddError(
				"Unexpected Format of ID",
				fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES where at least one rotation value is non-empty", id))

			return
		}
		state, err = parseMultiplePartId(idParts)
		if err != nil {
			resp.Diagnostics.AddError(
				"Import time rotating error",
				"The parameter that was supplied could not be parsed.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	state.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeRotatingModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := t.clock.Now().UTC()

	if plan.RFC3339.ValueString() != "" {
		rfc3339, diags := plan.RFC3339.ValueRFC3339Time()

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		timestamp = rfc3339
	}

	resp.Diagnostics.Append(setRotationValues(&plan, timestamp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state timeRotatingModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RotationRFC3339.IsNull() && state.RotationRFC3339.ValueString() != "" {
		now := t.clock.Now().UTC()
		rotationTimestamp, diags := state.RotationRFC3339.ValueRFC3339Time()

		resp.Diagnostics.Append(diags...)

		if resp.Diagnostics.HasError() {
			return
		}

		if now.After(rotationTimestamp) {
			log.Printf("[INFO] Expiration timestamp (%s) is after current timestamp (%s), removing from state", state.RotationRFC3339.ValueString(), now.Format(time.RFC3339))
			resp.State.RemoveResource(ctx)
			return
		}
	}

}

func (t *timeRotatingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state timeRotatingModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if state.RotationYears == plan.RotationYears &&
		state.RotationMonths == plan.RotationMonths &&
		state.RotationDays == plan.RotationDays &&
		state.RotationHours == plan.RotationHours &&
		state.RotationMinutes == plan.RotationMinutes &&
		state.RotationRFC3339 == plan.RotationRFC3339 {
		return
	}

	timestamp, diags := plan.ID.ValueRFC3339Time()

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(setRotationValues(&plan, timestamp)...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingResource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {

}

type timeRotatingModelV0 struct {
	Day             types.Int64       `tfsdk:"day"`
	RotationDays    types.Int64       `tfsdk:"rotation_days"`
	RotationHours   types.Int64       `tfsdk:"rotation_hours"`
	RotationMinutes types.Int64       `tfsdk:"rotation_minutes"`
	RotationMonths  types.Int64       `tfsdk:"rotation_months"`
	RotationRFC3339 timetypes.RFC3339 `tfsdk:"rotation_rfc3339"`
	RotationYears   types.Int64       `tfsdk:"rotation_years"`
	Hour            types.Int64       `tfsdk:"hour"`
	Triggers        types.Map         `tfsdk:"triggers"`
	Minute          types.Int64       `tfsdk:"minute"`
	Month           types.Int64       `tfsdk:"month"`
	RFC3339         timetypes.RFC3339 `tfsdk:"rfc3339"`
	Second          types.Int64       `tfsdk:"second"`
	Unix            types.Int64       `tfsdk:"unix"`
	Year            types.Int64       `tfsdk:"year"`
	ID              timetypes.RFC3339 `tfsdk:"id"`
}

func setRotationValues(plan *timeRotatingModelV0, timestamp time.Time) diag.Diagnostics {
	var diags diag.Diagnostics
	var rotationTimestamp time.Time

	if plan.RotationDays.ValueInt64() != 0 {
		rotationTimestamp = timestamp.AddDate(0, 0, int(plan.RotationDays.ValueInt64()))
	}

	if plan.RotationHours.ValueInt64() != 0 {
		hours := time.Duration(plan.RotationHours.ValueInt64()) * time.Hour
		rotationTimestamp = timestamp.Add(hours)
	}

	if plan.RotationMinutes.ValueInt64() != 0 {
		minutes := time.Duration(plan.RotationMinutes.ValueInt64()) * time.Minute
		rotationTimestamp = timestamp.Add(minutes)
	}

	if plan.RotationMonths.ValueInt64() != 0 {
		rotationTimestamp = timestamp.AddDate(0, int(plan.RotationMonths.ValueInt64()), 0)
	}

	if plan.RotationRFC3339.ValueString() != "" {
		rotationTimestamp, diags = plan.RotationRFC3339.ValueRFC3339Time()
	}

	if plan.RotationYears.ValueInt64() != 0 {
		rotationTimestamp = timestamp.AddDate(int(plan.RotationYears.ValueInt64()), 0, 0)
	}

	plan.RotationRFC3339 = timetypes.NewRFC3339TimeValue(rotationTimestamp)
	plan.Year = types.Int64Value(int64(rotationTimestamp.Year()))
	plan.Month = types.Int64Value(int64(rotationTimestamp.Month()))
	plan.Day = types.Int64Value(int64(rotationTimestamp.Day()))
	plan.Hour = types.Int64Value(int64(rotationTimestamp.Hour()))
	plan.Minute = types.Int64Value(int64(rotationTimestamp.Minute()))
	plan.Second = types.Int64Value(int64(rotationTimestamp.Second()))
	plan.RFC3339 = timetypes.NewRFC3339TimeValue(timestamp)
	plan.Unix = types.Int64Value(rotationTimestamp.Unix())
	plan.ID = timetypes.NewRFC3339TimeValue(timestamp)

	return diags
}

func parseTwoPartId(idParts []string) (timeRotatingModelV0, error) {

	baseRfc3339 := idParts[0]
	rotationRfc3339 := idParts[1]

	timestamp, err := time.Parse(time.RFC3339, baseRfc3339)
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	rotationTimestamp, err := time.Parse(time.RFC3339, rotationRfc3339)
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	return timeRotatingModelV0{
		Year:            types.Int64Value(int64(rotationTimestamp.Year())),
		Month:           types.Int64Value(int64(rotationTimestamp.Month())),
		Day:             types.Int64Value(int64(rotationTimestamp.Day())),
		Hour:            types.Int64Value(int64(rotationTimestamp.Hour())),
		Minute:          types.Int64Value(int64(rotationTimestamp.Minute())),
		Second:          types.Int64Value(int64(rotationTimestamp.Second())),
		RotationRFC3339: timetypes.NewRFC3339TimeValue(rotationTimestamp),
		RotationYears:   types.Int64Null(),
		RotationMonths:  types.Int64Null(),
		RotationDays:    types.Int64Null(),
		RotationHours:   types.Int64Null(),
		RotationMinutes: types.Int64Null(),
		RFC3339:         timetypes.NewRFC3339TimeValue(timestamp),
		Unix:            types.Int64Value(rotationTimestamp.Unix()),
		ID:              timetypes.NewRFC3339TimeValue(timestamp),
	}, nil
}

func parseMultiplePartId(idParts []string) (timeRotatingModelV0, error) {
	baseRfc3339 := idParts[0]

	rotationYears, err := rotationToInt64(idParts[1])
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	rotationMonths, err := rotationToInt64(idParts[2])
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	rotationDays, err := rotationToInt64(idParts[3])
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	rotationHours, err := rotationToInt64(idParts[4])
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	rotationMinutes, err := rotationToInt64(idParts[5])
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	timestamp, err := time.Parse(time.RFC3339, baseRfc3339)
	if err != nil {
		return timeRotatingModelV0{}, err
	}

	var rotationTimestamp time.Time

	if !rotationDays.IsNull() && rotationDays.ValueInt64() > 0 {
		rotationTimestamp = timestamp.AddDate(0, 0, int(rotationDays.ValueInt64()))
	}

	if !rotationHours.IsNull() && rotationHours.ValueInt64() > 0 {
		hours := time.Duration(rotationHours.ValueInt64()) * time.Hour
		rotationTimestamp = timestamp.Add(hours)
	}

	if !rotationMinutes.IsNull() && rotationMinutes.ValueInt64() > 0 {
		minutes := time.Duration(rotationMinutes.ValueInt64()) * time.Minute
		rotationTimestamp = timestamp.Add(minutes)
	}

	if !rotationMonths.IsNull() && rotationMonths.ValueInt64() > 0 {
		rotationTimestamp = timestamp.AddDate(0, int(rotationMonths.ValueInt64()), 0)
	}

	if !rotationYears.IsNull() && rotationYears.ValueInt64() > 0 {
		rotationTimestamp = timestamp.AddDate(int(rotationYears.ValueInt64()), 0, 0)
	}

	state := timeRotatingModelV0{
		Year:            types.Int64Value(int64(rotationTimestamp.Year())),
		Month:           types.Int64Value(int64(rotationTimestamp.Month())),
		Day:             types.Int64Value(int64(rotationTimestamp.Day())),
		Hour:            types.Int64Value(int64(rotationTimestamp.Hour())),
		Minute:          types.Int64Value(int64(rotationTimestamp.Minute())),
		Second:          types.Int64Value(int64(rotationTimestamp.Second())),
		RotationRFC3339: timetypes.NewRFC3339TimeValue(rotationTimestamp),
		RotationYears:   rotationYears,
		RotationMonths:  rotationMonths,
		RotationDays:    rotationDays,
		RotationHours:   rotationHours,
		RotationMinutes: rotationMinutes,
		RFC3339:         timetypes.NewRFC3339TimeValue(timestamp),
		Unix:            types.Int64Value(rotationTimestamp.Unix()),
		ID:              timetypes.NewRFC3339TimeValue(timestamp),
	}

	return state, nil
}

func rotationToInt64(rotationStr string) (types.Int64, error) {
	rotation := types.Int64Null()

	if rotationStr != "" {
		offsetInt, err := strconv.ParseInt(rotationStr, 10, 64)
		if err != nil {
			return rotation, fmt.Errorf("could not parse rotation (%q) as int: %w", rotationStr, err)
		}

		rotation = types.Int64Value(offsetInt)
	}

	return rotation, nil
}
