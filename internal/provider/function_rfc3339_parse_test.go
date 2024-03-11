// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestRFC3339Parse_UTC(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("2023-07-25T23:43:16Z")
				}
				`,
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
								"second":       knownvalue.Int64Exact(16),
								"unix":         knownvalue.Int64Exact(1690328596),
								"weekday":      knownvalue.Int64Exact(2),
								"weekday_name": knownvalue.StringExact("Tuesday"),
								"year":         knownvalue.Int64Exact(2023),
								"year_day":     knownvalue.Int64Exact(206),
							},
						)),
					},
				},
			},
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("2023-07-25T23:43:16-00:00")
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
					value = provider::time::rfc3339_parse("2023-07-25T23:43:16+00:00")
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

func TestRFC3339Parse_offset(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("1996-12-19T16:39:57-08:00")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test", knownvalue.ObjectExact(
							map[string]knownvalue.Check{
								"day":          knownvalue.Int64Exact(19),
								"hour":         knownvalue.Int64Exact(16),
								"iso_week":     knownvalue.Int64Exact(51),
								"iso_year":     knownvalue.Int64Exact(1996),
								"minute":       knownvalue.Int64Exact(39),
								"month":        knownvalue.Int64Exact(12),
								"month_name":   knownvalue.StringExact("December"),
								"second":       knownvalue.Int64Exact(57),
								"unix":         knownvalue.Int64Exact(851042397),
								"weekday":      knownvalue.Int64Exact(4),
								"weekday_name": knownvalue.StringExact("Thursday"),
								"year":         knownvalue.Int64Exact(1996),
								"year_day":     knownvalue.Int64Exact(354),
							},
						)),
					},
				},
			},
		},
	})
}

func TestRFC3339Parse_invalid(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// TODO: Replace with the stable v1.8.0 release when available
			tfversion.SkipBelow(version.Must(version.NewVersion("v1.8.0-beta1"))),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("abcdef")
				}
				`,
				ExpectError: regexp.MustCompile(`"abcdef" is not a valid RFC3339 timestamp.`),
			},
		},
	})
}
