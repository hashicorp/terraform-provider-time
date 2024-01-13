// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
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
	"iso_year":     types.Int64Type,
	"iso_week":     types.Int64Type,
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
		Summary:     "Parse an RFC3339 timestamp string",
		Description: "Given an RFC3339 timestamp string, will parse and return an object representation of that date and time.",

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
		// Intentionally not returning the Go parse error to practitioners
		tflog.Error(ctx, fmt.Sprintf("failed to parse RFC3339 timestamp, underlying time.Time error: %s", err.Error()))

		resp.Diagnostics.AddArgumentError(0, "Error parsing RFC3339 timestamp", fmt.Sprintf("%q is not a valid RFC3339 timestamp", timestamp))
		return
	}

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
