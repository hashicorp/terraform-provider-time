package tftime

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceTimeSleep() *schema.Resource {
	return &schema.Resource{
		Create: resourceTimeSleepCreate,
		Read:   schema.Noop,
		Update: schema.Noop,
		Delete: resourceTimeSleepDelete,

		Importer: &schema.ResourceImporter{
			State: func(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 2 || (idParts[0] == "" && idParts[1] == "") {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected CREATESECONDS,DESTROYSECONDS where at least one value is non-empty", d.Id())
				}

				createSeconds, _ := strconv.Atoi(idParts[0])
				destroySeconds, _ := strconv.Atoi(idParts[1])

				if createSeconds > 0 {
					if err := d.Set("create_seconds", createSeconds); err != nil {
						return nil, fmt.Errorf("error setting create_seconds: %s", err)
					}
				}

				if destroySeconds > 0 {
					if err := d.Set("destroy_seconds", destroySeconds); err != nil {
						return nil, fmt.Errorf("error setting destroy_seconds: %s", err)
					}
				}

				d.SetId(time.Now().UTC().Format(time.RFC3339))

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"create_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"create_seconds",
					"destroy_seconds",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"destroy_seconds": {
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"create_seconds",
					"destroy_seconds",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"triggers": {
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceTimeSleepCreate(d *schema.ResourceData, m interface{}) error {
	if v, ok := d.GetOk("create_seconds"); ok {
		time.Sleep(time.Duration(v.(int)) * time.Second)
	}

	d.SetId(time.Now().UTC().Format(time.RFC3339))

	return nil
}

func resourceTimeSleepDelete(d *schema.ResourceData, m interface{}) error {
	if v, ok := d.GetOk("destroy_seconds"); ok {
		time.Sleep(time.Duration(v.(int)) * time.Second)
	}

	return nil
}
