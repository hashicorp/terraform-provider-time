// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package timetesting

import (
	"encoding/json"
	"fmt"
	"regexp"

	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
)

var _ knownvalue.Check = numberRegularExpression{}

type numberRegularExpression struct {
	regex *regexp.Regexp
}

// CheckValue determines whether the passed value is of type json.Number, converts that to a string, and then
// checks if it contains a sequence of bytes that match the regular expression supplied to NumberRegularExpression.
func (v numberRegularExpression) CheckValue(other any) error {
	otherVal, ok := other.(json.Number)

	if !ok {
		return fmt.Errorf("expected json.Number value for NumberRegularExpression check, got: %T", other)
	}

	if !v.regex.MatchString(otherVal.String()) {
		return fmt.Errorf("expected regex match %s for NumberRegularExpression check, got: %s", v.regex.String(), otherVal)
	}

	return nil
}

// String returns the string representation of the value.
func (v numberRegularExpression) String() string {
	return v.regex.String()
}

// NumberRegularExpression returns a Check for asserting equality between the
// supplied regular expression and a value passed to the CheckValue method.
func NumberRegularExpression(regex *regexp.Regexp) numberRegularExpression {
	return numberRegularExpression{
		regex: regex,
	}
}
