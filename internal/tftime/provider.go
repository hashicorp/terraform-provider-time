package tftime

import (
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		ResourcesMap: map[string]*schema.Resource{
			"time_rotating": resourceTimeRotating(),
			"time_static":   resourceTimeStatic(),
		},
	}
}
