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
	time_knownvalue "github.com/hashicorp/terraform-provider-time/internal/testing/knownvalue"
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
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringRegularExpression(regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{1,2}$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d+$`))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), time_knownvalue.NumberRegularExpression(regexp.MustCompile(`^\d{4}$`))),
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
	var time1, time2 string
	resourceName := "time_static.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
					testExtractResourceAttr(resourceName, "rfc3339", &time1),
					testSleep(1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
					testExtractResourceAttr(resourceName, "rfc3339", &time2),
					testCheckAttributeValuesDiffer(&time1, &time2),
				),
			},
		},
	})
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "year", regexp.MustCompile(`^\d{4}$`)),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeStatic(),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeStatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "year", regexp.MustCompile(`^\d{4}$`)),
				),
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
