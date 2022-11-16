package timevalidator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var _ tfsdk.AttributeValidator = isRFC3339TimeValidator{}

// IsRFC3339Time validates if the provided value is of type string and a valid RFC3339Time.
type isRFC3339TimeValidator struct {
}

func (validator isRFC3339TimeValidator) Description(ctx context.Context) string {
	return "value must be a string in RFC3339 format"
}

func (validator isRFC3339TimeValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator isRFC3339TimeValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	// Only validate the attribute configuration value if it is known.
	if req.AttributeConfig.IsNull() || req.AttributeConfig.IsUnknown() {
		return
	}

	t := req.AttributeConfig.Type(ctx)
	if !t.Equal(types.StringType) {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			req.AttributePath,
			"Expected value of type string",
			t.String(),
		))
		return
	}

	s := req.AttributeConfig.(types.String)

	if _, err := time.Parse(time.RFC3339, s.ValueString()); err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			req.AttributePath,
			validator.MarkdownDescription(ctx),
			s.ValueString(),
		))
		return
	}
}

// IsRFC3339Time returns an AttributeValidator which ensures that any configured
// attribute value:
//
//   - Is a String.
//   - Is in RFC3339 Format.
//
// Null (unconfigured) and unknown (known after apply) values are skipped.
func IsRFC3339Time() tfsdk.AttributeValidator {
	return isRFC3339TimeValidator{}
}
