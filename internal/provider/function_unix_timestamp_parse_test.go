// Copyright IBM Corp. 2020, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestUnixTimestampParseFunction_Valid(t *testing.T) {
	t.Parallel()
	// Testing logic with known values
	knownUnixTime := 1690328596
	expectedKnownRFC3339 := "2023-07-25T23:43:16Z"

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`
output "test" {
  value = provider::time::unix_timestamp_parse(%d)
}
`, knownUnixTime),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test", knownvalue.ObjectExact(
							map[string]knownvalue.Check{
								"day":          knownvalue.Int64Exact(25),
								"hour":         knownvalue.Int64Exact(23),
								"iso_week":     knownvalue.Int64Exact(30),
								"iso_year":     knownvalue.Int64Exact(2023),
								"minute":       knownvalue.Int64Exact(43),
								"month":        knownvalue.Int64Exact(7),
								"month_name":   knownvalue.StringExact("July"),
								"rfc3339":      knownvalue.StringExact(expectedKnownRFC3339),
								"second":       knownvalue.Int64Exact(16),
								"weekday":      knownvalue.Int64Exact(2),
								"weekday_name": knownvalue.StringExact("Tuesday"),
								"year":         knownvalue.Int64Exact(2023),
								"year_day":     knownvalue.Int64Exact(206),
							},
						)),
					},
				},
			},
		},
	})
}

func TestUnixTimestampParseFunction_Null(t *testing.T) {
	t.Parallel()

	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
output "test" {
  value = provider::time::unix_timestamp_parse(null)
}
`,
				ExpectError: regexp.MustCompile(`Invalid value for "unix_timestamp" parameter: argument must not be null`),
			},
		},
	})
}
