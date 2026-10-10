// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

const (
	pkTimeSleepCloseDuration = "time_sleep_close_duration"
)

var (
	_ resource.Resource                = (*timeSleepResource)(nil)
	_ resource.ResourceWithImportState = (*timeSleepResource)(nil)
	_ resource.ResourceWithConfigure   = (*timeSleepResource)(nil)
)

var _ ephemeral.EphemeralResourceWithConfigure = (*timeSleepEphemeralResource)(nil)
var _ ephemeral.EphemeralResourceWithClose = (*timeSleepEphemeralResource)(nil)

func NewTimeSleepEphemeralResource() ephemeral.EphemeralResource {
	return &timeSleepEphemeralResource{}
}

type timeSleepEphemeralResource struct {
	clock clock.Clock
}

func (t *timeSleepEphemeralResource) Configure(ctx context.Context, req ephemeral.ConfigureRequest, resp *ephemeral.ConfigureResponse) {
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

func (t *timeSleepEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sleep"
}

func (t *timeSleepEphemeralResource) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an ephemeral resource that delays creation and/or destruction, typically for further resources. " +
			"This prevents cross-platform compatibility and destroy-time issues with using " +
			"the [`local-exec` provisioner](https://www.terraform.io/docs/provisioners/local-exec.html).",
		Attributes: map[string]schema.Attribute{
			"open_duration": schema.StringAttribute{
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource opening. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("close_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"close_duration": schema.StringAttribute{
				Description: "[Time duration](https://golang.org/pkg/time/#ParseDuration) to delay resource closing. " +
					"For example, `30s` for 30 seconds or `5m` for 5 minutes. Updating this value by itself will not trigger a delay. " +
					"This value or any updates to it must be successfully applied into the Terraform state before destroying this resource to take effect.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.AtLeastOneOf(path.MatchRoot("open_duration")),
					stringvalidator.RegexMatches(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`),
						"must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
				},
			},
			"outputs": schema.MapAttribute{
				Description: "(Optional) Arbitrary map of values that is used to be referenced by another resource/provider to introduce a dependency. " +
					"This is only useful for a provider config to indirectly reference an attribute from another ephemeral resource, via this resource to apply a delay.",
				ElementType: types.StringType,
				Optional:    true,
			},
		},
	}
}

func (t *timeSleepEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var config timeSleepEphemeralModelV0
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if config.OpenDuration.ValueString() != "" {
		duration, err := time.ParseDuration(config.OpenDuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Open time sleep error",
				"The open_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Open time sleep error",
				fmt.Sprintf("Original Error: %s", ctx.Err()),
			)
			return
		case <-time.After(duration):
		}
	}

	if !config.CloseDuration.IsNull() {
		b, err := json.Marshal(config.CloseDuration.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Marshal `close_duration` to JSON error",
				fmt.Sprintf("Original Error: %s", ctx.Err()),
			)
			return
		}
		diags = resp.Private.SetKey(ctx, pkTimeSleepCloseDuration, b)
		resp.Diagnostics.Append(diags...)
	}

	diags = resp.Result.Set(ctx, config)
	resp.Diagnostics.Append(diags...)
}

func (t *timeSleepEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	raw, diags := req.Private.GetKey(ctx, pkTimeSleepCloseDuration)
	resp.Diagnostics.Append(diags...)
	if diags.HasError() {
		return
	}

	if raw != nil {
		var rawStr string
		if err := json.Unmarshal(raw, &rawStr); err != nil {
			resp.Diagnostics.AddError(
				"Unmarshal `close_duration` from JSON error",
				fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
		duration, err := time.ParseDuration(rawStr)
		if err != nil {
			resp.Diagnostics.AddError(
				"Close time sleep error",
				"The close_duration cannot be parsed\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}

		select {
		case <-ctx.Done():
			resp.Diagnostics.AddError(
				"Close time sleep error",
				fmt.Sprintf("Original Error: %s", ctx.Err()),
			)
			return
		case <-time.After(duration):
		}
	}
}

type timeSleepEphemeralModelV0 struct {
	OpenDuration  types.String `tfsdk:"open_duration"`
	CloseDuration types.String `tfsdk:"close_duration"`
	Outputs       types.Map    `tfsdk:"outputs"`
}
