package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.ResourceType = (*timeRotatingResourceType)(nil)

type timeRotatingResourceType struct{}

func (t timeRotatingResourceType) GetSchema(ctx context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{
			"day": {
				Description: "Number day of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"rotation_days": {
				Description: "Number of days to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_hours": {
				Description: "Number of hours to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_minutes": {
				Description: "Number of minutes to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_months": {
				Description: "Number of months to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_rfc3339": {
				Description: "Configure the rotation timestamp with an " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format of the offset timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.StringType,
				Optional: true,
				Computed: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IsRFC3339Time,
			},
			"rotation_years": {
				Description: "Number of years to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     types.Int64Type,
				Optional: true,
				//AtLeastOneOf: []string{
				//	"rotation_days",
				//	"rotation_hours",
				//	"rotation_minutes",
				//	"rotation_months",
				//	"rotation_rfc3339",
				//	"rotation_years",
				//},
				//ValidateFunc: validation.IntAtLeast(1),
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        types.Int64Type,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved." +
					" These conditions recreate the resource in addition to other rotation arguments. " +
					"See [the main provider documentation](../index.md) for more information.",
				Type: types.MapType{
					ElemType: types.StringType,
				},
				Optional: true,
				//ForceNew: true,
				//Elem:     &schema.Schema{Type: schema.TypeString},
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
				//ForceNew:     true,
				//ValidateFunc: validation.IsRFC3339Time,
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
				Type:        types.Int64Type,
				Computed:    true,
			},
		},
	}, nil
}

func (t timeRotatingResourceType) NewResource(ctx context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	return &timeRotatingResource{}, nil
}

var (
	_ tfsdk.Resource                = (*timeRotatingResource)(nil)
	_ tfsdk.ResourceWithImportState = (*timeRotatingResource)(nil)
)

type timeRotatingResource struct {
}

func (t timeRotatingResource) ImportState(ctx context.Context, request tfsdk.ImportResourceStateRequest, response *tfsdk.ImportResourceStateResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeRotatingResource) Create(ctx context.Context, request tfsdk.CreateResourceRequest, response *tfsdk.CreateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeRotatingResource) Read(ctx context.Context, request tfsdk.ReadResourceRequest, response *tfsdk.ReadResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeRotatingResource) Update(ctx context.Context, request tfsdk.UpdateResourceRequest, response *tfsdk.UpdateResourceResponse) {
	//TODO implement me
	panic("implement me")
}

func (t timeRotatingResource) Delete(ctx context.Context, request tfsdk.DeleteResourceRequest, response *tfsdk.DeleteResourceResponse) {
	//TODO implement me
	panic("implement me")
}
