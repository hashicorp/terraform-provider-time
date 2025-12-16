// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var unixReturnAttrTypes = map[string]attr.Type{
	"year":         types.Int64Type,
	"year_day":     types.Int64Type,
	"day":          types.Int64Type,
	"month":        types.Int64Type,
	"month_name":   types.StringType,
	"weekday":      types.Int64Type,
	"weekday_name": types.StringType,
	"hour":         types.Int64Type,
	"minute":       types.Int64Type,
	"rfc3339":      types.StringType,
	"second":       types.Int64Type,
	"iso_year":     types.Int64Type,
	"iso_week":     types.Int64Type,
}

var _ function.Function = &UnixTimestampParseFunction{}

type UnixTimestampParseFunction struct{}

func (f *UnixTimestampParseFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "unix_timestamp_parse"
}

func (f *UnixTimestampParseFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:             "Parse a unix timestamp integer into an object",
		MarkdownDescription: "Given a unix timestamp integer, will parse and return an object representation of that date and time. A unix timestamp is the number of seconds elapsed since January 1, 1970 UTC.",

		Parameters: []function.Parameter{
			function.Int64Parameter{
				Name:                "unix_timestamp",
				MarkdownDescription: "Unix Timestamp integer to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: unixReturnAttrTypes,
		},
	}
}

func (f *UnixTimestampParseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var timestamp int64

	resp.Error = req.Arguments.Get(ctx, &timestamp)
	if resp.Error != nil {
		return
	}

	unixTime := time.Unix(timestamp, 0).UTC()

	isoYear, isoWeek := unixTime.ISOWeek()

	unixTimestampObject, diags := types.ObjectValue(
		unixReturnAttrTypes,
		map[string]attr.Value{
			"year":         types.Int64Value(int64(unixTime.Year())),
			"year_day":     types.Int64Value(int64(unixTime.YearDay())),
			"day":          types.Int64Value(int64(unixTime.Day())),
			"month":        types.Int64Value(int64(unixTime.Month())),
			"month_name":   types.StringValue(unixTime.Month().String()),
			"weekday":      types.Int64Value(int64(unixTime.Weekday())),
			"weekday_name": types.StringValue(unixTime.Weekday().String()),
			"hour":         types.Int64Value(int64(unixTime.Hour())),
			"minute":       types.Int64Value(int64(unixTime.Minute())),
			"rfc3339":      types.StringValue(unixTime.Format(time.RFC3339)),
			"second":       types.Int64Value(int64(unixTime.Second())),
			"iso_year":     types.Int64Value(int64(isoYear)),
			"iso_week":     types.Int64Value(int64(isoWeek)),
		},
	)

	resp.Error = function.FuncErrorFromDiags(ctx, diags)
	if resp.Error != nil {
		return
	}

	resp.Error = resp.Result.Set(ctx, &unixTimestampObject)

}

func NewUnixTimestampParseFunction() function.Function {
	return &UnixTimestampParseFunction{}
}
