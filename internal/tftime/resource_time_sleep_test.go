package tftime

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// Since the acceptance testing framework can introduce uncontrollable time delays,
// verify that sleeping works as expected via unit testing.
func TestResourceTimeSleepCreate(t *testing.T) {
	durationStr := "1s"
	expectedDuration, err := time.ParseDuration("1s")

	if err != nil {
		t.Fatalf("unable to parse test duration: %s", err)
	}

	d := resourceTimeSleep().Data(nil)
	d.SetType("time_sleep")
	d.SetId("test")
	err = d.Set("create_duration", durationStr)
	if err != nil {
		t.Fatalf("unable set create_duration to %s with error: %s", durationStr, err)
	}

	start := time.Now()
	resourceTimeSleepCreate(context.Background(), d, nil)
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

	d := resourceTimeSleep().Data(nil)
	d.SetType("time_sleep")
	d.SetId("test")
	err = d.Set("destroy_duration", durationStr)
	if err != nil {
		t.Fatalf("unable set destroy_duration to %s with error: %s", durationStr, err)
	}

	start := time.Now()
	resourceTimeSleepDelete(context.Background(), d, nil)
	end := time.Now()
	elapsed := end.Sub(start)

	if elapsed < expectedDuration {
		t.Errorf("did not sleep long enough, expected duration: %d got: %d", expectedDuration, elapsed)
	}
}

func TestAccTimeSleep_CreateDuration(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepCreateDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				),
			},
		},
	})
}

func TestAccTimeSleep_DestroyDuration(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepDestroyDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_duration", "1ms"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				),
			},
		},
	})
}

func TestAccTimeSleep_Triggers(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: testAccProviderFactories,
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "create_duration"),
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
				),
			},
		},
	})
}

func TestAccTimeSleep_Upgrade(t *testing.T) {
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
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
				ProtoV5ProviderFactories: testAccProviderFactories,
				Config:                   testAccConfigTimeSleepCreateDuration("1ms"),
				PlanOnly:                 true,
			},
			{
				ProtoV5ProviderFactories: testAccProviderFactories,
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
		ProtoV5ProviderFactories: testAccProviderFactories,
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
