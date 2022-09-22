package timemodifier

import (
	"context"
	"time"

	"github.com/bflad/terraform-plugin-framework-type-rfc3339/rfc3339type"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func ReplaceIfOutdated() tfsdk.AttributePlanModifier {
	return RequiresReplaceModifier{}
}

// RequiresReplaceModifier is an AttributePlanModifier that sets RequiresReplace
// on the attribute.
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

	var rotationRFC3339 rfc3339type.Value
	diags := tfsdk.ValueAs(ctx, req.AttributeState, &rotationRFC3339)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	now := time.Now().UTC()

	resp.RequiresReplace = now.After(rotationRFC3339.Time())
}
