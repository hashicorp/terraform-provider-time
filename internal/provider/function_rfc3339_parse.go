package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var rfc3339ReturnAttrTypes = map[string]attr.Type{
	"year":         types.Int64Type,
	"year_day":     types.Int64Type,
	"day":          types.Int64Type,
	"month":        types.Int64Type,
	"month_name":   types.StringType,
	"weekday":      types.Int64Type,
	"weekday_name": types.StringType,
	"hour":         types.Int64Type,
	"minute":       types.Int64Type,
	"second":       types.Int64Type,
	"unix":         types.Int64Type,

	// TODO: Zone name may be tricky to accurately determine, might need more research into this
	// https://stackoverflow.com/a/30741518
	"zone_name": types.StringType,

	// TODO: This is representing as # seconds east of UTC, is this too confusing for practitioners
	"zone_offset": types.Int64Type,
	"iso_year":    types.Int64Type,
	"iso_week":    types.Int64Type,
}

var _ function.Function = &RFC3339ParseFunction{}

type RFC3339ParseFunction struct{}

func NewRFC3339ParseFunction() function.Function {
	return &RFC3339ParseFunction{}
}

func (f *RFC3339ParseFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "rfc3339_parse"
}

func (f *RFC3339ParseFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary: "Parse an RFC3339 timestamp string",
		// TODO: better wording
		Description: "Given an RFC3339 timestamp string, will parse and return the object representation of that timestamp.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "timestamp",
				Description: "RFC3339 timestamp string to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: rfc3339ReturnAttrTypes,
		},
	}
}

func (f *RFC3339ParseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var timestamp string

	resp.Diagnostics.Append(req.Arguments.Get(ctx, &timestamp)...)
	if resp.Diagnostics.HasError() {
		return
	}

	rfc3339, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		// TODO: we probably shouldn't use error messages from time.Time for practitioners because they are hecka confusing :).
		// Ex: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006".
		resp.Diagnostics.AddArgumentError(0, "Error parsing RFC3339 timestamp", err.Error())
		return
	}

	zoneName, zoneOffset := rfc3339.Zone()
	isoYear, isoWeek := rfc3339.ISOWeek()

	rfc3339Obj, diags := types.ObjectValue(
		rfc3339ReturnAttrTypes,
		map[string]attr.Value{
			"year":         types.Int64Value(int64(rfc3339.Year())),
			"year_day":     types.Int64Value(int64(rfc3339.YearDay())),
			"day":          types.Int64Value(int64(rfc3339.Day())),
			"month":        types.Int64Value(int64(rfc3339.Month())),
			"month_name":   types.StringValue(rfc3339.Month().String()),
			"weekday":      types.Int64Value(int64(rfc3339.Weekday())),
			"weekday_name": types.StringValue(rfc3339.Weekday().String()),
			"hour":         types.Int64Value(int64(rfc3339.Hour())),
			"minute":       types.Int64Value(int64(rfc3339.Minute())),
			"second":       types.Int64Value(int64(rfc3339.Second())),
			"unix":         types.Int64Value(rfc3339.Unix()),
			"zone_name":    types.StringValue(zoneName),
			"zone_offset":  types.Int64Value(int64(zoneOffset)),
			"iso_year":     types.Int64Value(int64(isoYear)),
			"iso_week":     types.Int64Value(int64(isoWeek)),
		},
	)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.Result.Set(ctx, &rfc3339Obj)...)
}
