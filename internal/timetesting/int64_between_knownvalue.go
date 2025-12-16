// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package timetesting

import (
	"encoding/json"
	"fmt"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ knownvalue.Check = int64Between{}

type int64Between struct {
	min int64
	max int64
}

// CheckValue determines whether the passed value is of type json.Number, converts that to an int64, and then
// checks if the value is between the min and max int64 values (inclusive).
func (v int64Between) CheckValue(other any) error {
	jsonNum, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for Int64Between check, got: %T", other)
	}

	otherVal, err := jsonNum.Int64()

	if err != nil {
		return fmt.Errorf("expected json.Number to be parseable as int64 value for Int64Between check: %s", err)
	}

	if otherVal < v.min {
		return fmt.Errorf("received value: %d, which is less than the minimum value: %d", otherVal, v.min)
	}

	if otherVal > v.max {
		return fmt.Errorf("received value: %d, which is greater than the maximum value: %d", otherVal, v.max)
	}

	return nil
}

// String returns the string representation of the value.
func (v int64Between) String() string {
	return fmt.Sprintf("%d ≤ x ≤ %d", v.min, v.max)
}

// Int64Between returns a Check for asserting that a value is between the supplied min and max int64 values (inclusive).
func Int64Between(minVal, maxVal int64) int64Between {
	return int64Between{
		min: minVal,
		max: maxVal,
	}
}
