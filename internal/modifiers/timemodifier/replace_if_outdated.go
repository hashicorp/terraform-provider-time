package timemodifier

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func ReplaceIfOutdated() tfsdk.AttributePlanModifier {
	return RequiresReplaceModifier{}
}

// RequiresReplaceModifier is an AttributePlanModifier that sets RequiresReplace
// on the attribute if the current time is past the stored timestamp.
//
// This custom modifier is necessary because the resource.RequiresReplaceIf
// function uses special logic for Computed attributes which is not applicable
// this use case.
type RequiresReplaceModifier struct{}

func (r RequiresReplaceModifier) Description(ctx context.Context) string {
	return "value must be a string in RFC3339 format"
}

func (r RequiresReplaceModifier) MarkdownDescription(ctx context.Context) string {
	return r.Description(ctx)
}

func (r RequiresReplaceModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributeConfig == nil || req.AttributePlan == nil || req.AttributeState == nil {
		// shouldn't happen, but let's not panic if it does
		return
	}

	if req.State.Raw.IsNull() {
		// if we're creating the resource, no need to delete and
		// recreate it
		return
	}

	if req.Plan.Raw.IsNull() {
		// if we're deleting the resource, no need to delete and
		// recreate it
		return
	}

	var rotationRFC3339 types.String
	diags := tfsdk.ValueAs(ctx, req.AttributeState, &rotationRFC3339)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	rotationTimestamp, err := time.Parse(time.RFC3339, rotationRFC3339.Value)
	if err != nil {
		resp.Diagnostics.AddError(
			"replaceIfOutdated plan modifier error",
			"The rotation rfc3339 timestamp that was supplied could not be parsed as RFC3339.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	now := time.Now().UTC()

	resp.RequiresReplace = now.After(rotationTimestamp)
}
