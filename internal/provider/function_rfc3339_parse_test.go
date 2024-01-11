package provider

import (
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
)

func TestRFC3339Parse_basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("2023-07-25T23:43:16Z")
				}
				`,
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownOutputValue("test", knownvalue.ObjectValueExact(
							map[string]knownvalue.Check{
								"day":          knownvalue.Int64ValueExact(25),
								"hour":         knownvalue.Int64ValueExact(23),
								"iso_week":     knownvalue.Int64ValueExact(30),
								"iso_year":     knownvalue.Int64ValueExact(2023),
								"minute":       knownvalue.Int64ValueExact(43),
								"month":        knownvalue.Int64ValueExact(7),
								"month_name":   knownvalue.StringValueExact("July"),
								"second":       knownvalue.Int64ValueExact(16),
								"unix":         knownvalue.Int64ValueExact(1690328596),
								"weekday":      knownvalue.Int64ValueExact(2),
								"weekday_name": knownvalue.StringValueExact("Tuesday"),
								"year":         knownvalue.Int64ValueExact(2023),
								"year_day":     knownvalue.Int64ValueExact(206),
								"zone_name":    knownvalue.StringValueExact("UTC"),
								"zone_offset":  knownvalue.Int64ValueExact(0),
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
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: `
				output "test" {
					value = provider::time::rfc3339_parse("invalid")
				}
				`,
				// TODO: fix error message
				ExpectError: regexp.MustCompile(`blah`),
			},
		},
	})
}
