package tftime

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccTimeRotating_Triggers(t *testing.T) {
	var time1, time2 string
	resourceName := "time_rotating.test"

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					testExtractResourceAttr(resourceName, "rfc3339", &time1),
					testSleep(1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateIdFunc:       testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccConfigTimeRotatingTriggers1("key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					testExtractResourceAttr(resourceName, "rfc3339", &time2),
					testCheckAttributeValuesDiffer(&time1, &time2),
				),
			},
		},
	})
}

func TestAccTimeRotating_RotationDays_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationDays(timestamp.Format(time.RFC3339), 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(0, 0, 7).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationDays_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().AddDate(0, 0, -2)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationDays(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(0, 0, 1).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationHours_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationHours(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_hours", "3"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.Add(3*time.Hour).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationHours_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().Add(-2 * time.Hour)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationHours(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_hours", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.Add(1*time.Hour).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationMinutes_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationMinutes(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_minutes", "3"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.Add(3*time.Minute).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationMinutes_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().Add(-2 * time.Minute)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationMinutes(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_minutes", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.Add(1*time.Minute).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationMonths_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationMonths(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_months", "3"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(0, 3, 0).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationMonths_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().AddDate(0, -2, 0)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationMonths(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_months", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(0, 1, 0).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationRfc3339_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()
	rotationTimestamp := time.Now().UTC().AddDate(0, 0, 7)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationRfc3339(timestamp.Format(time.RFC3339), rotationTimestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationRfc3339_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().AddDate(0, 0, -2)
	rotationTimestamp := time.Now().UTC().AddDate(0, 0, -1)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationRfc3339(timestamp.Format(time.RFC3339), rotationTimestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationYears_basic(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationYears(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(3, 0, 0).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateIdFunc: testAccTimeRotatingImportStateIdFunc(),
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeRotating_RotationYears_expired(t *testing.T) {
	resourceName := "time_rotating.test"
	timestamp := time.Now().UTC().AddDate(-2, 0, 0)

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeRotatingRotationYears(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "rotation_years", "1"),
					resource.TestCheckResourceAttr(resourceName, "rotation_rfc3339", timestamp.AddDate(1, 0, 0).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
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

func testAccConfigTimeRotatingRotationDays(rfc3339 string, rotationDays int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_days = %[2]d
  rfc3339       = %[1]q
}
`, rfc3339, rotationDays)
}

func testAccConfigTimeRotatingRotationHours(rfc3339 string, rotationHours int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_hours = %[2]d
  rfc3339        = %[1]q
}
`, rfc3339, rotationHours)
}

func testAccConfigTimeRotatingRotationMinutes(rfc3339 string, rotationMinutes int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_minutes = %[2]d
  rfc3339          = %[1]q
}
`, rfc3339, rotationMinutes)
}

func testAccConfigTimeRotatingRotationMonths(rfc3339 string, rotationMonths int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_months = %[2]d
  rfc3339         = %[1]q
}
`, rfc3339, rotationMonths)
}

func testAccConfigTimeRotatingRotationYears(rfc3339 string, rotationYears int) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_years = %[2]d
  rfc3339        = %[1]q
}
`, rfc3339, rotationYears)
}

func testAccConfigTimeRotatingRotationRfc3339(rfc3339 string, rotationRfc3339 string) string {
	return fmt.Sprintf(`
resource "time_rotating" "test" {
  rotation_rfc3339 = %[2]q
  rfc3339          = %[1]q
}
`, rfc3339, rotationRfc3339)
}
