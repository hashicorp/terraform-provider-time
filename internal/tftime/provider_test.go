package tftime

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
	"github.com/hashicorp/terraform-provider-time/internal/provider"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//nolint:unparam
var testAccProviderFactories = map[string]func() (tfprotov5.ProviderServer, error){}

func init() {
	testAccProviderFactories["time"] = func() (tfprotov5.ProviderServer, error) {
		ctx := context.Background()

		// SDKv2 provider
		sdkv2 := Provider().GRPCProvider

		// Framework provider
		framework := providerserver.NewProtocol5(provider.New())

		muxServer, err := tf5muxserver.NewMuxServer(ctx, sdkv2, framework)
		if err != nil {
			return nil, err
		}

		return muxServer.ProviderServer(), nil
	}
}

func testCheckAttributeValuesDiffer(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) == testStringValue(j) {
			return fmt.Errorf("attribute values are the same")
		}

		return nil
	}
}

func testCheckAttributeValuesSame(i *string, j *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if testStringValue(i) != testStringValue(j) {
			return fmt.Errorf("attribute values are different")
		}

		return nil
	}
}

func testExtractResourceAttr(resourceName string, attributeName string, attributeValue *string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[resourceName]

		if !ok {
			return fmt.Errorf("resource name %s not found in state", resourceName)
		}

		attrValue, ok := rs.Primary.Attributes[attributeName]

		if !ok {
			return fmt.Errorf("attribute %s not found in resource %s state", attributeName, resourceName)
		}

		*attributeValue = attrValue

		return nil
	}
}

// Certain testing requires time differences that are too fast for unit testing.
// Sleeping for a second or two seems pragmatic in our testing.
func testSleep(seconds int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		time.Sleep(time.Duration(seconds) * time.Second)

		return nil
	}
}

func testStringValue(sPtr *string) string {
	if sPtr == nil {
		return ""
	}

	return *sPtr
}
