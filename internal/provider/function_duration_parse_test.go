// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestDurationParse_valid(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("1h")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test", knownvalue.ObjectExact(
							map[string]knownvalue.Check{
								"hours":        knownvalue.Float64Exact(1),
								"minutes":      knownvalue.Float64Exact(60),
								"seconds":      knownvalue.Float64Exact(3600),
								"milliseconds": knownvalue.Int64Exact(3600000),
								"microseconds": knownvalue.Int64Exact(3600000000),
								"nanoseconds":  knownvalue.Int64Exact(3600000000000),
							},
						)),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("60m")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("3600s")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("3600000ms")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("3600000000us")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("3600000000000ns")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
		},
	})
}

func TestDurationParse_invalid(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_8_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::duration_parse("abcdef")
				}
				`,
				ExpectError: regexp.MustCompile(`"abcdef" is not a valid duration string.`),
			},
		},
	})
}
