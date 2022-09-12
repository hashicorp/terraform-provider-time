package provider

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/int64validator"
	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-provider-time/internal/modifiers/timemodifier"
	"github.com/hashicorp/terraform-provider-time/internal/validators/timevalidator"
)

var _ tfsdk.ResourceType = (*timeRotatingResourceType)(nil)

type timeRotatingResourceType struct{}

func (t timeRotatingResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages a rotating time resource, which keeps a rotating UTC timestamp stored in the Terraform " +
			"state and proposes resource recreation when the locally sourced current time is beyond the rotation time. " +
			"This rotation only occurs when Terraform is executed, meaning there will be drift between the rotation " +
			"timestamp and actual rotation. The new rotation timestamp offset includes this drift. " +
			"This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html) " +
			"by only forcing a new value on the set cadence.",
		Attributes: map[string]tfsdk.Attribute{
			"day": {
				Description: "Number day of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"rotation_days": {
				Description: "Number of days to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_hours"),
						path.MatchRoot("rotation_minutes"),
						path.MatchRoot("rotation_months"),
						path.MatchRoot("rotation_rfc3339"),
						path.MatchRoot("rotation_years")),
					int64validator.AtLeast(1),
				},
			},
			"rotation_hours": {
				Description: "Number of hours to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_days"),
						path.MatchRoot("rotation_minutes"),
						path.MatchRoot("rotation_months"),
						path.MatchRoot("rotation_rfc3339"),
						path.MatchRoot("rotation_years")),
					int64validator.AtLeast(1),
				},
			},
			"rotation_minutes": {
				Description: "Number of minutes to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_days"),
						path.MatchRoot("rotation_hours"),
						path.MatchRoot("rotation_months"),
						path.MatchRoot("rotation_rfc3339"),
						path.MatchRoot("rotation_years")),
					int64validator.AtLeast(1),
				},
			},
			"rotation_months": {
				Description: "Number of months to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_days"),
						path.MatchRoot("rotation_hours"),
						path.MatchRoot("rotation_minutes"),
						path.MatchRoot("rotation_rfc3339"),
						path.MatchRoot("rotation_years")),
					int64validator.AtLeast(1),
				},
			},
			"rotation_rfc3339": {
				Description: "Configure the rotation timestamp with an " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format of the offset timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					timemodifier.ReplaceIfOutdated(),
				},
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_days"),
						path.MatchRoot("rotation_hours"),
						path.MatchRoot("rotation_minutes"),
						path.MatchRoot("rotation_months"),
						path.MatchRoot("rotation_years")),
					timevalidator.IsRFC3339Time(),
				},
			},
			"rotation_years": {
				Description: "Number of years to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("rotation_days"),
						path.MatchRoot("rotation_hours"),
						path.MatchRoot("rotation_minutes"),
						path.MatchRoot("rotation_months"),
						path.MatchRoot("rotation_rfc3339")),
					int64validator.AtLeast(1),
				},
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved." +
					" These conditions recreate the resource in addition to other rotation arguments. " +
					"See [the main provider documentation](../index.md) for more information.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
			},
			"minute": {
				Description: "Number minute of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"month": {
				Description: "Number month of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"rfc3339": {
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.RequiresReplace(),
				},
				Validators: []tfsdk.AttributeValidator{
					timevalidator.IsRFC3339Time(),
				},
			},
			"second": {
				Description: "Number second of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"unix": {
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"year": {
				Description: "Number year of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"id": {
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:        types.StringType,
				Computed:    true,
			},
		},
	}, nil
}

func (t timeRotatingResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &timeRotatingResource{}, nil
}

var (
	_ tfsdk.Resource                = (*timeRotatingResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeRotatingResource)(nil)
	_ tfsdk.ResourceWithModifyPlan  = (*timeRotatingResource)(nil)
)

type timeRotatingResource struct {
}

func (t timeRotatingResource) ModifyPlan(ctx context.Context, req tfsdk.ModifyResourcePlanRequest, resp *tfsdk.ModifyResourcePlanResponse) {
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

	var RFC3339, rotationRFC3339 types.String

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
	if RFC3339.Unknown {
		diags = req.State.GetAttribute(ctx, path.Root("rfc3339"), &RFC3339)
	}

	if rotationRFC3339.Unknown {
		diags = req.State.GetAttribute(ctx, path.Root("rotation_rfc3339"), &rotationRFC3339)
	}

	timestamp, err := time.Parse(time.RFC3339, RFC3339.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update time rotating error",
			"The rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	updatedPlan, err := setRotationValues(&plan, timestamp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update time rotating error",
			"The rotation_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	updatedPlan.Triggers = plan.Triggers

	diags = resp.Plan.Set(ctx, updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (t timeRotatingResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	id := req.ID
	state := timeRotatingModelV0{}

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

		state = parseTwoPartId(idParts, resp)

	} else {
		if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "") {
			resp.Diagnostics.AddError(
				"Unexpected Format of ID",
				fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES where at least one rotation value is non-empty", id))

			return
		}
		state = parseMultiplePartId(idParts, resp)

	}

	if state.ID.IsNull() {
		//If the ID of the state is not set, then something went wrong with parsing the import id
		return
	}

	state.Triggers.ElemType = types.StringType

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

type timeRotatingModelV0 struct {
	Day             types.Int64  `tfsdk:"day"`
	RotationDays    types.Int64  `tfsdk:"rotation_days"`
	RotationHours   types.Int64  `tfsdk:"rotation_hours"`
	RotationMinutes types.Int64  `tfsdk:"rotation_minutes"`
	RotationMonths  types.Int64  `tfsdk:"rotation_months"`
	RotationRFC3339 types.String `tfsdk:"rotation_rfc3339"`
	RotationYears   types.Int64  `tfsdk:"rotation_years"`
	Hour            types.Int64  `tfsdk:"hour"`
	Triggers        types.Map    `tfsdk:"triggers"`
	Minute          types.Int64  `tfsdk:"minute"`
	Month           types.Int64  `tfsdk:"month"`
	RFC3339         types.String `tfsdk:"rfc3339"`
	Second          types.Int64  `tfsdk:"second"`
	Unix            types.Int64  `tfsdk:"unix"`
	Year            types.Int64  `tfsdk:"year"`
	ID              types.String `tfsdk:"id"`
}

func (t timeRotatingResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan timeRotatingModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := time.Now().UTC()

	if plan.RFC3339.Value != "" {
		var err error

		if timestamp, err = time.Parse(time.RFC3339, plan.RFC3339.Value); err != nil {
			resp.Diagnostics.AddError(
				"Create time rotating error",
				"The rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	state, err := setRotationValues(&plan, timestamp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create time rotating error",
			"The rotation_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeRotatingResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state timeRotatingModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if !state.RotationRFC3339.IsNull() && state.RotationRFC3339.Value != "" {
		now := time.Now().UTC()
		rotationTimestamp, err := time.Parse(time.RFC3339, state.RotationRFC3339.Value)
		if err != nil {
			resp.Diagnostics.AddError(
				"Read time rotating error",
				"The rotation_rfc3339 that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

		if now.After(rotationTimestamp) {
			log.Printf("[INFO] Expiration timestamp (%s) is after current timestamp (%s), removing from state", state.RotationRFC3339.Value, now.Format(time.RFC3339))
			//TODO: Implement setting ID attribute
			//req.State.SetAttribute(ctx, path.MatchRoot(""), state)
			return
		}
	}

}

func (t timeRotatingResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var plan, state timeRotatingModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.Plan.Get(ctx, &state)
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

	timestamp, err := time.Parse(time.RFC3339, plan.ID.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Update time rotating error",
			"The ID that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	state, err = setRotationValues(&plan, timestamp)
	if err != nil {
		resp.Diagnostics.AddError(
			"Update time rotating error",
			"The rotation_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

}

func (t timeRotatingResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {

}

func setRotationValues(plan *timeRotatingModelV0, timestamp time.Time) (timeRotatingModelV0, error) {
	formattedTimestamp := timestamp.Format(time.RFC3339)

	var rotationTimestamp time.Time

	if plan.RotationDays.Value != 0 {
		rotationTimestamp = timestamp.AddDate(0, 0, int(plan.RotationDays.Value))
	}

	if plan.RotationHours.Value != 0 {
		hours := time.Duration(plan.RotationHours.Value) * time.Hour
		rotationTimestamp = timestamp.Add(hours)
	}

	if plan.RotationMinutes.Value != 0 {
		minutes := time.Duration(plan.RotationMinutes.Value) * time.Minute
		rotationTimestamp = timestamp.Add(minutes)
	}

	if plan.RotationMonths.Value != 0 {
		rotationTimestamp = timestamp.AddDate(0, int(plan.RotationMonths.Value), 0)
	}

	if plan.RotationRFC3339.Value != "" {
		var err error

		if rotationTimestamp, err = time.Parse(time.RFC3339, plan.RotationRFC3339.Value); err != nil {
			return timeRotatingModelV0{}, err
		}
	}

	if plan.RotationYears.Value != 0 {
		rotationTimestamp = timestamp.AddDate(int(plan.RotationYears.Value), 0, 0)
	}

	formattedRotationTimestamp := rotationTimestamp.Format(time.RFC3339)

	return timeRotatingModelV0{
		RotationRFC3339: types.String{Value: formattedRotationTimestamp},
		Triggers:        plan.Triggers,
		Year:            types.Int64{Value: int64(rotationTimestamp.Year())},
		Month:           types.Int64{Value: int64(rotationTimestamp.Month())},
		Day:             types.Int64{Value: int64(rotationTimestamp.Day())},
		Hour:            types.Int64{Value: int64(rotationTimestamp.Hour())},
		Minute:          types.Int64{Value: int64(rotationTimestamp.Minute())},
		Second:          types.Int64{Value: int64(rotationTimestamp.Second())},
		RotationYears:   plan.RotationYears,
		RotationMonths:  plan.RotationMonths,
		RotationDays:    plan.RotationDays,
		RotationHours:   plan.RotationHours,
		RotationMinutes: plan.RotationMinutes,
		RFC3339:         types.String{Value: formattedTimestamp},
		Unix:            types.Int64{Value: rotationTimestamp.Unix()},
		ID:              types.String{Value: formattedTimestamp},
	}, nil

}

func parseTwoPartId(idParts []string, resp *tfsdk.ImportResourceStateResponse) timeRotatingModelV0 {

	baseRfc3339 := idParts[0]
	rotationRfc3339 := idParts[1]

	timestamp, err := time.Parse(time.RFC3339, baseRfc3339)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import Timestamp Error",
			"The baseRfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	rotationTimestamp, err := time.Parse(time.RFC3339, rotationRfc3339)
	if err != nil {
		resp.Diagnostics.AddError(
			"Import Timestamp Error",
			"The rotationRfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	formattedTimestamp := timestamp.Format(time.RFC3339)

	return timeRotatingModelV0{
		Year:            types.Int64{Value: int64(rotationTimestamp.Year())},
		Month:           types.Int64{Value: int64(rotationTimestamp.Month())},
		Day:             types.Int64{Value: int64(rotationTimestamp.Day())},
		Hour:            types.Int64{Value: int64(rotationTimestamp.Hour())},
		Minute:          types.Int64{Value: int64(rotationTimestamp.Minute())},
		Second:          types.Int64{Value: int64(rotationTimestamp.Second())},
		RotationRFC3339: types.String{Value: rotationRfc3339},
		RotationYears:   types.Int64{Null: true},
		RotationMonths:  types.Int64{Null: true},
		RotationDays:    types.Int64{Null: true},
		RotationHours:   types.Int64{Null: true},
		RotationMinutes: types.Int64{Null: true},
		RFC3339:         types.String{Value: formattedTimestamp},
		Unix:            types.Int64{Value: rotationTimestamp.Unix()},
		ID:              types.String{Value: formattedTimestamp},
	}
}

func parseMultiplePartId(idParts []string, resp *tfsdk.ImportResourceStateResponse) timeRotatingModelV0 {
	baseRfc3339 := idParts[0]

	rotationYears, err := rotationToInt64(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time rotation error",
			"The rotation_years parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	rotationMonths, err := rotationToInt64(idParts[2])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time rotation error",
			"The rotation_months parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	rotationDays, err := rotationToInt64(idParts[3])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time rotation error",
			"The rotation_days parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	rotationHours, err := rotationToInt64(idParts[4])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time rotation error",
			"The rotation_hours parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	rotationMinutes, err := rotationToInt64(idParts[5])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time rotation error",
			"The rotation_minutes parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	timestamp, err := time.Parse(time.RFC3339, baseRfc3339)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create time rotation error",
			"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return timeRotatingModelV0{}
	}

	formattedTimestamp := timestamp.Format(time.RFC3339)

	var rotationTimestamp time.Time

	if !rotationDays.Null && rotationDays.Value > 0 {
		rotationTimestamp = timestamp.AddDate(0, 0, int(rotationDays.Value))
	}

	if !rotationHours.Null && rotationHours.Value > 0 {
		hours := time.Duration(rotationHours.Value) * time.Hour
		rotationTimestamp = timestamp.Add(hours)
	}

	if !rotationMinutes.Null && rotationMinutes.Value > 0 {
		minutes := time.Duration(rotationMinutes.Value) * time.Minute
		rotationTimestamp = timestamp.Add(minutes)
	}

	if !rotationMonths.Null && rotationMonths.Value > 0 {
		rotationTimestamp = timestamp.AddDate(0, int(rotationMonths.Value), 0)
	}

	if !rotationYears.Null && rotationYears.Value > 0 {
		rotationTimestamp = timestamp.AddDate(int(rotationYears.Value), 0, 0)
	}

	formattedRotationTimestamp := rotationTimestamp.Format(time.RFC3339)

	state := timeRotatingModelV0{
		Year:            types.Int64{Value: int64(rotationTimestamp.Year())},
		Month:           types.Int64{Value: int64(rotationTimestamp.Month())},
		Day:             types.Int64{Value: int64(rotationTimestamp.Day())},
		Hour:            types.Int64{Value: int64(rotationTimestamp.Hour())},
		Minute:          types.Int64{Value: int64(rotationTimestamp.Minute())},
		Second:          types.Int64{Value: int64(rotationTimestamp.Second())},
		RotationRFC3339: types.String{Value: formattedRotationTimestamp},
		RotationYears:   rotationYears,
		RotationMonths:  rotationMonths,
		RotationDays:    rotationDays,
		RotationHours:   rotationHours,
		RotationMinutes: rotationMinutes,
		RFC3339:         types.String{Value: formattedTimestamp},
		Unix:            types.Int64{Value: rotationTimestamp.Unix()},
		ID:              types.String{Value: baseRfc3339},
	}

	return state
}

func rotationToInt64(rotationStr string) (types.Int64, error) {
	rotation := types.Int64{Null: true}

	if rotationStr != "" {
		offsetInt, err := strconv.ParseInt(rotationStr, 10, 64)
		if err != nil {
			return rotation, fmt.Errorf("could not parse rotation (%q) as int: %w", rotationStr, err)
		}

		rotation.Value = offsetInt
		rotation.Null = false
	}

	return rotation, nil
}
