// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-time/internal/timetesting"
)

func TestAccTimeStatic_basic(t *testing.T) {
	resourceName := "time_static.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStatic(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringRegularExpression(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d+$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{4}$`))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeStatic_Triggers(t *testing.T) {

	resourceName := "time_static.test"

	captureTimeState1 := timetesting.NewExtractState(resourceName, tfjsonpath.New("rfc3339"))
	captureTimeState2 := timetesting.NewExtractState(resourceName, tfjsonpath.New("rfc3339"))

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
					captureTimeState1,
					timetesting.Sleep(2),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1updated")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
					captureTimeState2,
				},
			},
		},
	})

	if captureTimeState1.Value == captureTimeState2.Value {
		t.Fatal("attribute values are the same")
	}
}

func TestAccTimeStatic_Rfc3339(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	rfc3339 := knownvalue.StringExact(timestamp.Format(time.RFC3339))
	year := knownvalue.Int64Exact(int64(timestamp.Year()))
	month := knownvalue.Int64Exact(int64(timestamp.Month()))
	day := knownvalue.Int64Exact(int64(timestamp.Day()))
	hour := knownvalue.Int64Exact(int64(timestamp.Hour()))
	minute := knownvalue.Int64Exact(int64(timestamp.Minute()))
	second := knownvalue.Int64Exact(int64(timestamp.Second()))
	unix := knownvalue.Int64Exact(timestamp.Unix())

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticRfc3339(timestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), year),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), month),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), day),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), hour),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), minute),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), second),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), unix),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), rfc3339),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), rfc3339),
					},
				},

				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), year),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), month),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), day),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), hour),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), minute),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), second),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), unix),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), rfc3339),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), rfc3339),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeStatic_Upgrade(t *testing.T) {
	resourceName := "time_static.test"

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeStatic(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringRegularExpression(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d+$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{4}$`))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeStatic(),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeStatic(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringRegularExpression(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d+$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), timetesting.NumberRegularExpression(regexp.MustCompile(`^\d{4}$`))),
				},
			},
		},
	})
}

func TestAccTimeStatic_Validators(t *testing.T) {
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccConfigTimeStaticRfc3339(timestamp.Format(time.RFC850)),
				ExpectError: regexp.MustCompile(`.*Invalid RFC3339 String Value`),
			},
		},
	})
}

func testAccConfigTimeStatic() string {
	return `
resource "time_static" "test" {}
`
}

func testAccConfigTimeStaticTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
 triggers = {
   %[1]q = %[2]q
 }
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeStaticRfc3339(rfc3339 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
 rfc3339 = %[1]q
}
`, rfc3339)
}
