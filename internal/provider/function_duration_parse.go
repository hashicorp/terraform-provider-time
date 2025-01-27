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

var durationParseReturnAttrTypes = map[string]attr.Type{
	"hours":        types.Float64Type,
	"minutes":      types.Float64Type,
	"seconds":      types.Float64Type,
	"milliseconds": types.Int64Type,
	"microseconds": types.Int64Type,
	"nanoseconds":  types.Int64Type,
}

var _ function.Function = &DurationParseFunction{}

type DurationParseFunction struct{}

func NewDurationParseFunction() function.Function {
	return &DurationParseFunction{}
}

func (f *DurationParseFunction) Metadata(ctx context.Context, req function.MetadataRequest, resp *function.MetadataResponse) {
	resp.Name = "duration_parse"
}

func (f *DurationParseFunction) Definition(ctx context.Context, req function.DefinitionRequest, resp *function.DefinitionResponse) {
	resp.Definition = function.Definition{
		Summary:     "Parse a [Go duration string](https://pkg.go.dev/time#ParseDuration) into an object",
		Description: "Given a [Go duration string](https://pkg.go.dev/time#ParseDuration), will parse and return an object representation of that duration.",

		Parameters: []function.Parameter{
			function.StringParameter{
				Name:        "duration",
				Description: "Go time package duration string to parse",
			},
		},
		Return: function.ObjectReturn{
			AttributeTypes: durationParseReturnAttrTypes,
		},
	}
}

func (f *DurationParseFunction) Run(ctx context.Context, req function.RunRequest, resp *function.RunResponse) {
	var input string

	resp.Error = req.Arguments.Get(ctx, &input)
	if resp.Error != nil {
		return
	}

	duration, err := time.DurationParse(input)
	if err != nil {
		// Intentionally not including the Go parse error in the return diagnostic, as the message is based on a Go-specific
		// reference time that may be unfamiliar to practitioners
		tflog.Error(ctx, fmt.Sprintf("failed to parse duration string, underlying time.Duration error: %s", err.Error()))

		resp.Error = function.NewArgumentFuncError(0, fmt.Sprintf("Error parsing duration string: %q is not a valid duration string", input))
		return
	}

	durationObj, diags := types.ObjectValue(
		durationParseReturnAttrTypes,
		map[string]attr.Value{
			"hours":        types.Float64Value(duration.Hours()),
			"minutes":      types.Float64Value(duration.Minutes()),
			"seconds":      types.Float64Value(duration.Seconds()),
			"milliseconds": types.Int64Value(duration.Milliseconds()),
			"microseconds": types.Int64Value(duration.Microseconds()),
			"nanoseconds":  types.Int64Value(duration.Nanoseconds()),
		},
	)

	resp.Error = function.FuncErrorFromDiags(ctx, diags)
	if resp.Error != nil {
		return
	}

	resp.Error = resp.Result.Set(ctx, &durationObj)
}
