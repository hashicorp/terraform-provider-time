// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timevalidator

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func TestIsRFC3339TimeValidator(t *testing.T) {
	t.Parallel()

	type testCase struct {
		val         types.String
		expectError bool
	}

	tests := map[string]testCase{
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
			request := validator.StringRequest{
				Path:           path.Root("test"),
				PathExpression: path.MatchRoot("test"),
				ConfigValue:    test.val,
			}

			response := validator.StringResponse{}
			IsRFC3339Time().ValidateString(context.TODO(), request, &response)

			if !response.Diagnostics.HasError() && test.expectError {
				t.Fatal("expected error, got no error")
			}

			if response.Diagnostics.HasError() && !test.expectError {
				t.Fatalf("got unexpected error: %s", response.Diagnostics)
			}
		})
	}

}
