package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/schemavalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = (*timeOffsetResource)(nil)
	_ resource.ResourceWithImportState = (*timeOffsetResource)(nil)
	_ resource.ResourceWithModifyPlan  = (*timeOffsetResource)(nil)
)

func NewTimeOffsetResource() resource.Resource {
	return &timeOffsetResource{}
}

type timeOffsetResource struct{}

func (t timeOffsetResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Plan only needs modifying if the resource already exists as the purpose of
	// the plan modifier is to show updated attribute values on CLI.
	if req.State.Raw.IsNull() {
		return
	}

	var state, plan timeOffsetModelV0

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

	if state.OffsetYears == plan.OffsetYears &&
		state.OffsetMonths == plan.OffsetMonths &&
		state.OffsetDays == plan.OffsetDays &&
		state.OffsetHours == plan.OffsetHours &&
		state.OffsetMinutes == plan.OffsetMinutes &&
		state.OffsetSeconds == plan.OffsetSeconds {
		return
	}

	var baseRFC3339 types.String

	diags = req.Plan.GetAttribute(ctx, path.Root("base_rfc3339"), &baseRFC3339)

	resp.Diagnostics = append(resp.Diagnostics, diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// baseRFC3339 could be unknown if there is no value set in the config as the attribute is
	// optional and computed. If base_rfc3339 is not set in config then the previous value from
	// state is used and propagated to the update function.
	if baseRFC3339.Unknown {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("base_rfc3339"), &baseRFC3339)...)
	}

	timestamp, err := time.Parse(time.RFC3339, baseRFC3339.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create time offset error",
			"The base_rfc3339 timestamp could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}
	updatedPlan := setOffsetValues(&plan, timestamp)
	updatedPlan.Triggers = plan.Triggers

	diags = resp.Plan.Set(ctx, updatedPlan)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	var importedState timeOffsetModelV0
	var err error

	id := req.ID

	idParts := strings.Split(id, ",")

	if len(idParts) != 7 {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES,SECONDS", id))

		return
	}

	if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "" && idParts[6] == "") {
		resp.Diagnostics.AddError(
			"Unexpected Format of ID",
			fmt.Sprintf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES,SECONDS where at least one offset value is non-empty", id))

		return
	}

	baseRfc3339 := idParts[0]

	importedState.OffsetYears, err = offsetToInt64(idParts[1])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_years parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	importedState.OffsetMonths, err = offsetToInt64(idParts[2])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_months parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	importedState.OffsetDays, err = offsetToInt64(idParts[3])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_days parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	importedState.OffsetHours, err = offsetToInt64(idParts[4])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_hours parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	importedState.OffsetMinutes, err = offsetToInt64(idParts[5])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_minutes parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	importedState.OffsetSeconds, err = offsetToInt64(idParts[6])
	if err != nil {
		resp.Diagnostics.AddError(
			"Import time offset error",
			"The offset_seconds parameter that was supplied could not be parsed as Int64.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	timestamp, err := time.Parse(time.RFC3339, baseRfc3339)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create time offset error",
			"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	state := setOffsetValues(&importedState, timestamp)
	state.Triggers.ElemType = types.StringType

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_offset"
}

func (t timeOffsetResource) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Manages an offset time resource, which keeps an UTC timestamp stored in the Terraform state that is" +
			" offset from a locally sourced base timestamp. This prevents perpetual differences caused " +
			"by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).",
		Attributes: map[string]tfsdk.Attribute{
			"base_rfc3339": {
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
			},
			"day": {
				Description: "Number day of offset timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"hour": {
				Description: "Number hour of offset timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value " +
					"to be saved. See [the main provider documentation](../index.md) for more information.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				PlanModifiers: []tfsdk.AttributePlanModifier{
					resource.RequiresReplace(),
				},
			},
			"minute": {
				Description: "Number minute of offset timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"month": {
				Description: "Number month of offset timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"offset_days": {
				Description: "Number of days to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_years"),
						path.MatchRoot("offset_months"),
						path.MatchRoot("offset_hours"),
						path.MatchRoot("offset_minutes"),
						path.MatchRoot("offset_seconds")),
				},
			},
			"offset_hours": {
				Description: " Number of hours to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_years"),
						path.MatchRoot("offset_months"),
						path.MatchRoot("offset_days"),
						path.MatchRoot("offset_minutes"),
						path.MatchRoot("offset_seconds")),
				},
			},
			"offset_minutes": {
				Description: "Number of minutes to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_years"),
						path.MatchRoot("offset_months"),
						path.MatchRoot("offset_days"),
						path.MatchRoot("offset_hours"),
						path.MatchRoot("offset_seconds")),
				},
			},
			"offset_months": {
				Description: "Number of months to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_years"),
						path.MatchRoot("offset_days"),
						path.MatchRoot("offset_hours"),
						path.MatchRoot("offset_minutes"),
						path.MatchRoot("offset_seconds")),
				},
			},
			"offset_seconds": {
				Description: "Number of seconds to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_years"),
						path.MatchRoot("offset_months"),
						path.MatchRoot("offset_days"),
						path.MatchRoot("offset_hours"),
						path.MatchRoot("offset_minutes")),
				},
			},
			"offset_years": {
				Description: "Number of years to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				Validators: []tfsdk.AttributeValidator{
					schemavalidator.AtLeastOneOf(path.MatchRoot("offset_months"),
						path.MatchRoot("offset_days"),
						path.MatchRoot("offset_hours"),
						path.MatchRoot("offset_minutes"),
						path.MatchRoot("offset_seconds")),
				},
			},
			"rfc3339": {
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:        types.StringType,
				Computed:    true,
			},
			"second": {
				Description: "Number second of offset timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"unix": {
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"year": {
				Description: "Number year of offset timestamp.",
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

func (t timeOffsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeOffsetModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := time.Now().UTC()

	if plan.BaseRFC3339.Value != "" {
		var err error

		if timestamp, err = time.Parse(time.RFC3339, plan.BaseRFC3339.Value); err != nil {
			resp.Diagnostics.AddError(
				"Create time offset error",
				"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	state := setOffsetValues(&plan, timestamp)
	state.Triggers = plan.Triggers
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

}

func (t timeOffsetResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan timeOffsetModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp, err := time.Parse(time.RFC3339, plan.BaseRFC3339.Value)

	if err != nil {
		resp.Diagnostics.AddError(
			"Update time offset error",
			"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	state := setOffsetValues(&plan, timestamp)
	state.Triggers = plan.Triggers
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

}

type timeOffsetModelV0 struct {
	BaseRFC3339   types.String `tfsdk:"base_rfc3339"`
	Triggers      types.Map    `tfsdk:"triggers"`
	Year          types.Int64  `tfsdk:"year"`
	Month         types.Int64  `tfsdk:"month"`
	Day           types.Int64  `tfsdk:"day"`
	Hour          types.Int64  `tfsdk:"hour"`
	Minute        types.Int64  `tfsdk:"minute"`
	Second        types.Int64  `tfsdk:"second"`
	OffsetYears   types.Int64  `tfsdk:"offset_years"`
	OffsetMonths  types.Int64  `tfsdk:"offset_months"`
	OffsetDays    types.Int64  `tfsdk:"offset_days"`
	OffsetHours   types.Int64  `tfsdk:"offset_hours"`
	OffsetMinutes types.Int64  `tfsdk:"offset_minutes"`
	OffsetSeconds types.Int64  `tfsdk:"offset_seconds"`
	RFC3339       types.String `tfsdk:"rfc3339"`
	Unix          types.Int64  `tfsdk:"unix"`
	ID            types.String `tfsdk:"id"`
}

func setOffsetValues(plan *timeOffsetModelV0, timestamp time.Time) timeOffsetModelV0 {
	formattedTimestamp := timestamp.Format(time.RFC3339)

	var offsetTimestamp time.Time

	if plan.OffsetDays.Value != 0 {
		offsetTimestamp = timestamp.AddDate(0, 0, int(plan.OffsetDays.Value))
	}

	if plan.OffsetHours.Value != 0 {
		hours := time.Duration(plan.OffsetHours.Value) * time.Hour
		offsetTimestamp = timestamp.Add(hours)
	}

	if plan.OffsetMinutes.Value != 0 {
		minutes := time.Duration(plan.OffsetMinutes.Value) * time.Minute
		offsetTimestamp = timestamp.Add(minutes)
	}

	if plan.OffsetMonths.Value != 0 {
		offsetTimestamp = timestamp.AddDate(0, int(plan.OffsetMonths.Value), 0)
	}

	if plan.OffsetSeconds.Value != 0 {
		seconds := time.Duration(plan.OffsetSeconds.Value) * time.Second
		offsetTimestamp = timestamp.Add(seconds)
	}

	if plan.OffsetYears.Value != 0 {
		offsetTimestamp = timestamp.AddDate(int(plan.OffsetYears.Value), 0, 0)
	}

	formattedOffsetTimestamp := offsetTimestamp.Format(time.RFC3339)

	return timeOffsetModelV0{
		BaseRFC3339:   types.String{Value: formattedTimestamp},
		Year:          types.Int64{Value: int64(offsetTimestamp.Year())},
		Month:         types.Int64{Value: int64(offsetTimestamp.Month())},
		Day:           types.Int64{Value: int64(offsetTimestamp.Day())},
		Hour:          types.Int64{Value: int64(offsetTimestamp.Hour())},
		Minute:        types.Int64{Value: int64(offsetTimestamp.Minute())},
		Second:        types.Int64{Value: int64(offsetTimestamp.Second())},
		OffsetYears:   plan.OffsetYears,
		OffsetMonths:  plan.OffsetMonths,
		OffsetDays:    plan.OffsetDays,
		OffsetHours:   plan.OffsetHours,
		OffsetMinutes: plan.OffsetMinutes,
		OffsetSeconds: plan.OffsetSeconds,
		RFC3339:       types.String{Value: formattedOffsetTimestamp},
		Unix:          types.Int64{Value: offsetTimestamp.Unix()},
		ID:            types.String{Value: formattedTimestamp},
	}

}

func offsetToInt64(offsetStr string) (types.Int64, error) {
	offset := types.Int64{Null: true}

	if offsetStr != "" {
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return offset, fmt.Errorf("could not parse offset (%q) as int: %w", offsetStr, err)
		}

		offset.Value = offsetInt
		offset.Null = false
	}

	return offset, nil
}
