package tftime

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/helper/validation"
)

func resourceTimeStatic() *schema.Resource {
	return &schema.Resource{
		Create: resourceTimeStaticCreate,
		Read:   resourceTimeStaticRead,
		Update: resourceTimeStaticUpdate,
		Delete: schema.Noop,

		CustomizeDiff: customdiff.Sequence(
			customdiff.If(resourceTimeStaticConditionExpirationChange,
				func(diff *schema.ResourceDiff, meta interface{}) error {
					if diff.Id() == "" {
						return nil
					}

					timestamp, err := time.Parse(time.RFC3339, diff.Id())

					if err != nil {
						return fmt.Errorf("error parsing timestamp (%s): %s", diff.Id(), err)
					}

					var expirationTimestamp *time.Time

					if v, ok := diff.GetOk("expiration_days"); ok {
						expirationTimestamp = timePtr(timestamp.AddDate(0, 0, v.(int)))
					}

					if v, ok := diff.GetOk("expiration_hours"); ok {
						expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Hour))
					}

					if v, ok := diff.GetOk("expiration_minutes"); ok {
						expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Minute))
					}

					if v, ok := diff.GetOk("expiration_months"); ok {
						expirationTimestamp = timePtr(timestamp.AddDate(0, v.(int), 0))
					}

					if v, ok := diff.GetOk("expiration_years"); ok {
						expirationTimestamp = timePtr(timestamp.AddDate(v.(int), 0, 0))
					}

					if expirationTimestamp != nil {
						if err := diff.SetNew("expiration_rfc3339", expirationTimestamp.Format(time.RFC3339)); err != nil {
							return fmt.Errorf("error setting new expiration_rfc3339: %s", err)
						}
					}

					return nil
				},
			),
			customdiff.ForceNewIf("expiration_rfc3339", func(diff *schema.ResourceDiff, meta interface{}) bool {
				now := time.Now().UTC()
				expirationTimestamp, err := time.Parse(time.RFC3339, diff.Get("expiration_rfc3339").(string))

				if err != nil {
					return false
				}

				return now.After(expirationTimestamp)
			}),
		),

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Schema: map[string]*schema.Schema{
			"day": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"expiration_days": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"expiration_hours": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"expiration_days",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"expiration_minutes": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_months",
					"expiration_rfc3339",
					"expiration_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"expiration_months": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_rfc3339",
					"expiration_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"expiration_rfc3339": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ConflictsWith: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_years",
				},
				ValidateFunc: validation.IsRFC3339Time,
			},
			"expiration_years": {
				Type:     schema.TypeInt,
				Optional: true,
				ConflictsWith: []string{
					"expiration_days",
					"expiration_hours",
					"expiration_minutes",
					"expiration_months",
					"expiration_rfc3339",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"hour": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"minute": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"month": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"rfc822": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rfc822z": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rfc850": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rfc1123": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rfc1123z": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rfc3339": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ForceNew:     true,
				ValidateFunc: validation.IsRFC3339Time,
			},
			"second": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unix": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"unixdate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"year": {
				Type:     schema.TypeInt,
				Computed: true,
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

	var expirationTimestamp *time.Time

	if v, ok := d.GetOk("expiration_days"); ok {
		expirationTimestamp = timePtr(timestamp.AddDate(0, 0, v.(int)))
	}

	if v, ok := d.GetOk("expiration_hours"); ok {
		expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Hour))
	}

	if v, ok := d.GetOk("expiration_minutes"); ok {
		expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Minute))
	}

	if v, ok := d.GetOk("expiration_months"); ok {
		expirationTimestamp = timePtr(timestamp.AddDate(0, v.(int), 0))
	}

	if v, ok := d.GetOk("expiration_years"); ok {
		expirationTimestamp = timePtr(timestamp.AddDate(v.(int), 0, 0))
	}

	if expirationTimestamp != nil {
		if err := d.Set("expiration_rfc3339", expirationTimestamp.Format(time.RFC3339)); err != nil {
			return fmt.Errorf("error setting expiration_rfc3339: %s", err)
		}
	}

	return resourceTimeStaticRead(d, m)
}

func resourceTimeStaticRead(d *schema.ResourceData, m interface{}) error {
	timestamp, err := time.Parse(time.RFC3339, d.Id())

	if err != nil {
		return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
	}

	if v, ok := d.GetOk("expiration_rfc3339"); ok && !d.IsNewResource() {
		now := time.Now().UTC()
		expirationTimestamp, err := time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return fmt.Errorf("error parsing expiration_rfc3339 (%s): %s", v.(string), err)
		}

		if now.After(expirationTimestamp) {
			log.Printf("[INFO] Expiration timestamp (%s) is after current timestamp (%s), removing from state", v.(string), now.Format(time.RFC3339))
			d.SetId("")
			return nil
		}
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

	if err := d.Set("rfc822", timestamp.Format(time.RFC822)); err != nil {
		return fmt.Errorf("error setting rfc822: %s", err)
	}

	if err := d.Set("rfc822z", timestamp.Format(time.RFC822Z)); err != nil {
		return fmt.Errorf("error setting rfc822z: %s", err)
	}

	if err := d.Set("rfc850", timestamp.Format(time.RFC850)); err != nil {
		return fmt.Errorf("error setting rfc850: %s", err)
	}

	if err := d.Set("rfc1123", timestamp.Format(time.RFC1123)); err != nil {
		return fmt.Errorf("error setting rfc1123: %s", err)
	}

	if err := d.Set("rfc1123z", timestamp.Format(time.RFC1123Z)); err != nil {
		return fmt.Errorf("error setting rfc1123z: %s", err)
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

	if err := d.Set("unixdate", timestamp.Format(time.UnixDate)); err != nil {
		return fmt.Errorf("error setting unixdate: %s", err)
	}

	if err := d.Set("year", timestamp.Year()); err != nil {
		return fmt.Errorf("error setting year: %s", err)
	}

	return nil
}

func resourceTimeStaticUpdate(d *schema.ResourceData, m interface{}) error {
	if d.HasChanges("expiration_days", "expiration_hours", "expiration_minutes", "expiration_months", "expiration_years") {
		timestamp, err := time.Parse(time.RFC3339, d.Id())

		if err != nil {
			return fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
		}

		var expirationTimestamp *time.Time

		if v, ok := d.GetOk("expiration_days"); ok {
			expirationTimestamp = timePtr(timestamp.AddDate(0, 0, v.(int)))
		}

		if v, ok := d.GetOk("expiration_hours"); ok {
			expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Hour))
		}

		if v, ok := d.GetOk("expiration_minutes"); ok {
			expirationTimestamp = timePtr(timestamp.Add(time.Duration(v.(int)) * time.Minute))
		}

		if v, ok := d.GetOk("expiration_months"); ok {
			expirationTimestamp = timePtr(timestamp.AddDate(0, v.(int), 0))
		}

		if v, ok := d.GetOk("expiration_years"); ok {
			expirationTimestamp = timePtr(timestamp.AddDate(v.(int), 0, 0))
		}

		if expirationTimestamp != nil {
			if err := d.Set("expiration_rfc3339", expirationTimestamp.Format(time.RFC3339)); err != nil {
				return fmt.Errorf("error setting expiration_rfc3339: %s", err)
			}
		}
	}

	return resourceTimeStaticRead(d, m)
}

func resourceTimeStaticConditionExpirationChange(diff *schema.ResourceDiff, meta interface{}) bool {
	return diff.HasChange("expiration_days") ||
		diff.HasChange("expiration_hours") ||
		diff.HasChange("expiration_minutes") ||
		diff.HasChange("expiration_months") ||
		diff.HasChange("expiration_rfc3339") ||
		diff.HasChange("expiration_years")
}

func timePtr(t time.Time) *time.Time {
	return &t
}
