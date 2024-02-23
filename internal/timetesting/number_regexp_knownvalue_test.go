// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetesting_test

import (
	"encoding/json"
	"fmt"
	"regexp"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-provider-time/internal/timetesting"
)

func TestNumberRegularExpression_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-nil": {
			self:          timetesting.NumberRegularExpression(regexp.MustCompile("")),
			expectedError: fmt.Errorf("expected json.Number value for NumberRegularExpression check, got: <nil>"),
		},
		"zero-other": {
			self:  timetesting.NumberRegularExpression(regexp.MustCompile("")),
			other: json.Number(""), // checking against the underlying value field zero-value
		},
		"nil": {
			self:          timetesting.NumberRegularExpression(regexp.MustCompile("1.23")),
			expectedError: fmt.Errorf("expected json.Number value for NumberRegularExpression check, got: <nil>"),
		},
		"wrong-type": {
			self:          timetesting.NumberRegularExpression(regexp.MustCompile("1.23")),
			other:         "1.23",
			expectedError: fmt.Errorf("expected json.Number value for NumberRegularExpression check, got: string"),
		},
		"not-equal": {
			self:          timetesting.NumberRegularExpression(regexp.MustCompile("1.23")),
			other:         json.Number("1.24"),
			expectedError: fmt.Errorf("expected regex match 1.23 for NumberRegularExpression check, got: 1.24"),
		},
		"equal": {
			self:  timetesting.NumberRegularExpression(regexp.MustCompile("1.23")),
			other: json.Number("1.23"),
		},
	}

	for name, testCase := range testCases {
		name, testCase := name, testCase

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.self.CheckValue(testCase.other)

			if diff := cmp.Diff(got, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestNumberRegularExpression_String(t *testing.T) {
	t.Parallel()

	got := timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`)).String()

	if diff := cmp.Diff(got, `^\d{1,2}$`); diff != "" {
		t.Errorf("unexpected difference: %s", diff)
	}
}

// equateErrorMessage reports errors to be equal if both are nil
// or both have the same message.
var equateErrorMessage = cmp.Comparer(func(x, y error) bool {
	if x == nil || y == nil {
		return x == nil && y == nil
	}
	return x.Error() == y.Error()
})
