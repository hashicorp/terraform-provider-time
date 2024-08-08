// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"regexp"
	"testing"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	"github.com/hashicorp/terraform-plugin-testing/config"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestAccTimeRotating_Triggers(t *testing.T) {
	resourceName := "time_rotating.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingTriggers1("key1", "value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				// Ensures a time difference when running unit tests in CI
				PreConfig: func() {
					time.Sleep(time.Duration(1) * time.Second)
				},
				Config: testAccConfigTimeRotatingTriggers1("key1", "value1updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1updated")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationDays_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationDays := 7

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationDays(rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationDate := now.AddDate(0, 0, -8)
					mockClock.Increment(mockClock.Since(rotationDate)) // 8 days
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationDays(rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(0, 0, rotationDays+8).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.AddDate(0, 0, 8).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationDays_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationDays := 7

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationDays(baseTimestamp.Format(time.RFC3339), rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationDate := baseTimestamp.AddDate(0, 0, -8)
					mockClock.Increment(mockClock.Since(rotationDate)) // 8 days
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationDays(baseTimestamp.Format(time.RFC3339), rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is computed using the configured rfc3339 and rotation_days attributes,
					// the value of rotation_rfc3339 is exactly the same as before the rotation. Since the rotation_rfc3339
					// timestamp will be expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationDays_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	expiredTimestamp := time.Now().UTC().AddDate(0, 0, -2)
	rotationDays := 1

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationDays(expiredTimestamp.Format(time.RFC3339), rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(expiredTimestamp.AddDate(0, 0, 1).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationHours_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationHours := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationHours(rotationHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.Add(time.Duration(rotationHours)*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					mockClock.Increment(4 * time.Hour)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationHours(rotationHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.Add((4+time.Duration(rotationHours))*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Add(4*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationHours_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationHours := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationHours(baseTimestamp.Format(time.RFC3339), rotationHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.Add(time.Duration(rotationHours)*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					mockClock.Increment(4 * time.Hour)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationHours(baseTimestamp.Format(time.RFC3339), rotationHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is computed using the configured rfc3339 and rotation_hours attributes,
					// the value of rotation_rfc3339 is exactly the same as before the rotation. Since the rotation_rfc3339
					// timestamp will be expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.Add(time.Duration(rotationHours)*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationHours_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	expiredTimestamp := time.Now().UTC().Add(-5 * time.Hour)
	rotationHours := 3

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationHours(expiredTimestamp.Format(time.RFC3339), rotationHours),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(expiredTimestamp.Add(time.Duration(rotationHours)*time.Hour).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Int64Exact(int64(rotationHours))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationMinutes_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationMinutes := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationMinutes(rotationMinutes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.Add(time.Duration(rotationMinutes)*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					mockClock.Increment(4 * time.Minute)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationMinutes(rotationMinutes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.Add((time.Duration(rotationMinutes)+4)*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Add(4*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationMinutes_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationMinutes := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationMinutes(baseTimestamp.Format(time.RFC3339), rotationMinutes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.Add(time.Duration(rotationMinutes)*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					mockClock.Increment(4 * time.Minute)
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationMinutes(baseTimestamp.Format(time.RFC3339), rotationMinutes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is computed using the configured rfc3339 and rotation_minutes attributes,
					// the value of rotation_rfc3339 is exactly the same as before the rotation. Since the rotation_rfc3339
					// timestamp will be expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.Add(time.Duration(rotationMinutes)*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationMinutes_expired(t *testing.T) {
	resourceName := "time_rotating.test"

	expiredTimestamp := time.Now().UTC().Add(-2 * time.Minute)
	rotationMinutes := 1

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationMinutes(expiredTimestamp.Format(time.RFC3339), rotationMinutes),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(expiredTimestamp.Add(time.Duration(rotationMinutes)*time.Minute).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Int64Exact(int64(rotationMinutes))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationMonths_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationMonths := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationMonths(rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationDate := now.AddDate(0, -4, 0)
					mockClock.Increment(mockClock.Since(rotationDate)) // 4 months
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationMonths(rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(0, rotationMonths+4, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.AddDate(0, 4, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationMonths_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationMonths := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationMonths(baseTimestamp.Format(time.RFC3339), rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationTimestamp := baseTimestamp.AddDate(0, -4, 0)
					mockClock.Increment(mockClock.Since(rotationTimestamp)) // 4 months
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationMonths(baseTimestamp.Format(time.RFC3339), rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is computed using the configured rfc3339 and rotation_months attributes,
					// the value of rotation_rfc3339 is exactly the same as before the rotation. Since the rotation_rfc3339
					// timestamp will be expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationMonths_expired(t *testing.T) {
	resourceName := "time_rotating.test"

	expiredTimestamp := time.Now().UTC().AddDate(0, -2, 0)
	rotationMonths := 1

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationMonths(expiredTimestamp.Format(time.RFC3339), rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(expiredTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationRfc3339_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationTimestamp := time.Now().UTC().AddDate(0, 0, 7)

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationRfc3339(rotationTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			{
				PreConfig: func() {
					rotationDate := now.AddDate(0, 0, -8)
					mockClock.Increment(mockClock.Since(rotationDate)) // 8 days
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationRfc3339(rotationTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is configured, it's value is exactly
					// the same as before the rotation. Since the rotation_rfc3339 timestamp will be
					// expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.AddDate(0, 0, 8).Format(time.RFC3339))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationRfc3339_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationTimestamp := time.Now().UTC().AddDate(0, 0, 7)

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationRfc3339(baseTimestamp.Format(time.RFC3339), rotationTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			{
				PreConfig: func() {
					rotationTimestamp := baseTimestamp.AddDate(0, 0, -8)
					mockClock.Increment(mockClock.Since(rotationTimestamp)) // 8 days
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationRfc3339(baseTimestamp.Format(time.RFC3339), rotationTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 and rfc3339 are configured, their values are exactly
					// the same as before the rotation. Since the rotation_rfc3339 timestamp will be
					// expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC339RotationRfc3339_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	baseTimestamp := time.Now().UTC().AddDate(0, 0, -2)
	rotationTimestamp := time.Now().UTC().AddDate(0, 0, -1)

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationRfc3339(baseTimestamp.Format(time.RFC3339), rotationTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 and rfc3339 are configured, their values are exactly
					// the same as before the rotation. Since the rotation_rfc3339 timestamp will be
					// expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(rotationTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339RotationYears_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	now := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(now)
	rotationYears := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationYears(rotationYears),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(rotationYears, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationTimestamp := now.AddDate(-4, 0, 0)
					mockClock.Increment(mockClock.Since(rotationTimestamp)) // 4 years
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationYears(rotationYears),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(now.AddDate(4+rotationYears, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(now.AddDate(4, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationYears_basic(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationYears := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationYears(baseTimestamp.Format(time.RFC3339), rotationYears),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(rotationYears, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				ResourceName:             resourceName,
				ImportState:              true,
				ImportStateIdFunc:        testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:        true,
			},
			// Trigger a rotation
			{
				PreConfig: func() {
					rotationTimestamp := baseTimestamp.AddDate(-4, 0, 0)
					mockClock.Increment(mockClock.Since(rotationTimestamp)) // 4 years
				},
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRFC3339RotationYears(baseTimestamp.Format(time.RFC3339), rotationYears),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						// Rotations have the "Create" action since the rotation checking logic is run
						// during ReadResource() and the resource is removed from state.
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					// Because rotation_rfc3339 is computed using the configured rfc3339 and rotation_years attributes,
					// the value of rotation_rfc3339 is exactly the same as before the rotation. Since the rotation_rfc3339
					// timestamp will be expired after this first rotation, every subsequent refresh of the resource
					// will trigger a rotation.
					// Ref: https://github.com/hashicorp/terraform-provider-time/issues/44
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(rotationYears, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339RotationYears_expired(t *testing.T) {
	resourceName := "time_rotating.test"

	expiredTimestamp := time.Now().UTC().AddDate(-2, 0, 0)
	rotationYears := 1

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationYears(expiredTimestamp.Format(time.RFC3339), rotationYears),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
					// Since the rotation_rfc3339 timestamp is expired, every subsequent refresh of the resource
					// will trigger a rotation, creating the resource with the exact same plan/state values.
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(expiredTimestamp.AddDate(rotationYears, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(expiredTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(int64(rotationYears))),
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_ComputedRFC3339_RotationDays_ToRotationMonths(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	mockClock := fakeclock.NewFakeClock(baseTimestamp)
	rotationDays := 7
	rotationMonths := 3

	resource.UnitTest(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationDays(rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactoriesTestProvider(mockClock),
				Config:                   testAccConfigTimeRotatingRotationMonths(rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
			},
		},
	})
}

func TestAccTimeRotating_ConfiguredRFC3339_RotationDays_ToRotationMonths(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	rotationDays := 7
	rotationMonths := 3

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationDays(baseTimestamp.Format(time.RFC3339), rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
			{
				Config: testAccConfigTimeRotatingRFC3339RotationMonths(baseTimestamp.Format(time.RFC3339), rotationMonths),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, rotationMonths, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Int64Exact(int64(rotationMonths))),
				},
			},
		},
	})
}

// When the resource is being updated, the "rotation_rfc3339" value
// is computed during ModifyPlan(). If any of the "rotation_" attributes
// are unknown during the plan, then an incorrect value will be calculated
// for "rotation_rfc3339" during the initial plan which will differ from
// the final plan, causing Terraform core to throw an error
// Ref: https://github.com/hashicorp/terraform-provider-time/issues/227
func TestAccTimeRotating_UpdateUnknownValue(t *testing.T) {
	resourceName := "time_rotating.test"

	baseTimestamp := time.Now().UTC()
	rotationDays := 7

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRFC3339RotationDays(baseTimestamp.Format(time.RFC3339), rotationDays),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionCreate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_rfc3339")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
					},
				},
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(baseTimestamp.AddDate(0, 0, rotationDays).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Int64Exact(int64(rotationDays))),
				},
			},
			{
				Config: fmt.Sprintf(`resource "time_static" "unknown" {}

					resource "time_rotating" "test" {
  							rotation_days   = time_static.unknown.day
  							rfc3339         = %q
						}`, baseTimestamp.Format(time.RFC3339)),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionUpdate),
						plancheck.ExpectUnknownValue(resourceName, tfjsonpath.New("rotation_days")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact("0001-01-01T00:00:00Z")),
						plancheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.StringExact(baseTimestamp.Format(time.RFC3339))),
					},
				},
				ExpectError: regexp.MustCompile(`.*Error: Provider produced inconsistent final plan`),
			},
		},
	})
}

// The "time_rotating" resource predates the `replace_triggered_by`
// lifestyle argument introduced in Terraform v1.2.0. The `replace_triggered_by` argument looks
// for an update or replacement of the supplied resource instance. Because the "time_rotating" rotation
// checking logic is run during ReadResource() and the resource is removed from state,
// a rotation is considered to be a creation of a new resource rather than an update or replacement.
// Ref: https://github.com/hashicorp/terraform-provider-time/issues/118
func TestAccTimeRotating_LifecycleReplaceTriggeredBy(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			// "terraform_data" resource is only available in Terraform v1.4.0 and above.
			tfversion.SkipBelow(tfversion.Version1_4_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				ConfigFile: config.TestNameFile("test.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("time_rotating.configured_rfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.configured_rotationrfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.computed_rotation", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rotationrfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("terraform_data.test_computed_rotation", plancheck.ResourceActionCreate),
					},
				},
			},
			{
				ConfigFile: config.TestNameFile("test.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectEmptyPlan(),
					},
				},
			},
			{
				PreConfig: func() {
					time.Sleep(time.Duration(1) * time.Minute)
				},
				ConfigFile: config.TestNameFile("test.tf"),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("time_rotating.configured_rfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.configured_rotationrfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.computed_rotation", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rfc3339", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rotationrfc3339", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.test_computed_rotation", plancheck.ResourceActionNoop),
					},
					PostApplyPostRefresh: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction("time_rotating.configured_rfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.configured_rotationrfc3339", plancheck.ResourceActionCreate),
						plancheck.ExpectResourceAction("time_rotating.computed_rotation", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.test_computed_rotation", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rfc3339", plancheck.ResourceActionNoop),
						plancheck.ExpectResourceAction("terraform_data.test_configured_rotationrfc3339", plancheck.ResourceActionNoop),
					},
				},
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotation_Upgrade(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()
	expiredTimestamp := time.Now().UTC().AddDate(-2, 0, 0)

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeRotatingRFC3339RotationYears(timestamp.Format(time.RFC3339), 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(timestamp.AddDate(3, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeRotatingRFC3339RotationYears(timestamp.Format(time.RFC3339), 3),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeRotatingRFC3339RotationYears(timestamp.Format(time.RFC3339), 3),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_years"), knownvalue.Int64Exact(3)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_rfc3339"), knownvalue.StringExact(timestamp.AddDate(3, 0, 0).Format(time.RFC3339))),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_months"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_days"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_hours"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rotation_minutes"), knownvalue.Null()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("rfc3339"), knownvalue.NotNull()),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeRotatingRFC3339RotationYears(expiredTimestamp.Format(time.RFC3339), 3),
				PlanOnly:                 true,
				ExpectNonEmptyPlan:       true,
			},
		},
	})
}

func TestAccTimeRotating_Validators(t *testing.T) {
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "time_rotating" "test" {
                     rfc3339 = %q
                  }`, timestamp.Format(time.RFC3339)),
				ExpectError: regexp.MustCompile(`.*Error: Missing Attribute Configuration`),
			},
			{
				Config:      testAccConfigTimeRotatingRFC3339RotationMinutes(timestamp.Format(time.RFC822), 1),
				ExpectError: regexp.MustCompile(`.*Invalid RFC3339 String Value`),
			},
			{
				Config:      testAccConfigTimeRotatingRFC3339RotationMinutes(timestamp.Format(time.RFC3339), 0),
				ExpectError: regexp.MustCompile(`.*must be at least 1`),
			},
		},
	})
}

func testAccTimeRotatingImportStateIdFunc() resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		resourceName := "time_rotating.test"
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		rotationYears := rs.Primary.Attributes["rotation_years"]
		rotationMonths := rs.Primary.Attributes["rotation_months"]
		rotationDays := rs.Primary.Attributes["rotation_days"]
		rotationHours := rs.Primary.Attributes["rotation_hours"]
		rotationMinutes := rs.Primary.Attributes["rotation_minutes"]

		if rotationYears != "" || rotationMonths != "" || rotationDays != "" || rotationHours != "" || rotationMinutes != "" {
			return fmt.Sprintf("%s,%s,%s,%s,%s,%s", rs.Primary.ID, rotationYears, rotationMonths, rotationDays, rotationHours, rotationMinutes), nil
		}

		return fmt.Sprintf("%s,%s", rs.Primary.ID, rs.Primary.Attributes["rotation_rfc3339"]), nil
	}
}

func testAccConfigTimeRotatingTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  triggers = {
    %[1]q = %[2]q
  }
  rotation_days = 1
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeRotatingRFC3339RotationDays(rfc3339 string, rotationDays int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_days = %[2]d
  rfc3339       = %[1]q
}
`, rfc3339, rotationDays)
}

func testAccConfigTimeRotatingRotationDays(rotationDays int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_days = %d
}
`, rotationDays)
}

func testAccConfigTimeRotatingRFC3339RotationHours(rfc3339 string, rotationHours int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_hours = %[2]d
  rfc3339        = %[1]q
}
`, rfc3339, rotationHours)
}

func testAccConfigTimeRotatingRotationHours(rotationHours int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_hours = %d
}
`, rotationHours)
}

func testAccConfigTimeRotatingRFC3339RotationMinutes(rfc3339 string, rotationMinutes int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_minutes = %[2]d
  rfc3339          = %[1]q
}
`, rfc3339, rotationMinutes)
}

func testAccConfigTimeRotatingRotationMinutes(rotationMinutes int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_minutes = %d
}
`, rotationMinutes)
}

func testAccConfigTimeRotatingRotationMonths(rotationMonths int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_months = %d
}
`, rotationMonths)
}

func testAccConfigTimeRotatingRFC3339RotationMonths(rfc3339 string, rotationMonths int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_months = %[2]d
  rfc3339         = %[1]q
}
`, rfc3339, rotationMonths)
}

func testAccConfigTimeRotatingRotationYears(rotationYears int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_years = %d
}
`, rotationYears)
}

func testAccConfigTimeRotatingRFC3339RotationYears(rfc3339 string, rotationYears int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_years = %[2]d
  rfc3339        = %[1]q
}
`, rfc3339, rotationYears)
}

func testAccConfigTimeRotatingRFC3339RotationRfc3339(rfc3339 string, rotationRfc3339 string) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_rfc3339 = %[2]q
  rfc3339          = %[1]q
}
`, rfc3339, rotationRfc3339)
}

func testAccConfigTimeRotatingRotationRfc3339(rotationRfc3339 string) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_rfc3339 = %q
}
`, rotationRfc3339)
}
