// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ action.Action = &TimeSleepAction{}

func NewTimeSleepAction() action.Action {
	return &TimeSleepAction{}
}

// TimeSleepAction defines the action implementation.
type TimeSleepAction struct{}

// TimeSleepActionModel describes the action data model.
type TimeSleepActionModel struct {
	Duration types.String `tfsdk:"duration"`
}

func (a *TimeSleepAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sleep"
}

func (a *TimeSleepAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Waits for the specified duration to elapse before completing.",
		MarkdownDescription: "Waits for the specified duration to elapse before completing. This action is useful for introducing delays in Terraform workflows.",
		Attributes: map[string]schema.Attribute{
			"duration": schema.StringAttribute{
				Description: "Time duration to wait. Must be a valid Go duration string (e.g., \"30s\", \"5m\", \"1h\").",
				MarkdownDescription: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to wait. Must be a number immediately followed by a unit: `ms`" +
					" (milliseconds), `s` (seconds), `m` (minutes), or `h` (hours). For example, `\"30s\"` for 30 seconds or `\"1.5m\"` for 90 seconds.",
				Required: true,
				Validators: []validator.String{
					stringvalidator.RegexMatches(
						regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds.",
					),
				},
			},
		},
	}
}

func (a *TimeSleepAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var data TimeSleepActionModel

	diags := req.Config.Get(ctx, &data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	duration, err := time.ParseDuration(data.Duration.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid Duration",
			fmt.Sprintf("The duration value could not be parsed: %s\n\nOriginal Error: %s", data.Duration.ValueString(), err),
		)
		return
	}

	select {
	case <-ctx.Done():
		resp.Diagnostics.AddError(
			"Action Cancelled",
			fmt.Sprintf("The sleep action was cancelled: %s", ctx.Err()),
		)
		return
	case <-time.After(duration):
		// Sleep completed successfully
	}
}
