package timevalidator

import (
	"context"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-validators/helpers/validatordiag"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

var _ validator.String = isRFC3339TimeValidator{}

// IsRFC3339Time validates if the provided value is of type string and a valid RFC3339Time.
type isRFC3339TimeValidator struct {
}

func (validator isRFC3339TimeValidator) Description(ctx context.Context) string {
	return "value must be a string in RFC3339 format"
}

func (validator isRFC3339TimeValidator) MarkdownDescription(ctx context.Context) string {
	return validator.Description(ctx)
}

func (validator isRFC3339TimeValidator) ValidateString(ctx context.Context, req validator.StringRequest, resp *validator.StringResponse) {
	// Only validate the attribute configuration value if it is known.
	if req.ConfigValue.IsNull() || req.ConfigValue.IsUnknown() {
		return
	}

	if _, err := time.Parse(time.RFC3339, req.ConfigValue.ValueString()); err != nil {
		resp.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			req.Path,
			validator.MarkdownDescription(ctx),
			req.ConfigValue.ValueString(),
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
func IsRFC3339Time() validator.String {
	return isRFC3339TimeValidator{}
}
