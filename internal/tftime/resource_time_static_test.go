package tftime

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

func TestAccTimeStatic_basic(t *testing.T) {
	resourceName := "time_static.test"

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestCheckNoResourceAttr(resourceName, "expiration_rfc3339"),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc822", regexp.MustCompile(`^\d{2} \w+ \d{2} \d{2}:\d{2} \w+$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc822z", regexp.MustCompile(`^\d{2} \w+ \d{2} \d{2}:\d{2} \+\d{4}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc850", regexp.MustCompile(`^\w+, \d{2}-\w+-\d{2} \d{2}:\d{2}:\d{2} \w+$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc1123", regexp.MustCompile(`^\w+, \d{2} \w+ \d{4} \d{2}:\d{2}:\d{2} \w+$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc1123z", regexp.MustCompile(`^\w+, \d{2} \w+ \d{4} \d{2}:\d{2}:\d{2} \+\d{4}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "unixdate", regexp.MustCompile(`^\w+ \w+\s+\d{1,2} \d{2}:\d{2}:\d{2} \w+ \d{4}$`)),
					resource.TestMatchResourceAttr(resourceName, "year", regexp.MustCompile(`^\d{4}$`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationDays_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationDays(timestamp.Format(time.RFC3339), 7),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_days", "7"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(0, 0, 7).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationDays_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().AddDate(0, 0, -2)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationDays(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_days", "1"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(0, 0, 1).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationHours_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationHours(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_hours", "3"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.Add(3*time.Hour).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationHours_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().Add(-2 * time.Hour)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationHours(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_hours", "1"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.Add(1*time.Hour).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationMinutes_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationMinutes(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_minutes", "3"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.Add(3*time.Minute).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationMinutes_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().Add(-2 * time.Minute)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationMinutes(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_minutes", "1"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.Add(1*time.Minute).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationMonths_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationMonths(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_months", "3"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(0, 3, 0).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationMonths_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().AddDate(0, -2, 0)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationMonths(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_months", "1"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(0, 1, 0).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationRfc3339_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()
	expirationTimestamp := time.Now().UTC().AddDate(0, 0, 7)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationRfc3339(timestamp.Format(time.RFC3339), expirationTimestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", expirationTimestamp.Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationRfc3339_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().AddDate(0, 0, -2)
	expirationTimestamp := time.Now().UTC().AddDate(0, 0, -1)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationRfc3339(timestamp.Format(time.RFC3339), expirationTimestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", expirationTimestamp.Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_ExpirationYears_basic(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationYears(timestamp.Format(time.RFC3339), 3),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_years", "3"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(3, 0, 0).Format(time.RFC3339)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
			},
		},
	})
}

func TestAccTimeStatic_ExpirationYears_expired(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC().AddDate(-2, 0, 0)

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticExpirationYears(timestamp.Format(time.RFC3339), 1),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "expiration_years", "1"),
					resource.TestCheckResourceAttr(resourceName, "expiration_rfc3339", timestamp.AddDate(1, 0, 0).Format(time.RFC3339)),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func TestAccTimeStatic_Rfc3339(t *testing.T) {
	resourceName := "time_static.test"
	timestamp := time.Now().UTC()
	day := strconv.Itoa(timestamp.Day())
	hour := strconv.Itoa(timestamp.Hour())
	minute := strconv.Itoa(timestamp.Minute())
	month := strconv.Itoa(int(timestamp.Month()))
	second := strconv.Itoa(timestamp.Second())
	unix := strconv.Itoa(int(timestamp.Unix()))
	year := strconv.Itoa(timestamp.Year())

	resource.UnitTest(t, resource.TestCase{
		Providers:    testAccProviders,
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticRfc3339(timestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "day", day),
					resource.TestCheckNoResourceAttr(resourceName, "expiration_rfc3339"),
					resource.TestCheckResourceAttr(resourceName, "hour", hour),
					resource.TestCheckResourceAttr(resourceName, "minute", minute),
					resource.TestCheckResourceAttr(resourceName, "month", month),
					resource.TestCheckResourceAttr(resourceName, "rfc822", timestamp.Format(time.RFC822)),
					resource.TestCheckResourceAttr(resourceName, "rfc822z", timestamp.Format(time.RFC822Z)),
					resource.TestCheckResourceAttr(resourceName, "rfc850", timestamp.Format(time.RFC850)),
					resource.TestCheckResourceAttr(resourceName, "rfc1123", timestamp.Format(time.RFC1123)),
					resource.TestCheckResourceAttr(resourceName, "rfc1123z", timestamp.Format(time.RFC1123Z)),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", second),
					resource.TestCheckResourceAttr(resourceName, "unix", unix),
					resource.TestCheckResourceAttr(resourceName, "unixdate", timestamp.Format(time.UnixDate)),
					resource.TestCheckResourceAttr(resourceName, "year", year),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccConfigTimeStatic() string {
	return fmt.Sprintf(`
resource "time_static" "test" {}
`)
}

func testAccConfigTimeStaticExpirationDays(rfc3339 string, expirationDays int) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_days = %[2]d
  rfc3339         = %[1]q
}
`, rfc3339, expirationDays)
}

func testAccConfigTimeStaticExpirationHours(rfc3339 string, expirationHours int) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_hours = %[2]d
  rfc3339          = %[1]q
}
`, rfc3339, expirationHours)
}

func testAccConfigTimeStaticExpirationMinutes(rfc3339 string, expirationMinutes int) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_minutes = %[2]d
  rfc3339            = %[1]q
}
`, rfc3339, expirationMinutes)
}

func testAccConfigTimeStaticExpirationMonths(rfc3339 string, expirationMonths int) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_months = %[2]d
  rfc3339           = %[1]q
}
`, rfc3339, expirationMonths)
}

func testAccConfigTimeStaticExpirationYears(rfc3339 string, expirationYears int) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_years = %[2]d
  rfc3339          = %[1]q
}
`, rfc3339, expirationYears)
}

func testAccConfigTimeStaticExpirationRfc3339(rfc3339 string, expirationRfc3339 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  expiration_rfc3339 = %[2]q
  rfc3339            = %[1]q
}
`, rfc3339, expirationRfc3339)
}

func testAccConfigTimeStaticRfc3339(rfc3339 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
  rfc3339 = %[1]q
}
`, rfc3339)
}
