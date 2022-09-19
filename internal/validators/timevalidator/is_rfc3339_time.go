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

func (validator isRFC3339TimeValidator) Validate(ctx context.Context, request tfsdk.ValidateAttributeRequest, response *tfsdk.ValidateAttributeResponse) {
	t := request.AttributeConfig.Type(ctx)
	if t != types.StringType {
		response.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			request.AttributePath,
			"Expected value of type string",
			t.String(),
		))
		return
	}

	s := request.AttributeConfig.(types.String)
	if s.Unknown || s.Null {
		return
	}

	if _, err := time.Parse(time.RFC3339, s.Value); err != nil {
		response.Diagnostics.Append(validatordiag.InvalidAttributeTypeDiagnostic(
			request.AttributePath,
			validator.MarkdownDescription(ctx),
			s.Value,
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
