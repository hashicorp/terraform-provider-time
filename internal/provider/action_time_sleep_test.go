// Copyright IBM Corp. 2020, 2026
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// Unit test to verify that sleeping works as expected for actions
func TestActionTimeSleepInvoke(t *testing.T) {
	durationStr := "1s"
	expectedDuration, err := time.ParseDuration("1s")

	if err != nil {
		t.Fatalf("unable to parse test duration: %s", err)
	}

	sleepAction := NewTimeSleepAction()

	m := map[string]tftypes.Value{
		"duration": tftypes.NewValue(tftypes.String, durationStr),
	}
	config := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"duration": tftypes.String,
		},
	}, m)

	schemaResponse := action.SchemaResponse{}
	sleepAction.Schema(context.Background(), action.SchemaRequest{}, &schemaResponse)

	req := action.InvokeRequest{
		Config: tfsdk.Config{
			Raw:    config,
			Schema: schemaResponse.Schema,
		},
	}

	resp := action.InvokeResponse{
		Diagnostics: nil,
	}

	start := time.Now()
	sleepAction.Invoke(context.Background(), req, &resp)
	end := time.Now()
	elapsed := end.Sub(start)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error during invoke: %v", resp.Diagnostics)
	}

	if elapsed < expectedDuration {
		t.Errorf("did not sleep long enough, expected duration: %d got: %d", expectedDuration, elapsed)
	}
}

func TestAccTimeSleepAction_Basic(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepAction("1ms"),
			},
			{
				Config: testAccConfigTimeSleepAction("2ms"),
			},
		},
	})
}

func TestAccTimeSleepAction_Validators(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config:      testAccConfigTimeSleepAction("1"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Value Match`),
			},
			{
				Config:      testAccConfigTimeSleepAction("invalid"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Value Match`),
			},
			{
				Config:      testAccConfigTimeSleepAction("30"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Value Match`),
			},
		},
	})
}

func TestAccTimeSleepAction_ValidDurations(t *testing.T) {
	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepAction("100ms"),
			},
			{
				Config: testAccConfigTimeSleepAction("1s"),
			},
			{
				Config: testAccConfigTimeSleepAction("1m"),
			},
			{
				Config: testAccConfigTimeSleepAction("1h"),
			},
			{
				Config: testAccConfigTimeSleepAction("1.5s"),
			},
		},
	})
}

func testAccConfigTimeSleepAction(duration string) string {
	return fmt.Sprintf(`
terraform {
  required_providers {
    time = {
      source = "hashicorp/time"
    }
  }
}

provider "time" {}

action "time_sleep" "test" {
  config {
    duration = %[1]q
  }
}
`, duration)
}
