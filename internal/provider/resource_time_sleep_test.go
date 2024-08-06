// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-framework-timetypes/timetypes"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	r "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
	"github.com/hashicorp/terraform-provider-time/internal/timetesting"
)

// Since the acceptance testing framework can introduce uncontrollable time delays,
// verify that sleeping works as expected via unit testing.
func TestResourceTimeSleepCreate(t *testing.T) {
	durationStr := "1s"
	expectedDuration, err := time.ParseDuration("1s")

	if err != nil {
		t.Fatalf("unable to parse test duration: %s", err)
	}

	sleepResource := NewTimeSleepResource()

	m := map[string]tftypes.Value{
		"create_duration":  tftypes.NewValue(tftypes.String, durationStr),
		"destroy_duration": tftypes.NewValue(tftypes.String, nil),
		"id":               tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"triggers":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	}
	config := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"create_duration":  tftypes.String,
			"destroy_duration": tftypes.String,
			"id":               tftypes.String,
			"triggers":         tftypes.Map{ElementType: tftypes.String},
		},
		OptionalAttributes: map[string]struct{}{
			"create_duration":  {},
			"destroy_duration": {},
			"triggers":         {},
		},
	}, m)
	plan := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"create_duration":  tftypes.String,
			"destroy_duration": tftypes.String,
			"id":               tftypes.String,
			"triggers":         tftypes.Map{ElementType: tftypes.String},
		},
		OptionalAttributes: map[string]struct{}{
			"create_duration":  {},
			"destroy_duration": {},
			"triggers":         {},
		},
	}, m)

	schemaResponse := r.SchemaResponse{}
	sleepResource.Schema(context.Background(), r.SchemaRequest{}, &schemaResponse)

	req := r.CreateRequest{
		Config: tfsdk.Config{
			Raw:    config,
			Schema: schemaResponse.Schema,
		},
		Plan: tfsdk.Plan{
			Raw:    plan,
			Schema: schemaResponse.Schema,
		},
		ProviderMeta: tfsdk.Config{},
	}

	resp := r.CreateResponse{
		State: tfsdk.State{
			Schema: schemaResponse.Schema,
		},
		Diagnostics: nil,
	}

	start := time.Now()
	sleepResource.Create(context.Background(), req, &resp)
	end := time.Now()
	elapsed := end.Sub(start)

	if elapsed < expectedDuration {
		t.Errorf("did not sleep long enough, expected duration: %d got: %d", expectedDuration, elapsed)
	}
}

// Since the acceptance testing framework can introduce uncontrollable time delays,
// verify that sleeping works as expected via unit testing.
func TestResourceTimeSleepDelete(t *testing.T) {
	durationStr := "1s"
	expectedDuration, err := time.ParseDuration("1s")

	if err != nil {
		t.Fatalf("unable to parse test duration: %s", err)
	}

	sleepResource := NewTimeSleepResource()

	m := map[string]tftypes.Value{
		"create_duration":  tftypes.NewValue(tftypes.String, nil),
		"destroy_duration": tftypes.NewValue(tftypes.String, durationStr),
		"id":               tftypes.NewValue(tftypes.String, tftypes.UnknownValue),
		"triggers":         tftypes.NewValue(tftypes.Map{ElementType: tftypes.String}, nil),
	}

	config := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"create_duration":  tftypes.String,
			"destroy_duration": tftypes.String,
			"id":               tftypes.String,
			"triggers":         tftypes.Map{ElementType: tftypes.String},
		},
		OptionalAttributes: map[string]struct{}{
			"create_duration":  {},
			"destroy_duration": {},
			"triggers":         {},
		},
	}, m)

	schemaResponse := r.SchemaResponse{}
	sleepResource.Schema(context.Background(), r.SchemaRequest{}, &schemaResponse)

	req := r.DeleteRequest{
		State: tfsdk.State{
			Raw:    config,
			Schema: schemaResponse.Schema,
		},
		ProviderMeta: tfsdk.Config{},
	}

	resp := r.DeleteResponse{
		State: tfsdk.State{
			Schema: schemaResponse.Schema,
		},
		Diagnostics: nil,
	}

	start := time.Now()
	sleepResource.Delete(context.Background(), req, &resp)
	end := time.Now()
	elapsed := end.Sub(start)

	if elapsed < expectedDuration {
		t.Errorf("did not sleep long enough, expected duration: %d got: %d", expectedDuration, elapsed)
	}
}

func TestAccTimeSleep_CreateDuration(t *testing.T) {
	resourceName := "time_sleep.test"

	// These ID comparisons can eventually be replaced by the multiple value checks once released
	// in terraform-plugin-testing: https://github.com/hashicorp/terraform-plugin-testing/issues/295
	captureTimeState1 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))
	captureTimeState2 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepCreateDuration("1ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.StringExact("1ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					captureTimeState1,
				},
			},
			// This test may work in local execution but typically does not work in CI because of its reliance
			// on the current time stamp in the ID. We will also need to revisit this test later once TF core allows
			// multiple parameters in Import
			//{
			//	ResourceName:      resourceName,
			//	ImportState:       true,
			//	ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
			//	ImportStateVerify: true,
			//},
			{
				Config: testAccConfigTimeSleepCreateDuration("2ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.StringExact("2ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					captureTimeState2,
				},
			},
		},
	})

	// Ensure the id time value is different due to the sleep
	if captureTimeState1.Value == captureTimeState2.Value {
		t.Fatal("attribute values are the same")
	}
}

func TestAccTimeSleep_DestroyDuration(t *testing.T) {
	resourceName := "time_sleep.test"

	// These ID comparisons can eventually be replaced by the multiple value checks once released
	// in terraform-plugin-testing: https://github.com/hashicorp/terraform-plugin-testing/issues/295
	captureTimeState1 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))
	captureTimeState2 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepDestroyDuration("1ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("destroy_duration"), knownvalue.StringExact("1ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					captureTimeState1,
				},
			},
			// This test may work in local execution but typically does not work in CI because of its reliance
			// on the current time stamp in the ID. We will also need to revisit this test later once TF core allows
			// multiple parameters in Import
			//{
			//	ResourceName:      resourceName,
			//	ImportState:       true,
			//	ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
			//	ImportStateVerify: true,
			//},
			{
				Config: testAccConfigTimeSleepDestroyDuration("2ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("destroy_duration"), knownvalue.StringExact("2ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					captureTimeState2,
				},
			},
		},
	})

	// Ensure the id time value is different due to the sleep
	if captureTimeState1.Value == captureTimeState2.Value {
		t.Fatal("attribute values are the same")
	}
}

func TestAccTimeSleep_Triggers(t *testing.T) {
	resourceName := "time_sleep.test"

	// These ID comparisons can eventually be replaced by the multiple value checks once released
	// in terraform-plugin-testing: https://github.com/hashicorp/terraform-plugin-testing/issues/295
	captureTimeState1 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))
	captureTimeState2 := timetesting.NewExtractState(resourceName, tfjsonpath.New("id"))

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.NotNull()),
					captureTimeState1,
				},
			},
			// This test may work in local execution but typically does not work in CI because of its reliance
			// on the current time stamp in the ID. We will also need to revisit this test later once TF core allows
			// multiple parameters in Import
			//{
			//	ResourceName:            resourceName,
			//	ImportState:             true,
			//	ImportStateIdFunc:       testAccTimeSleepImportStateIdFunc(resourceName),
			//	ImportStateVerify:       true,
			//	ImportStateVerifyIgnore: []string{"triggers"},
			//},
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1updated"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1updated")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.NotNull()),
					captureTimeState2,
				},
			},
		},
	})

	// Ensure the id time value is different due to the sleep
	if captureTimeState1.Value == captureTimeState2.Value {
		t.Fatal("attribute values are the same")
	}
}

func TestAccTimeSleep_Upgrade(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeSleepCreateDuration("1ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.StringExact("1ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeSleepCreateDuration("1ms"),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeSleepCreateDuration("1ms"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.StringExact("1ms")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccTimeSleep_Validators(t *testing.T) {

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(`resource "time_sleep" "test" {
                     triggers = {
						%[1]q = %[2]q
					  }
                  }`, "key1", "value1"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				Config:      testAccConfigTimeSleepCreateDuration("1"),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Value Match`),
			},
		},
	})
}

// Validate that importing works as expected
func TestResourceTimeSleepImport(t *testing.T) {

	// want := timeSleepModelV0{
	// 	CreateDuration:  types.StringValue("1s"),
	// 	DestroyDuration: types.StringValue("1s"),
	// 	Triggers:        types.MapValueMust(types.StringType, triggers),
	// 	ID:              timetypes.NewRFC3339TimeValue(time.Now().UTC()),
	// }
	var tests = []struct {
		name    string
		id      string
		want    timeSleepModelV0
		error   bool
		summary string
	}{
		{
			name: "just_createduration",
			id:   "1s,",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringNull(),
				Triggers:        types.MapValueMust(types.StringType, map[string]attr.Value{}),
				ID:              timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "create_and_destroy_duration_only",
			id:   "1s,20s",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringValue("20s"),
				Triggers:        types.MapValueMust(types.StringType, map[string]attr.Value{}),
				ID:              timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "destroy_duration_only",
			id:   ",20s",
			want: timeSleepModelV0{
				CreateDuration:  types.StringNull(),
				DestroyDuration: types.StringValue("20s"),
				Triggers:        types.MapValueMust(types.StringType, map[string]attr.Value{}),
				ID:              timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "create_destroy_single_trigger_only",
			id:   "1s,1s,test=testvalue",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringValue("1s"),
				Triggers: types.MapValueMust(types.StringType, map[string]attr.Value{
					"test": types.StringValue("testvalue"),
				}),
				ID: timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "create_destroy_multi",
			id:   "1s,1s,test1=testvalue1,test2=testvalue2",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringValue("1s"),
				Triggers: types.MapValueMust(types.StringType, map[string]attr.Value{
					"test1": types.StringValue("testvalue1"),
					"test2": types.StringValue("testvalue2"),
				}),
				ID: timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "create_no_destroy_single_only",
			id:   "1s,,test1=testvalue1",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringNull(),
				Triggers: types.MapValueMust(types.StringType, map[string]attr.Value{
					"test1": types.StringValue("testvalue1"),
				}),
				ID: timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "no_create_destroy_single_only",
			id:   ",1s,test1=testvalue1",
			want: timeSleepModelV0{
				CreateDuration:  types.StringNull(),
				DestroyDuration: types.StringValue("1s"),
				Triggers: types.MapValueMust(types.StringType, map[string]attr.Value{
					"test1": types.StringValue("testvalue1"),
				}),
				ID: timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error: false,
		},
		{
			name: "create_destroy_invalid_trigger_format",
			id:   "1s,1s,test=testvalue,test2==testvalue2",
			want: timeSleepModelV0{
				CreateDuration:  types.StringValue("1s"),
				DestroyDuration: types.StringValue("1s"),
				Triggers: types.MapValueMust(types.StringType, map[string]attr.Value{
					"test": types.StringValue("testvalue"),
				}),
				ID: timetypes.NewRFC3339TimeValue(time.Now().UTC()),
			},
			error:   true,
			summary: "Trigger import error",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			sleepResource := NewTimeSleepResource()

			schemaResponse := r.SchemaResponse{}
			sleepResource.Schema(context.Background(), r.SchemaRequest{}, &schemaResponse)

			req := r.ImportStateRequest{
				ID: test.id,
			}

			resp := r.ImportStateResponse{
				State: tfsdk.State{
					Schema: schemaResponse.Schema,
				},
				Diagnostics: nil,
			}

			sleepResource.(r.ResourceWithImportState).ImportState(context.Background(), req, &resp)
			if resp.Diagnostics.HasError() {
				if !test.error {
					t.Fatalf("Diags was not empty: %+v", resp.Diagnostics)
				}

				for _, diag := range resp.Diagnostics {
					if diag.Summary() != test.summary {
						t.Fatalf("Diags had additional errors that were not expected: %+v", resp.Diagnostics)
					}
				}
				// Passed the test
				return
			}

			state := timeSleepModelV0{}
			resp.State.Get(context.Background(), &state)

			// This is a bit of a cheat, because the ID of the state object is generated inside the function
			// so we need to set the `want` structure to have the right time.
			test.want.ID = state.ID

			if !reflect.DeepEqual(state, test.want) {
				t.Fatal(fmt.Sprintf("Trigger map was not what we expected, want `%+v` got `%+v`", test.want, state))
			}
		})
	}
}

func TestAccTimeSleepImport_Triggers1(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepImportTriggers1("key1", "value1"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(1)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.NotNull()),
				},
			},
		},
	})
}

func TestAccTimeSleepImport_Triggers2(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepImportTriggers2("key1", "value1", "key2", "value2"),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers"), knownvalue.MapSizeExact(2)),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key1"), knownvalue.StringExact("value1")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("triggers").AtMapKey("key2"), knownvalue.StringExact("value2")),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("id"), knownvalue.NotNull()),
					statecheck.ExpectKnownValue(resourceName, tfjsonpath.New("create_duration"), knownvalue.NotNull()),
				},
			},
		},
	})
}

//func testAccTimeSleepImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
//	return func(s *terraform.State) (string, error) {
//		rs, ok := s.RootModule().Resources[resourceName]
//		if !ok {
//			return "", fmt.Errorf("Not found: %s", resourceName)
//		}
//
//		createDuration := rs.Primary.Attributes["create_duration"]
//		destroyDuration := rs.Primary.Attributes["destroy_duration"]
//
//		return fmt.Sprintf("%s,%s", createDuration, destroyDuration), nil
//	}
//}

func testAccConfigTimeSleepCreateDuration(createDuration string) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  create_duration = %[1]q
}
`, createDuration)
}

func testAccConfigTimeSleepDestroyDuration(destroyDuration string) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  destroy_duration = %[1]q
}
`, destroyDuration)
}

func testAccConfigTimeSleepTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  create_duration = "1s"

  triggers = {
    %[1]q = %[2]q
  }
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeSleepImportTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
import {
  to = time_sleep.test
  id = "1s,,%s=%s"
}

resource "time_sleep" "test" {
  create_duration = "1s"

  triggers = {
    %[1]q = %[2]q
  }
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeSleepImportTriggers2(keeperKey1 string, keeperKey2 string, keeperKey3 string, keeperKey4 string) string {
	return fmt.Sprintf(`
import {
  to = time_sleep.test
  id = "1s,,%s=%s,%s=%s"
}

resource "time_sleep" "test" {
  create_duration = "1s"

  triggers = {
    %[1]q = %[2]q
    %[3]q = %[4]q
  }
}
`, keeperKey1, keeperKey2, keeperKey3, keeperKey4)
}
