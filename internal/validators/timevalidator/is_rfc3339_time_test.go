package timevalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsRFC3339TimeValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         attr.Value
		expectError bool
	}

	tests := map[string]testCase{
		"not a String": {
			val:         types.BoolValue(true),
			expectError: true,
		},
		"String unknown": {
			val:         types.StringUnknown(),
			expectError: false,
		},
		"String null": {
			val:         types.StringNull(),
			expectError: false,
		},
		"not in RFC3339 format": {
			val:         types.StringValue("testString"),
			expectError: true,
		},
		"success scenario": {
			val:         types.StringValue("2022-09-06T17:47:31+00:00"),
			expectError: false,
		},
	}

	for name, test := range tests {
		name, test := name, test
		t.Run(name, func(t *testing.T) {
			request := tfsdk.ValidateAttributeRequest{
				AttributePath:           path.Root("test"),
				AttributePathExpression: path.MatchRoot("test"),
				AttributeConfig:         test.val,
			}

			response := tfsdk.ValidateAttributeResponse{}
			IsRFC3339Time().Validate(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}

}
