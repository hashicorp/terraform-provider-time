package tftime

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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
	d.Set("create_duration", durationStr)

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
	d.Set("destroy_duration", durationStr)

	start := time.Now()
	resourceTimeSleepDelete(context.Background(), d, nil)
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepCreateDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "1ms"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
			// 	ImportStateVerify: true,
			// },
			{
				Config: testAccConfigTimeSleepCreateDuration("2ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_duration", "2ms"),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepDestroyDuration("1ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_duration", "1ms"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			// {
			// 	ResourceName:      resourceName,
			// 	ImportState:       true,
			// 	ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
			// 	ImportStateVerify: true,
			// },
			{
				Config: testAccConfigTimeSleepDestroyDuration("2ms"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_duration", "2ms"),
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
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			// {
			// 	ResourceName:            resourceName,
			// 	ImportState:             true,
			// 	ImportStateIdFunc:       testAccTimeSleepImportStateIdFunc(resourceName),
			// 	ImportStateVerify:       true,
			// 	ImportStateVerifyIgnore: []string{"triggers"},
			// },
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					testExtractResourceAttr(resourceName, "id", &time2),
					testCheckAttributeValuesDiffer(&time1, &time2),
				),
			},
		},
	})
}

func testAccTimeSleepImportStateIdFunc(resourceName string) resource.ImportStateIdFunc {
	return func(s *terraform.State) (string, error) {
		rs, ok := s.RootModule().Resources[resourceName]
		if !ok {
			return "", fmt.Errorf("Not found: %s", resourceName)
		}

		createDuration := rs.Primary.Attributes["create_duration"]
		destroyDuration := rs.Primary.Attributes["destroy_duration"]

		return fmt.Sprintf("%s,%s", createDuration, destroyDuration), nil
	}
}

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
