package provider

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-provider-time/internal/validators/timevalidator"
	"time"
)

var _ tfsdk.ResourceType = (*timeStaticResourceType)(nil)

type timeStaticResourceType struct{}

func (t timeStaticResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"day": {
				Description: "Number day of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. " +
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
				PlanModifiers: []tfsdk.AttributePlanModifier{
					tfsdk.UseStateForUnknown(),
				},
			},
		},
	}, nil
}

func (t timeStaticResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &timeStaticResource{}, nil
}

var (
	_ tfsdk.Resource                = (*timeStaticResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeStaticResource)(nil)
)

type timeStaticResource struct {
}

func (t timeStaticResource) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

type timeStaticModelV0 struct {
	Day      types.Int64  `tfsdk:"day"`
	Hour     types.Int64  `tfsdk:"hour"`
	Triggers types.Map    `tfsdk:"triggers"`
	Minute   types.Int64  `tfsdk:"minute"`
	Month    types.Int64  `tfsdk:"month"`
	RFC3339  types.String `tfsdk:"rfc3339"`
	Second   types.Int64  `tfsdk:"second"`
	Unix     types.Int64  `tfsdk:"unix"`
	Year     types.Int64  `tfsdk:"year"`
	ID       types.String `tfsdk:"id"`
}

func (t timeStaticResource) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	var plan timeStaticModelV0

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

	formattedTimestamp := timestamp.Format(time.RFC3339)

	state := timeStaticModelV0{
		Triggers: plan.Triggers,
		Year:     types.Int64{Value: int64(timestamp.Year())},
		Month:    types.Int64{Value: int64(timestamp.Month())},
		Day:      types.Int64{Value: int64(timestamp.Day())},
		Hour:     types.Int64{Value: int64(timestamp.Hour())},
		Minute:   types.Int64{Value: int64(timestamp.Minute())},
		Second:   types.Int64{Value: int64(timestamp.Second())},
		RFC3339:  types.String{Value: formattedTimestamp},
		Unix:     types.Int64{Value: timestamp.Unix()},
		ID:       types.String{Value: formattedTimestamp},
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (t timeStaticResource) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {

}

func (t timeStaticResource) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	var data timeStaticModelV0

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (t timeStaticResource) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {

}
