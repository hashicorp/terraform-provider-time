package tftime

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccTimeSleep_CreateSeconds(t *testing.T) {
	var time1, time2 string
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepCreateSeconds(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_seconds", "1"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeSleepCreateSeconds(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "create_seconds", "2"),
					testExtractResourceAttr(resourceName, "id", &time2),
					testCheckAttributeValuesSame(&time1, &time2),
				),
			},
		},
	})
}

func TestAccTimeSleep_DestroySeconds(t *testing.T) {
	var time1, time2 string
	resourceName := "time_sleep.test"

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepDestroySeconds(1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_seconds", "1"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeSleepImportStateIdFunc(resourceName),
				ImportStateVerify: true,
			},
			{
				Config: testAccConfigTimeSleepDestroySeconds(2),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "destroy_seconds", "2"),
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
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeSleepTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					testExtractResourceAttr(resourceName, "id", &time1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccTimeSleepImportStateIdFunc(resourceName),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
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

		createSeconds := rs.Primary.Attributes["create_seconds"]
		destroySeconds := rs.Primary.Attributes["destroy_seconds"]

		return fmt.Sprintf("%s,%s", createSeconds, destroySeconds), nil
	}
}

func testAccConfigTimeSleepCreateSeconds(createSeconds int) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  create_seconds = %[1]d
}
`, createSeconds)
}

func testAccConfigTimeSleepDestroySeconds(destroySeconds int) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  destroy_seconds = %[1]d
}
`, destroySeconds)
}

func testAccConfigTimeSleepTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_sleep" "test" {
  create_seconds = 1

  triggers = {
    %[1]q = %[2]q
  }
}
`, keeperKey1, keeperKey2)
}
