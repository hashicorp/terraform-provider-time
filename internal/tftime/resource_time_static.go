package tftime

import (
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTimeStatic() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a static time resource, which keeps a locally sourced UTC timestamp stored in the Terraform state. " +
			"This prevents perpetual differences caused by using " +
			"the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html).",
		Create: resourceTimeStaticCreate,
		Read:   resourceTimeStaticRead,
		Delete: schema.Noop,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"day": {
				Description: "Number day of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved. " +
					"See [the main provider documentation](../index.md) for more information.",
				Type:     schema.TypeMap,
				Optional: true,
				ForceNew: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"minute": {
				Description: "Number minute of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"month": {
				Description: "Number month of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"rfc3339": {
				Description:  "UTC RFC3339 format of timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"second": {
				Description: "Number second of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"unix": {
				Description: "Number of seconds since epoch time, e.g. `1581489373`.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"year": {
				Description: "Number year of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
		},
	}
}

func resourceTimeStaticCreate(d *schema.ResourceData, m interface{}) error {
	timestamp := time.Now().UTC()

	if v, ok := d.GetOk("rfc3339"); ok {
		var err error
		timestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing rfc3339 (%s): %s", v.(string), err)
		}
	}

	d.SetId(timestamp.Format(time.RFC3339))

	return resourceTimeStaticRead(d, m)
}

func resourceTimeStaticRead(d *schema.ResourceData, m interface{}) error {
	timestamp, err := time.Parse(time.RFC3339, d.Id())

	if err != nil {
		return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
	}

	if err := d.Set("day", timestamp.Day()); err != nil {
		return fmt.Errorf("error setting day: %s", err)
	}

	if err := d.Set("hour", timestamp.Hour()); err != nil {
		return fmt.Errorf("error setting hour: %s", err)
	}

	if err := d.Set("minute", timestamp.Minute()); err != nil {
		return fmt.Errorf("error setting minute: %s", err)
	}

	if err := d.Set("month", int(timestamp.Month())); err != nil {
		return fmt.Errorf("error setting month: %s", err)
	}

	if err := d.Set("rfc3339", timestamp.Format(time.RFC3339)); err != nil {
		return fmt.Errorf("error setting rfc3339: %s", err)
	}

	if err := d.Set("second", timestamp.Second()); err != nil {
		return fmt.Errorf("error setting second: %s", err)
	}

	if err := d.Set("unix", timestamp.Unix()); err != nil {
		return fmt.Errorf("error setting unix: %s", err)
	}

	if err := d.Set("year", timestamp.Year()); err != nil {
		return fmt.Errorf("error setting year: %s", err)
	}

	return nil
}
