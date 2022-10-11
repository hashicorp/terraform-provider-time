package provider

import (
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTimeStatic_basic(t *testing.T) {
	resourceName := "time_static.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
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

func TestAccTimeStatic_Triggers(t *testing.T) {
	var time1, time2 string
	resourceName := "time_static.test"

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
					testExtractResourceAttr(resourceName, "rfc3339", &time1),
					testSleep(1),
				),
			},
			{
				ResourceName:            resourceName,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{"triggers"},
			},
			{
				Config: testAccConfigTimeStaticTriggers1("key1", "value1updated"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "triggers.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "triggers.key1", "value1updated"),
					resource.TestCheckResourceAttrSet(resourceName, "rfc3339"),
					testExtractResourceAttr(resourceName, "rfc3339", &time2),
					testCheckAttributeValuesDiffer(&time1, &time2),
				),
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
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config: testAccConfigTimeStaticRfc3339(timestamp.Format(time.RFC3339)),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "day", day),
					resource.TestCheckResourceAttr(resourceName, "hour", hour),
					resource.TestCheckResourceAttr(resourceName, "minute", minute),
					resource.TestCheckResourceAttr(resourceName, "month", month),
					resource.TestCheckResourceAttr(resourceName, "rfc3339", timestamp.Format(time.RFC3339)),
					resource.TestCheckResourceAttr(resourceName, "second", second),
					resource.TestCheckResourceAttr(resourceName, "unix", unix),
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

func TestAccTimeStatic_Upgrade(t *testing.T) {
	resourceName := "time_static.test"

	resource.Test(t, resource.TestCase{
		CheckDestroy: nil,
		Steps: []resource.TestStep{
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeStatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "year", regexp.MustCompile(`^\d{4}$`)),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigTimeStatic(),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion080(),
				Config:            testAccConfigTimeStatic(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestMatchResourceAttr(resourceName, "day", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "hour", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "minute", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "month", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "rfc3339", regexp.MustCompile(`^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}Z$`)),
					resource.TestMatchResourceAttr(resourceName, "second", regexp.MustCompile(`^\d{1,2}$`)),
					resource.TestMatchResourceAttr(resourceName, "unix", regexp.MustCompile(`^\d+$`)),
					resource.TestMatchResourceAttr(resourceName, "year", regexp.MustCompile(`^\d{4}$`)),
				),
			},
		},
	})
}

func TestAccTimeStatic_Validators(t *testing.T) {
	timestamp := time.Now().UTC()

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []resource.TestStep{
			{
				Config:      testAccConfigTimeStaticRfc3339(timestamp.Format(time.RFC850)),
				ExpectError: regexp.MustCompile(`.*Value must be a string in RFC3339 format`),
			},
		},
	})
}

func testAccConfigTimeStatic() string {
	return `
resource "time_static" "test" {}
`
}

func testAccConfigTimeStaticTriggers1(keeperKey1 string, keeperKey2 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
 triggers = {
   %[1]q = %[2]q
 }
}
`, keeperKey1, keeperKey2)
}

func testAccConfigTimeStaticRfc3339(rfc3339 string) string {
	return fmt.Sprintf(`
resource "time_static" "test" {
 rfc3339 = %[1]q
}
`, rfc3339)
}
