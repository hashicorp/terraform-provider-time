// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	r "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
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
	var time1, time2 string
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepCreateDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "2ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testExtractResourceAttr(resourceName, "id", &time2),
					testCheckAttributeValuesSame(&time1, &time2),
				),
			},
		},
	})
}

func TestAccTimeSleep_DestroyDuration(t *testing.T) {
	var time1, time2 string
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepDestroyDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_duration", "2ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					testExtractResourceAttr(resourceName, "id", &time2),
					testCheckAttributeValuesSame(&time1, &time2),
				),
			},
		},
	})
}

func TestAccTimeSleep_Triggers(t *testing.T) {
	var time1, time2 string
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_duration"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
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
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_duration"),
					testExtractResourceAttr(resourceName, "id", &time2),
					testCheckAttributeValuesDiffer(&time1, &time2),
				),
			},
		},
	})
}

func TestAccTimeSleep_Upgrade(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeSleepCreateDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeSleepCreateDuration("1ms"),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeSleepCreateDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
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
