package provider

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                     = (*timeOffsetResource)(nil)
	_ resource.ResourceWithImportState      = (*timeOffsetResource)(nil)
	_ resource.ResourceWithModifyPlan       = (*timeOffsetResource)(nil)
	_ resource.ResourceWithConfigValidators = (*timeOffsetResource)(nil)
)

func NewTimeOffsetResource() resource.Resource {
	return &timeOffsetResource{}
}

type timeOffsetResource struct{}

func (t timeOffsetResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_offset"
}

func (t timeOffsetResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an offset time resource, which keeps an UTC timestamp stored in the Terraform state that is" +
			" offset from a locally sourced base timestamp. This prevents perpetual differences caused " +
			"by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).",
		Attributes: map[string]schema.Attribute{
			"base_rfc3339": schema.StringAttribute{
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Optional: true,
				Computed: true,
			},
			"day": schema.Int64Attribute{
				Description: "Number day of offset timestamp.",
				Computed:    true,
			},
			"hour": schema.Int64Attribute{
				Description: "Number hour of offset timestamp.",
				Computed:    true,
			},
			"triggers": schema.MapAttribute{
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value " +
					"to be saved. See [the main provider documentation](../index.md) for more information.",
				ElementType: types.StringType,
				Optional:    true,
				PlanModifiers: []planmodifier.Map{
					mapplanmodifier.RequiresReplace(),
				},
			},
			"minute": schema.Int64Attribute{
				Description: "Number minute of offset timestamp.",
				Computed:    true,
			},
			"month": schema.Int64Attribute{
				Description: "Number month of offset timestamp.",
				Computed:    true,
			},
			"offset_days": schema.Int64Attribute{
				Description: "Number of days to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"offset_hours": schema.Int64Attribute{
				Description: " Number of hours to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"offset_minutes": schema.Int64Attribute{
				Description: "Number of minutes to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"offset_months": schema.Int64Attribute{
				Description: "Number of months to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"offset_seconds": schema.Int64Attribute{
				Description: "Number of seconds to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"offset_years": schema.Int64Attribute{
				Description: "Number of years to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Optional:    true,
			},
			"rfc3339": schema.StringAttribute{
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Computed:    true,
			},
			"second": schema.Int64Attribute{
				Description: "Number second of offset timestamp.",
				Computed:    true,
			},
			"unix": schema.Int64Attribute{
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Computed:    true,
			},
			"year": schema.Int64Attribute{
				Description: "Number year of offset timestamp.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Computed:    true,
			},
		},
	}
}

func (t timeOffsetResource) ConfigValidators(_ context.Context) []resource.ConfigValidator {
	return []resource.ConfigValidator{
		resourcevalidator.AtLeastOneOf(
			path.MatchRoot("offset_seconds"),
			path.MatchRoot("offset_minutes"),
			path.MatchRoot("offset_hours"),
			path.MatchRoot("offset_days"),
			path.MatchRoot("offset_months"),
			path.MatchRoot("offset_years"),
		),
	}
}

func (t timeOffsetResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// Plan does not need to be modified when the resource is being destroyed.
	if req.Plan.Raw.IsNull() {
		return
	}

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
	if baseRFC3339.IsUnknown() {
		resp.Diagnostics.Append(req.State.GetAttribute(ctx, path.Root("base_rfc3339"), &baseRFC3339)...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	timestamp, err := time.Parse(time.RFC3339, baseRFC3339.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Create time offset error",
			"The base_rfc3339 timestamp could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}
	setOffsetValues(&plan, timestamp)

	diags = resp.Plan.Set(ctx, plan)
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

	setOffsetValues(&importedState, timestamp)
	importedState.Triggers = types.MapValueMust(types.StringType, map[string]attr.Value{})

	diags := resp.State.Set(ctx, importedState)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan timeOffsetModelV0

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	timestamp := time.Now().UTC()

	if plan.BaseRFC3339.ValueString() != "" {
		var err error

		if timestamp, err = time.Parse(time.RFC3339, plan.BaseRFC3339.ValueString()); err != nil {
			resp.Diagnostics.AddError(
				"Create time offset error",
				"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	setOffsetValues(&plan, timestamp)
	diags = resp.State.Set(ctx, plan)
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

	timestamp, err := time.Parse(time.RFC3339, plan.BaseRFC3339.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Update time offset error",
			"The base_rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	setOffsetValues(&plan, timestamp)
	diags = resp.State.Set(ctx, plan)
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

func setOffsetValues(plan *timeOffsetModelV0, timestamp time.Time) {
	formattedTimestamp := timestamp.Format(time.RFC3339)

	var offsetTimestamp time.Time

	if plan.OffsetDays.ValueInt64() != 0 {
		offsetTimestamp = timestamp.AddDate(0, 0, int(plan.OffsetDays.ValueInt64()))
	}

	if plan.OffsetHours.ValueInt64() != 0 {
		hours := time.Duration(plan.OffsetHours.ValueInt64()) * time.Hour
		offsetTimestamp = timestamp.Add(hours)
	}

	if plan.OffsetMinutes.ValueInt64() != 0 {
		minutes := time.Duration(plan.OffsetMinutes.ValueInt64()) * time.Minute
		offsetTimestamp = timestamp.Add(minutes)
	}

	if plan.OffsetMonths.ValueInt64() != 0 {
		offsetTimestamp = timestamp.AddDate(0, int(plan.OffsetMonths.ValueInt64()), 0)
	}

	if plan.OffsetSeconds.ValueInt64() != 0 {
		seconds := time.Duration(plan.OffsetSeconds.ValueInt64()) * time.Second
		offsetTimestamp = timestamp.Add(seconds)
	}

	if plan.OffsetYears.ValueInt64() != 0 {
		offsetTimestamp = timestamp.AddDate(int(plan.OffsetYears.ValueInt64()), 0, 0)
	}

	formattedOffsetTimestamp := offsetTimestamp.Format(time.RFC3339)

	plan.BaseRFC3339 = types.StringValue(formattedTimestamp)
	plan.Year = types.Int64Value(int64(offsetTimestamp.Year()))
	plan.Month = types.Int64Value(int64(offsetTimestamp.Month()))
	plan.Day = types.Int64Value(int64(offsetTimestamp.Day()))
	plan.Hour = types.Int64Value(int64(offsetTimestamp.Hour()))
	plan.Minute = types.Int64Value(int64(offsetTimestamp.Minute()))
	plan.Second = types.Int64Value(int64(offsetTimestamp.Second()))
	plan.RFC3339 = types.StringValue(formattedOffsetTimestamp)
	plan.Unix = types.Int64Value(offsetTimestamp.Unix())
	plan.ID = types.StringValue(formattedTimestamp)
}

func offsetToInt64(offsetStr string) (types.Int64, error) {
	offset := types.Int64Null()

	if offsetStr != "" {
		offsetInt, err := strconv.ParseInt(offsetStr, 10, 64)
		if err != nil {
			return offset, fmt.Errorf("could not parse offset (%q) as int: %w", offsetStr, err)
		}

		offset = types.Int64Value(offsetInt)
	}

	return offset, nil
}
