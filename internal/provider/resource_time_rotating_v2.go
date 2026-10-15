// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"log"
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
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

var (
	_ resource.Resource                     = (*timeRotatingV2Resource)(nil)
	_ resource.ResourceWithImportState      = (*timeRotatingV2Resource)(nil)
	_ resource.ResourceWithModifyPlan       = (*timeRotatingV2Resource)(nil)
	_ resource.ResourceWithConfigValidators = (*timeRotatingV2Resource)(nil)
	_ resource.ResourceWithConfigure        = (*timeRotatingV2Resource)(nil)
)

func NewTimeRotatingV2Resource() resource.Resource {
	return &timeRotatingV2Resource{}
}

type timeRotatingV2Resource struct {
	clock clock.Clock
}

func (t *timeRotatingV2Resource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (t *timeRotatingV2Resource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_rotating_v2"
}

func (t *timeRotatingV2Resource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages a rotating time resource, which keeps a rotating UTC timestamp stored in the Terraform " +
			"state and triggers resource replacement when the locally sourced current time is beyond the rotation time. " +
			"This is an improved version of time_rotating that properly integrates with Terraform's replace_triggered_by " +
			"lifecycle argument. The rotation_* units are cumulative (e.g., rotation_years=1 and rotation_days=1 means " +
			"1 year plus 1 day). Drift behavior is controlled by the first_rotation_rfc3339 attribute.",
		Attributes: map[string]schema.Attribute{
			"day": schema.Int64Attribute{
				Description: "Number day of next rotation timestamp.",
				Computed:    true,
			},
			"rotation_days": schema.Int64Attribute{
				Description: "Number of days to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger replacement. " +
					"At least one of the 'rotation_' arguments must be configured. " +
					"All rotation units are cumulative.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_hours": schema.Int64Attribute{
				Description: "Number of hours to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger replacement. " +
					"At least one of the 'rotation_' arguments must be configured. " +
					"All rotation units are cumulative.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_minutes": schema.Int64Attribute{
				Description: "Number of minutes to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger replacement. " +
					"At least one of the 'rotation_' arguments must be configured. " +
					"All rotation units are cumulative.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_months": schema.Int64Attribute{
				Description: "Number of months to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger replacement. " +
					"At least one of the 'rotation_' arguments must be configured. " +
					"All rotation units are cumulative.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"rotation_years": schema.Int64Attribute{
				Description: "Number of years to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger replacement. " +
					"At least one of the 'rotation_' arguments must be configured. " +
					"All rotation units are cumulative.",
				Optional: true,
				Validators: []validator.Int64{
					int64validator.AtLeast(1),
				},
			},
			"first_rotation_rfc3339": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Description: "Configure the first rotation timestamp with an " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format timestamp. " +
					"When configured, enables drift mode where subsequent rotations advance from actual rotation time. " +
					"When omitted, both first_rotation_rfc3339 and next_rotation_rfc3339 advance together (non-drift mode). " +
					"Defaults to current time plus rotation duration.",
				Optional: true,
				Computed: true,
			},
			"next_rotation_rfc3339": schema.StringAttribute{
				CustomType: timetypes.RFC3339Type{},
				Description: "Computed timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"indicating when the next rotation will occur. When current time exceeds this value, " +
					"the resource triggers replacement.",
				Computed: true,
			},
			"hour": schema.Int64Attribute{
				Description: "Number hour of next rotation timestamp.",
				Computed:    true,
			},
			"triggers": schema.MapAttribute{
				Description: "Arbitrary map of values that, when changed, will trigger a new rotation timestamp to be calculated. " +
					"These conditions recreate the resource in addition to rotation time expiration. " +
					"See [the main provider documentation](../index.md) for more information.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"minute": schema.Int64Attribute{
				Description: "Number minute of next rotation timestamp.",
				Computed:    true,
			},
			"month": schema.Int64Attribute{
				Description: "Number month of next rotation timestamp.",
				Computed:    true,
			},
			"second": schema.Int64Attribute{
				Description: "Number second of next rotation timestamp.",
				Computed:    true,
			},
			"unix": schema.Int64Attribute{
				Description: "Number of seconds since epoch time for next rotation, e.g. `1581489373`.",
				Computed:    true,
			},
			"year": schema.Int64Attribute{
				Description: "Number year of next rotation timestamp.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				CustomType:  timetypes.RFC3339Type{},
				Description: "RFC3339 format of the next rotation timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Computed:    true,
			},
		},
	}
}

func (t *timeRotatingV2Resource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("rotation_minutes"),
			path.MatchRoot("rotation_hours"),
			path.MatchRoot("rotation_days"),
			path.MatchRoot("rotation_months"),
			path.MatchRoot("rotation_years"),
		),
	}
}

func (t *timeRotatingV2Resource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Plan does not need to be modified when the resource is being destroyed.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Plan only needs modifying if the resource already exists.
	if req.State.Raw.IsNull() {
		return
	}

	var state, plan timeRotatingV2Model
	var config timeRotatingV2Model

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

	diags = req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Check if rotation has expired
	now := t.clock.Now().UTC()
	nextRotation, diags := state.NextRotationRFC3339.ValueRFC3339Time()
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if now.After(nextRotation) {
		log.Printf("[INFO] Rotation timestamp (%s) has passed current time (%s), triggering replacement",
			state.NextRotationRFC3339.ValueString(), now.Format(time.RFC3339))

		// Check if user explicitly configured first_rotation_rfc3339
		userConfiguredFirstRotation := !config.FirstRotationRFC3339.IsNull()

		if userConfiguredFirstRotation {
			// DRIFT MODE: next_rotation advances from now, first_rotation stays fixed
			newNextRotation := addRotationDuration(now, plan)
			plan.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(newNextRotation)
			// Keep first_rotation_rfc3339 from state (it was user-configured)
			plan.FirstRotationRFC3339 = state.FirstRotationRFC3339
		} else {
			// NON-DRIFT MODE: both first_rotation and next_rotation advance together
			newRotation := addRotationDuration(now, plan)
			plan.FirstRotationRFC3339 = timetypes.NewRFC3339TimeValue(newRotation)
			plan.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(newRotation)
		}

		// Update computed fields from next rotation timestamp
		nextRotationTime, diags := plan.NextRotationRFC3339.ValueRFC3339Time()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		setComputedTimestampFields(&plan, nextRotationTime)

		diags = resp.Plan.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// THIS IS THE KEY FIX: Use RequiresReplace instead of RemoveResource
		// This generates a Replace action that triggers replace_triggered_by
		resp.RequiresReplace = append(resp.RequiresReplace, path.Root("next_rotation_rfc3339"))
	}

	// Handle changes to rotation units (similar to v1)
	if state.RotationYears != plan.RotationYears ||
		state.RotationMonths != plan.RotationMonths ||
		state.RotationDays != plan.RotationDays ||
		state.RotationHours != plan.RotationHours ||
		state.RotationMinutes != plan.RotationMinutes {

		// Recalculate next_rotation based on first_rotation + new rotation units
		firstRotation, diags := state.FirstRotationRFC3339.ValueRFC3339Time()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		newNextRotation := addRotationDuration(firstRotation, plan)
		plan.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(newNextRotation)
		plan.FirstRotationRFC3339 = state.FirstRotationRFC3339

		setComputedTimestampFields(&plan, newNextRotation)

		diags = resp.Plan.Set(ctx, plan)
		resp.Diagnostics.Append(diags...)
	}
}

func (t *timeRotatingV2Resource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id := req.ID
	var state timeRotatingV2Model

	idParts := strings.Split(id, ",")

	if len(idParts) != 6 {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected FIRST_ROTATION_RFC3339,YEARS,MONTHS,DAYS,HOURS,MINUTES", id))
		return
	}

	if idParts[0] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), first rotation timestamp cannot be empty", id))
		return
	}

	if idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "" {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), at least one rotation value must be non-empty", id))
		return
	}

	// Parse first rotation timestamp
	firstRotation, err := time.Parse(time.RFC3339, idParts[0])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			"The first rotation timestamp could not be parsed as RFC3339.\n\n"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	// Parse rotation units
	rotationYears, err := rotationToInt64(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			fmt.Sprintf("Could not parse rotation_years: %s", err),
		)
		return
	}

	rotationMonths, err := rotationToInt64(idParts[2])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			fmt.Sprintf("Could not parse rotation_months: %s", err),
		)
		return
	}

	rotationDays, err := rotationToInt64(idParts[3])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			fmt.Sprintf("Could not parse rotation_days: %s", err),
		)
		return
	}

	rotationHours, err := rotationToInt64(idParts[4])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			fmt.Sprintf("Could not parse rotation_hours: %s", err),
		)
		return
	}

	rotationMinutes, err := rotationToInt64(idParts[5])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time_rotating_v2 error",
			fmt.Sprintf("Could not parse rotation_minutes: %s", err),
		)
		return
	}

	state.FirstRotationRFC3339 = timetypes.NewRFC3339TimeValue(firstRotation)
	state.RotationYears = rotationYears
	state.RotationMonths = rotationMonths
	state.RotationDays = rotationDays
	state.RotationHours = rotationHours
	state.RotationMinutes = rotationMinutes
	state.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	// Calculate next rotation
	nextRotation := addRotationDuration(firstRotation, state)
	state.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(nextRotation)

	setComputedTimestampFields(&state, nextRotation)

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeRotatingV2Model

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	now := t.clock.Now().UTC()

	var firstRotation time.Time
	var nextRotation time.Time

	// Check if user provided first_rotation_rfc3339
	if !plan.FirstRotationRFC3339.IsNull() && plan.FirstRotationRFC3339.ValueString() != "" {
		// User configured first_rotation (drift mode)
		var diags diag.Diagnostics
		firstRotation, diags = plan.FirstRotationRFC3339.ValueRFC3339Time()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		// If we're already past the first rotation (can happen on import or recreation),
		// set next_rotation to now + rotation_duration
		if now.After(firstRotation) {
			nextRotation = addRotationDuration(now, plan)
		} else {
			nextRotation = firstRotation
		}
	} else {
		// Compute first_rotation as now + rotation duration (non-drift mode)
		firstRotation = addRotationDuration(now, plan)
		plan.FirstRotationRFC3339 = timetypes.NewRFC3339TimeValue(firstRotation)
		nextRotation = firstRotation
	}

	plan.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(nextRotation)
	setComputedTimestampFields(&plan, nextRotation)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingV2Resource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Completely passive - do nothing
	// DO NOT check rotation expiry here
	// DO NOT call RemoveResource()
	// ModifyPlan handles rotation detection and triggers replacement
}

func (t *timeRotatingV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan, state timeRotatingV2Model

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

	// Handle changes to rotation unit attributes
	if state.RotationYears != plan.RotationYears ||
		state.RotationMonths != plan.RotationMonths ||
		state.RotationDays != plan.RotationDays ||
		state.RotationHours != plan.RotationHours ||
		state.RotationMinutes != plan.RotationMinutes {

		// Preserve first_rotation_rfc3339 from state
		plan.FirstRotationRFC3339 = state.FirstRotationRFC3339

		// Recalculate next_rotation = first_rotation + new rotation units
		firstRotation, diags := state.FirstRotationRFC3339.ValueRFC3339Time()
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}

		nextRotation := addRotationDuration(firstRotation, plan)
		plan.NextRotationRFC3339 = timetypes.NewRFC3339TimeValue(nextRotation)

		setComputedTimestampFields(&plan, nextRotation)
	}

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (t *timeRotatingV2Resource) Delete(ctx context.Context, request resource.DeleteRequest, response *resource.DeleteResponse) {
	// Standard cleanup - state is automatically removed
}

type timeRotatingV2Model struct {
	Day                  types.Int64       `tfsdk:"day"`
	RotationDays         types.Int64       `tfsdk:"rotation_days"`
	RotationHours        types.Int64       `tfsdk:"rotation_hours"`
	RotationMinutes      types.Int64       `tfsdk:"rotation_minutes"`
	RotationMonths       types.Int64       `tfsdk:"rotation_months"`
	RotationYears        types.Int64       `tfsdk:"rotation_years"`
	FirstRotationRFC3339 timetypes.RFC3339 `tfsdk:"first_rotation_rfc3339"`
	NextRotationRFC3339  timetypes.RFC3339 `tfsdk:"next_rotation_rfc3339"`
	Hour                 types.Int64       `tfsdk:"hour"`
	Triggers             types.Map         `tfsdk:"triggers"`
	Minute               types.Int64       `tfsdk:"minute"`
	Month                types.Int64       `tfsdk:"month"`
	Second               types.Int64       `tfsdk:"second"`
	Unix                 types.Int64       `tfsdk:"unix"`
	Year                 types.Int64       `tfsdk:"year"`
	ID                   timetypes.RFC3339 `tfsdk:"id"`
}

// addRotationDuration adds all rotation units cumulatively to baseTime.
// This fixes the v1 bug where only one rotation unit was used.
func addRotationDuration(baseTime time.Time, model timeRotatingV2Model) time.Time {
	result := baseTime

	// Add calendar units cumulatively (years, months, days)
	years := int(model.RotationYears.ValueInt64())
	months := int(model.RotationMonths.ValueInt64())
	days := int(model.RotationDays.ValueInt64())

	if years != 0 || months != 0 || days != 0 {
		result = result.AddDate(years, months, days)
	}

	// Add time duration cumulatively (hours, minutes)
	var duration time.Duration
	if !model.RotationHours.IsNull() && model.RotationHours.ValueInt64() > 0 {
		duration += time.Duration(model.RotationHours.ValueInt64()) * time.Hour
	}
	if !model.RotationMinutes.IsNull() && model.RotationMinutes.ValueInt64() > 0 {
		duration += time.Duration(model.RotationMinutes.ValueInt64()) * time.Minute
	}
	if duration > 0 {
		result = result.Add(duration)
	}

	return result
}

// setComputedTimestampFields sets all computed timestamp fields from a time value.
func setComputedTimestampFields(model *timeRotatingV2Model, timestamp time.Time) {
	model.Year = types.Int64Value(int64(timestamp.Year()))
	model.Month = types.Int64Value(int64(timestamp.Month()))
	model.Day = types.Int64Value(int64(timestamp.Day()))
	model.Hour = types.Int64Value(int64(timestamp.Hour()))
	model.Minute = types.Int64Value(int64(timestamp.Minute()))
	model.Second = types.Int64Value(int64(timestamp.Second()))
	model.Unix = types.Int64Value(timestamp.Unix())
	model.ID = timetypes.NewRFC3339TimeValue(timestamp)
}
