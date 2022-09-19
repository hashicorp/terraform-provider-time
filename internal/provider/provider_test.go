package provider

import (
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

//nolint:unparam
func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"time": providerserver.NewProtocol5WithError(New()),
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
