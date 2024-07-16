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
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestAccTimeOffset_Triggers(t *testing.T) {
	resourceName := "time_offset.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeOffsetTriggers1("key1", "value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_days"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_seconds"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				// Ensures a time difference when running unit tests in CI
				PreConfig: func() {
					time.Sleep(time.Duration(1) * time.Second)
				},
				Config: testAccConfigTimeOffsetTriggers1("key1", "value1updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1updated")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_days"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_seconds"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_days"), knownvalue.Int64Exact(7)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetDays(timestamp.Format(time.RFC3339), 8),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_days"), knownvalue.Int64Exact(8)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_hours"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetHours(timestamp.Format(time.RFC3339), 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_hours"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_minutes"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetMinutes(timestamp.Format(time.RFC3339), 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_minutes"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetMonths(timestamp.Format(time.RFC3339), 4),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_seconds"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetSeconds(timestamp.Format(time.RFC3339), 2),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_seconds"), knownvalue.Int64Exact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 4),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeOffsetImportStateIdFunc(),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeOffsetOffsetYearsAndMonths(timestamp.Format(time.RFC3339), 4, 4),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_months"), knownvalue.Int64Exact(4)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestampUpdated.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestampUpdated.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestampUpdated.Year()))),
				},
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
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeOffsetOffsetYears(timestamp.Format(time.RFC3339), 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("base_rfc3339"), knownvalue.StringExact(timestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("day"), knownvalue.Int64Exact(int64(offsetTimestamp.Day()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("hour"), knownvalue.Int64Exact(int64(offsetTimestamp.Hour()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("minute"), knownvalue.Int64Exact(int64(offsetTimestamp.Minute()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("month"), knownvalue.Int64Exact(int64(offsetTimestamp.Month()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("offset_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(offsetTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("second"), knownvalue.Int64Exact(int64(offsetTimestamp.Second()))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("unix"), knownvalue.Int64Exact(offsetTimestamp.Unix())),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("year"), knownvalue.Int64Exact(int64(offsetTimestamp.Year()))),
				},
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
