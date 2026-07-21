// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/echoprovider"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"

	"github.com/hashicorp/terraform-provider-time/internal/clock"
)

// Since the acceptance testing framework can introduce uncontrollable time delays,
// verify that sleeping works as expected via unit testing.
func TestEphemeralResourceTimeSleep_Open(t *testing.T) {
	t.Parallel()

	durationStr := "1s"
	expectedDuration, err := time.ParseDuration(durationStr)
	if err != nil {
		t.Fatalf("unable to parse test duration: %s", err)
	}

	sleepResource := NewTimeSleepEphemeralResource().(*timeSleepEphemeralResource)
	configureReq := ephemeral.ConfigureRequest{
		ProviderData: clock.NewClock(),
	}
	sleepResource.Configure(context.Background(), configureReq, &ephemeral.ConfigureResponse{})

	m := map[string]tftypes.Value{
		"open_duration":  tftypes.NewValue(tftypes.String, durationStr),
		"close_duration": tftypes.NewValue(tftypes.String, durationStr),
		"outputs":        tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	}
	typ := tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"open_duration":  tftypes.String,
			"close_duration": tftypes.String,
			"outputs":        tftypes.Map{ElementType: tftypes.String},
		},
		OptionalAttributes: map[string]struct{}{
			"open_duration":  {},
			"close_duration": {},
			"outputs":        {},
		},
	}
	config := tftypes.NewValue(typ, m)

	schemaResponse := ephemeral.SchemaResponse{}
	sleepResource.Schema(context.Background(), ephemeral.SchemaRequest{}, &schemaResponse)

	req := ephemeral.OpenRequest{
		Config: tfsdk.Config{
			Raw:    config,
			Schema: schemaResponse.Schema,
		},
	}
	resp := ephemeral.OpenResponse{
		Result: tfsdk.EphemeralResultData{
			Raw:    tftypes.NewValue(schemaResponse.Schema.Type().TerraformType(context.Background()), tftypes.UnknownValue),
			Schema: schemaResponse.Schema,
		},
		// There isn't a way to construct the Private field here, hence we can't test for the close_duration
	}
	start := time.Now()
	sleepResource.Open(context.Background(), req, &resp)
	end := time.Now()
	elapsed := end.Sub(start)
	if elapsed < expectedDuration {
		t.Errorf("open did not sleep long enough, expected duration: %v got: %v", expectedDuration, elapsed)
	}
}

func TestAccEphemeralTimeSleep_Outputs(t *testing.T) {
	t.Parallel()

	factories := protoV6ProviderFactories()
	factories["echo"] = echoprovider.NewProviderServer()

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: factories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigEphemeralTimeSleepOutputs(),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.test", tfjsonpath.New("data").AtMapKey("foo"), knownvalue.StringExact("bar")),
				},
			},
		},
	})
}

func testAccConfigEphemeralTimeSleepOutputs() string {
	return `
ephemeral "time_sleep" "test" {
  close_duration = "0.1s"

  outputs = {
	foo = "bar"
  }
}

provider "echo" {
  data = ephemeral.time_sleep.test.outputs
}

resource "echo" "test" {}
`
}
