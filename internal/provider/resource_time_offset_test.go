// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
)

func TestAccTimeOffset_Triggers(t *testing.T) {
	resourceName := "time_offset.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					resource.TestCheckResourceAttr(resourceName, "offset_days", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_months"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_hours"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_minutes"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
					testSleep(1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccConfigTimeOffsetTriggers1("key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "offset_days", "1"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_months"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_hours"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_minutes"),
					resource.TestCheckNoResourceAttr(resourceName, "offset_seconds"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetDays(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.AddDate(0, 0, 7)
	offsetTimestampUpdated := timestamp.AddDate(0, 0, 8)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetDays(timestamp.Format(time.RFC3339), 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetDays(timestamp.Format(time.RFC3339), 8),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_days", "8"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetHours(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.Add(1 * time.Hour)
	offsetTimestampUpdated := timestamp.Add(2 * time.Hour)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetHours(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_hours", "1"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetHours(timestamp.Format(time.RFC3339), 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_hours", "2"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetMinutes(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.Add(1 * time.Minute)
	offsetTimestampUpdated := timestamp.Add(2 * time.Minute)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetMinutes(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_minutes", "1"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetMinutes(timestamp.Format(time.RFC3339), 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_minutes", "2"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetMonths(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.AddDate(0, 3, 0)
	offsetTimestampUpdated := timestamp.AddDate(0, 4, 0)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetMonths(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_months", "3"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetMonths(timestamp.Format(time.RFC3339), 4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_months", "4"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetSeconds(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.Add(1 * time.Second)
	offsetTimestampUpdated := timestamp.Add(2 * time.Second)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetSeconds(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_seconds", "1"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetSeconds(timestamp.Format(time.RFC3339), 2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_seconds", "2"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetYears(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.AddDate(3, 0, 0)
	offsetTimestampUpdated := timestamp.AddDate(4, 0, 0)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "4"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_OffsetYearsAndMonths(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.AddDate(3, 3, 0)
	offsetTimestampUpdated := timestamp.AddDate(4, 4, 0)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetOffsetYearsAndMonths(timestamp.Format(time.RFC3339), 3, 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "offset_months", "3"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetYearsAndMonths(timestamp.Format(time.RFC3339), 4, 4),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestampUpdated.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestampUpdated.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestampUpdated.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestampUpdated.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "4"),
					resource.TestCheckResourceAttr(resourceName, "offset_months", "4"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestampUpdated.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestampUpdated.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestampUpdated.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestampUpdated.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_Upgrade(t *testing.T) {
	resourceName := "time_offset.test"
	timestamp := time.Now().UTC()
	offsetTimestamp := timestamp.AddDate(3, 0, 0)

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "base_rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "day", strconv.Itoa(offsetTimestamp.Day())),
					resource.TestCheckResourceAttr(resourceName, "hour", strconv.Itoa(offsetTimestamp.Hour())),
					resource.TestCheckResourceAttr(resourceName, "minute", strconv.Itoa(offsetTimestamp.Minute())),
					resource.TestCheckResourceAttr(resourceName, "month", strconv.Itoa(int(offsetTimestamp.Month()))),
					resource.TestCheckResourceAttr(resourceName, "offset_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", offsetTimestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", strconv.Itoa(offsetTimestamp.Second())),
					resource.TestCheckResourceAttr(resourceName, "unix", strconv.Itoa(int(offsetTimestamp.Unix()))),
					resource.TestCheckResourceAttr(resourceName, "year", strconv.Itoa(offsetTimestamp.Year())),
				),
			},
		},
	})
}

func TestAccTimeOffset_Validators(t *testing.T) {
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "time_offset" "test" {
                     base_rfc3339 = %q
                  }`, timestamp.Format(time.RFC3339)),
				ExpectError: regexp.MustCompile(`.*Error: Missing Attribute Configuration`),
			},
		},
	})
}

func testAccTimeOffsetImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		resourceName := "time_offset.test"
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		offsetYears := rs.Primary.Attributes["offset_years"]
		offsetMonths := rs.Primary.Attributes["offset_months"]
		offsetDays := rs.Primary.Attributes["offset_days"]
		offsetHours := rs.Primary.Attributes["offset_hours"]
		offsetMinutes := rs.Primary.Attributes["offset_minutes"]
		offsetSeconds := rs.Primary.Attributes["offset_seconds"]

		return fmt.Sprintf("%s,%s,%s,%s,%s,%s,%s", rs.Primary.ID, offsetYears, offsetMonths, offsetDays, offsetHours, offsetMinutes, offsetSeconds), nil
	}
}

func testAccConfigTimeOffsetTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  triggers = {
    %[1]q = %[2]q
  }
  offset_days = 1
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeOffsetOffsetDays(baseRfc3339 string, offsetDays int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339 = %[1]q
  offset_days  = %[2]d
}
`, baseRfc3339, offsetDays)
}

func testAccConfigTimeOffsetOffsetHours(baseRfc3339 string, offsetHours int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339 = %[1]q
  offset_hours = %[2]d
}
`, baseRfc3339, offsetHours)
}

func testAccConfigTimeOffsetOffsetMinutes(baseRfc3339 string, offsetMinutes int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339   = %[1]q
  offset_minutes = %[2]d
}
`, baseRfc3339, offsetMinutes)
}

func testAccConfigTimeOffsetOffsetMonths(baseRfc3339 string, offsetMonths int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339  = %[1]q
  offset_months = %[2]d
}
`, baseRfc3339, offsetMonths)
}

func testAccConfigTimeOffsetOffsetSeconds(baseRfc3339 string, offsetSeconds int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339   = %[1]q
  offset_seconds = %[2]d
}
`, baseRfc3339, offsetSeconds)
}

func testAccConfigTimeOffsetOffsetYears(baseRfc3339 string, offsetYears int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339 = %[1]q
  offset_years = %[2]d
}
`, baseRfc3339, offsetYears)
}

func testAccConfigTimeOffsetOffsetYearsAndMonths(baseRfc3339 string, offsetYears int, offsetMonths int) string {
	return fmt.Sprintf(`
resource "time_offset" "test" {
  base_rfc3339 = %[1]q
  offset_years = %[2]d
  offset_months = %[3]d
}
`, baseRfc3339, offsetYears, offsetMonths)
}