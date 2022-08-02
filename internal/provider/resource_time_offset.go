package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"strconv"
	"strings"
	"time"
)

var _ tfsdk.ResourceType = (*timeOffsetResourceType)(nil)

type timeOffsetResourceType struct{}

func (t timeOffsetResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"base_rfc3339": {
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				//ForceNew:     true,
				//ValidateFunc: validation.IsRFC3339Time,
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
					tfsdk.RequiresReplace(),
				},
				//ForceNew: true,
				//Elem:     &schema.Schema{Type: schema.TypeString},
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
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
			},
			"offset_hours": {
				Description: " Number of hours to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
			},
			"offset_minutes": {
				Description: "Number of minutes to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
			},
			"offset_months": {
				Description: "Number of months to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
			},
			"offset_seconds": {
				Description: "Number of seconds to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
			},
			"offset_years": {
				Description: "Number of years to offset the base timestamp. At least one of the 'offset_' arguments must be configured.",
				Type:        types.Int64Type,
				Optional:    true,
				//AtLeastOneOf: []string{
				//	"offset_days",
				//	"offset_hours",
				//	"offset_minutes",
				//	"offset_months",
				//	"offset_seconds",
				//	"offset_years",
				//},
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

func (t timeOffsetResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &timeOffsetResource{}, nil
}

var (
	_ tfsdk.Resource                = (*timeOffsetResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeOffsetResource)(nil)
)

type timeOffsetResource struct {
}

func (t timeOffsetResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
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

	offsetYears := types.Int64{Null: true}
	if idParts[1] != "" {
		offset, err := strconv.ParseInt(idParts[1], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetYears.Value = offset
		offsetYears.Null = false
	}

	offsetMonths := types.Int64{Null: true}
	if idParts[2] != "" {
		offset, err := strconv.ParseInt(idParts[2], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetMonths.Value = offset
		offsetMonths.Null = false
	}

	offsetDays := types.Int64{Null: true}
	if idParts[3] != "" {
		offset, err := strconv.ParseInt(idParts[3], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetDays.Value = offset
		offsetDays.Null = false
	}

	offsetHours := types.Int64{Null: true}
	if idParts[4] != "" {
		offset, err := strconv.ParseInt(idParts[4], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetHours.Value = offset
		offsetHours.Null = false
	}

	offsetMinutes := types.Int64{Null: true}
	if idParts[5] != "" {
		offset, err := strconv.ParseInt(idParts[5], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetMinutes.Value = offset
		offsetMinutes.Null = false
	}

	offsetSeconds := types.Int64{Null: true}
	if idParts[6] != "" {
		offset, err := strconv.ParseInt(idParts[6], 10, 64)
		if err != nil {
			// return diagnostic
		}

		offsetSeconds.Value = offset
		offsetSeconds.Null = false
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

	formattedTimestamp := timestamp.Format(time.RFC3339)

	var offsetTimestamp time.Time

	if !offsetDays.Null && offsetDays.Value > 0 {
		offsetTimestamp = timestamp.AddDate(0, 0, int(offsetDays.Value))
	}

	if !offsetHours.Null && offsetHours.Value != 0 {
		hours := time.Duration(offsetHours.Value)
		offsetTimestamp = timestamp.Add(hours * time.Hour)
	}

	if !offsetMinutes.Null && offsetMinutes.Value != 0 {
		minutes := time.Duration(offsetMinutes.Value)
		offsetTimestamp = timestamp.Add(minutes * time.Minute)
	}

	if !offsetMonths.Null && offsetMonths.Value != 0 {
		offsetTimestamp = timestamp.AddDate(0, int(offsetMonths.Value), 0)
	}

	if !offsetSeconds.Null && offsetSeconds.Value != 0 {
		seconds := time.Duration(offsetSeconds.Value)
		offsetTimestamp = timestamp.Add(seconds * time.Second)
	}

	if !offsetYears.Null && offsetYears.Value != 0 {
		offsetTimestamp = timestamp.AddDate(int(offsetYears.Value), 0, 0)
	}

	formattedOffsetTimestamp := offsetTimestamp.Format(time.RFC3339)

	state := timeOffsetModelV0{
		BaseRFC3339: types.String{Value: formattedTimestamp},
		//Triggers:      plan.Triggers,
		Year:   types.Int64{Value: int64(offsetTimestamp.Year())},
		Month:  types.Int64{Value: int64(offsetTimestamp.Month())},
		Day:    types.Int64{Value: int64(offsetTimestamp.Day())},
		Hour:   types.Int64{Value: int64(offsetTimestamp.Hour())},
		Minute: types.Int64{Value: int64(offsetTimestamp.Minute())},
		Second: types.Int64{Value: int64(offsetTimestamp.Second())},
		// Need to handle instances where the ID string passed into the import function contains empty string
		// for the offset (e.g., years). If so, we need to set Null on the type, as no value has been supplied.
		OffsetYears:   offsetYears,
		OffsetMonths:  offsetMonths,
		OffsetDays:    offsetDays,
		OffsetHours:   offsetHours,
		OffsetMinutes: offsetMinutes,
		OffsetSeconds: offsetSeconds,
		RFC3339:       types.String{Value: formattedOffsetTimestamp},
		Unix:          types.Int64{Value: offsetTimestamp.Unix()},
		ID:            types.String{Value: formattedTimestamp},
	}

	state.Triggers.ElemType = types.StringType

	diags := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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

func (t timeOffsetResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
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
	formattedTimestamp := timestamp.Format(time.RFC3339)

	var offsetTimestamp time.Time

	if plan.OffsetDays.Value != 0 {
		offsetTimestamp = timestamp.AddDate(0, 0, int(plan.OffsetDays.Value))
	}

	if plan.OffsetHours.Value != 0 {
		hours := time.Duration(plan.OffsetHours.Value)
		offsetTimestamp = timestamp.Add(hours * time.Hour)
	}

	if plan.OffsetMinutes.Value != 0 {
		minutes := time.Duration(plan.OffsetMinutes.Value)
		offsetTimestamp = timestamp.Add(minutes * time.Minute)
	}

	if plan.OffsetMonths.Value != 0 {
		offsetTimestamp = timestamp.AddDate(0, int(plan.OffsetMonths.Value), 0)
	}

	if plan.OffsetSeconds.Value != 0 {
		seconds := time.Duration(plan.OffsetSeconds.Value)
		offsetTimestamp = timestamp.Add(seconds * time.Second)
	}

	if plan.OffsetYears.Value != 0 {
		offsetTimestamp = timestamp.AddDate(int(plan.OffsetYears.Value), 0, 0)
	}

	formattedOffsetTimestamp := offsetTimestamp.Format(time.RFC3339)

	state := timeOffsetModelV0{
		BaseRFC3339:   types.String{Value: formattedTimestamp},
		Triggers:      plan.Triggers,
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

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeOffsetResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {

}

func (t timeOffsetResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeOffsetResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {

}
