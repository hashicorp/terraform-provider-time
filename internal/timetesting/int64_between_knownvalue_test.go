// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetesting_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-provider-time/internal/timetesting"
)

func TestInt64Between_CheckValue(t *testing.T) {
	t.Parallel()

	testCases := map[string]struct {
		self          knownvalue.Check
		other         any
		expectedError error
	}{
		"zero-other": {
			self:  timetesting.Int64Between(0, 0),
			other: json.Number("0"), // checking against the underlying value field zero-value
		},
		"nil": {
			self:          timetesting.Int64Between(0, 100),
			expectedError: fmt.Errorf("expected json.Number value for Int64Between check, got: <nil>"),
		},
		"wrong-type": {
			self:          timetesting.Int64Between(0, 100),
			other:         json.Number("str"),
			expectedError: fmt.Errorf("expected json.Number to be parseable as int64 value for Int64Between check: strconv.ParseInt: parsing \"str\": invalid syntax"),
		},
		"less-than-min": {
			self:          timetesting.Int64Between(1, 50),
			other:         json.Number("0"),
			expectedError: fmt.Errorf("received value: 0, which is less than the minimum value: 1"),
		},
		"greater-than-max": {
			self:          timetesting.Int64Between(1, 50),
			other:         json.Number("51"),
			expectedError: fmt.Errorf("received value: 51, which is greater than the maximum value: 50"),
		},
		"between": {
			self:  timetesting.Int64Between(1, 50),
			other: json.Number("35"),
		},
		"between-equal-to-min": {
			self:  timetesting.Int64Between(1, 50),
			other: json.Number("1"),
		},
		"between-equal-to-max": {
			self:  timetesting.Int64Between(1, 50),
			other: json.Number("50"),
		},
	}

	for name, testCase := range testCases {

		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got := testCase.self.CheckValue(testCase.other)

			if diff := cmp.Diff(got, testCase.expectedError, equateErrorMessage); diff != "" {
				t.Errorf("unexpected difference: %s", diff)
			}
		})
	}
}

func TestInt64Between_String(t *testing.T) {
	t.Parallel()

	got := timetesting.Int64Between(0, 100).String()

	if diff := cmp.Diff(got, `0 ≤ x ≤ 100`); diff != "" {
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
