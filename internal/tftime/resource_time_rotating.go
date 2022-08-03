package tftime

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/customdiff"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

func resourceTimeRotating() *schema.Resource {
	return &schema.Resource{
		Description: "Manages a rotating time resource, which keeps a rotating UTC timestamp stored in the Terraform " +
			"state and proposes resource recreation when the locally sourced current time is beyond the rotation time. " +
			"This rotation only occurs when Terraform is executed, meaning there will be drift between the rotation " +
			"timestamp and actual rotation. The new rotation timestamp offset includes this drift. " +
			"This prevents perpetual differences caused by using the [`timestamp()` function](https://www.terraform.io/docs/configuration/functions/timestamp.html) " +
			"by only forcing a new value on the set cadence.",
		CreateContext: resourceTimeRotatingCreate,
		ReadContext:   resourceTimeRotatingRead,
		UpdateContext: resourceTimeRotatingUpdate,
		DeleteContext: schema.NoopContext,

		CustomizeDiff: customdiff.Sequence(
			customdiff.If(resourceTimeRotatingConditionExpirationChange,
				func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) error {
					if diff.Id() == "" {
						return nil
					}

					timestamp, err := time.Parse(time.RFC3339, diff.Id())

					if err != nil {
						return fmt.Errorf("error parsing timestamp (%s): %s", diff.Id(), err)
					}

					var rotationTimestamp time.Time

					if v, ok := diff.GetOk("rotation_days"); ok {
						rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
					}

					if v, ok := diff.GetOk("rotation_hours"); ok {
						rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
					}

					if v, ok := diff.GetOk("rotation_minutes"); ok {
						rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
					}

					if v, ok := diff.GetOk("rotation_months"); ok {
						rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
					}

					if v, ok := diff.GetOk("rotation_rfc3339"); ok {
						var err error
						rotationTimestamp, err = time.Parse(time.RFC3339, v.(string))

						if err != nil {
							return fmt.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
						}
					}

					if v, ok := diff.GetOk("rotation_years"); ok {
						rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
					}

					if err := diff.SetNew("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
						return fmt.Errorf("error setting new rotation_rfc3339: %s", err)
					}

					return nil
				},
			),
			customdiff.ForceNewIf("rotation_rfc3339", func(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
				now := time.Now().UTC()
				rotationTimestamp, err := time.Parse(time.RFC3339, diff.Get("rotation_rfc3339").(string))

				if err != nil {
					return false
				}

				return now.After(rotationTimestamp)
			}),
		),

		Importer: &schema.ResourceImporter{
			StateContext: func(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
				idParts := strings.Split(d.Id(), ",")

				if len(idParts) != 2 && len(idParts) != 6 {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES or BASETIMESTAMP,ROTATIONTIMESTAMP", d.Id())
				}

				if len(idParts) == 2 {
					if idParts[0] == "" || idParts[1] == "" {
						return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,ROTATIONTIMESTAMP", d.Id())
					}

					baseRfc3339 := idParts[0]
					rotationRfc3339 := idParts[1]

					d.SetId(baseRfc3339)

					if err := d.Set("rotation_rfc3339", rotationRfc3339); err != nil {
						return nil, fmt.Errorf("error setting rotation_rfc3339: %s", err)
					}

					return []*schema.ResourceData{d}, nil
				}

				if idParts[0] == "" || (idParts[1] == "" && idParts[2] == "" && idParts[3] == "" && idParts[4] == "" && idParts[5] == "") {
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected BASETIMESTAMP,YEARS,MONTHS,DAYS,HOURS,MINUTES, where at least one rotation value is non-empty", d.Id())
				}

				baseRfc3339 := idParts[0]
				rotationYears, _ := strconv.Atoi(idParts[1])
				rotationMonths, _ := strconv.Atoi(idParts[2])
				rotationDays, _ := strconv.Atoi(idParts[3])
				rotationHours, _ := strconv.Atoi(idParts[4])
				rotationMinutes, _ := strconv.Atoi(idParts[5])

				d.SetId(baseRfc3339)

				if rotationYears > 0 {
					if err := d.Set("rotation_years", rotationYears); err != nil {
						return nil, fmt.Errorf("error setting rotation_years: %s", err)
					}
				}

				if rotationMonths > 0 {
					if err := d.Set("rotation_months", rotationMonths); err != nil {
						return nil, fmt.Errorf("error setting rotation_months: %s", err)
					}
				}

				if rotationDays > 0 {
					if err := d.Set("rotation_days", rotationDays); err != nil {
						return nil, fmt.Errorf("error setting rotation_days: %s", err)
					}
				}

				if rotationHours > 0 {
					if err := d.Set("rotation_hours", rotationHours); err != nil {
						return nil, fmt.Errorf("error setting rotation_hours: %s", err)
					}
				}

				if rotationMinutes > 0 {
					if err := d.Set("rotation_minutes", rotationMinutes); err != nil {
						return nil, fmt.Errorf("error setting rotation_minutes: %s", err)
					}
				}

				timestamp, err := time.Parse(time.RFC3339, d.Id())

				if err != nil {
					return nil, fmt.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
				}

				var rotationTimestamp time.Time

				if v, ok := d.GetOk("rotation_days"); ok {
					rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
				}

				if v, ok := d.GetOk("rotation_hours"); ok {
					rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
				}

				if v, ok := d.GetOk("rotation_minutes"); ok {
					rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
				}

				if v, ok := d.GetOk("rotation_months"); ok {
					rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
				}

				if v, ok := d.GetOk("rotation_years"); ok {
					rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
				}

				if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
					return nil, fmt.Errorf("error setting rotation_rfc3339: %s", err)
				}

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"day": {
				Description: "Number day of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"rotation_days": {
				Description: "Number of days to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_hours": {
				Description: "Number of hours to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_minutes": {
				Description: "Number of minutes to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_months": {
				Description: "Number of months to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"rotation_rfc3339": {
				Description: "Configure the rotation timestamp with an " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format of the offset timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IsRFC3339Time,
			},
			"rotation_years": {
				Description: "Number of years to add to the base timestamp to configure the rotation timestamp. " +
					"When the current time has passed the rotation timestamp, the resource will trigger recreation. " +
					"At least one of the 'rotation_' arguments must be configured.",
				Type:     schema.TypeInt,
				Optional: true,
				AtLeastOneOf: []string{
					"rotation_days",
					"rotation_hours",
					"rotation_minutes",
					"rotation_months",
					"rotation_rfc3339",
					"rotation_years",
				},
				ValidateFunc: validation.IntAtLeast(1),
			},
			"hour": {
				Description: "Number hour of timestamp.",
				Type:        schema.TypeInt,
				Computed:    true,
			},
			"triggers": {
				Description: "Arbitrary map of values that, when changed, will trigger a new base timestamp value to be saved." +
					" These conditions recreate the resource in addition to other rotation arguments. " +
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
				Description: "Base timestamp in " +
					"[RFC3339](https://datatracker.ietf.org/doc/html/rfc3339#section-5.8) format " +
					"(see [RFC3339 time string](https://tools.ietf.org/html/rfc3339#section-5.8) e.g., " +
					"`YYYY-MM-DDTHH:MM:SSZ`). Defaults to the current time.",
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
			"id": {
				Description: "RFC3339 format of the offset timestamp, e.g. `2020-02-12T06:36:13Z`.",
				Type:        schema.TypeString,
				Computed:    true,
			},
		},
	}
}

func resourceTimeRotatingCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	timestamp := time.Now().UTC()

	if v, ok := d.GetOk("rfc3339"); ok {
		var err error
		timestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return diag.Errorf("error parsing rfc3339 (%s): %s", v.(string), err)
		}
	}

	d.SetId(timestamp.Format(time.RFC3339))

	var rotationTimestamp time.Time

	if v, ok := d.GetOk("rotation_days"); ok {
		rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
	}

	if v, ok := d.GetOk("rotation_hours"); ok {
		rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
	}

	if v, ok := d.GetOk("rotation_minutes"); ok {
		rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
	}

	if v, ok := d.GetOk("rotation_months"); ok {
		rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
	}

	if v, ok := d.GetOk("rotation_rfc3339"); ok {
		var err error
		rotationTimestamp, err = time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return diag.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
		}
	}

	if v, ok := d.GetOk("rotation_years"); ok {
		rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
	}

	if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
		return diag.Errorf("error setting rotation_rfc3339: %s", err)
	}

	return resourceTimeRotatingRead(ctx, d, m)
}

func resourceTimeRotatingRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	timestamp, err := time.Parse(time.RFC3339, d.Id())

	if err != nil {
		return diag.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
	}

	if v, ok := d.GetOk("rotation_rfc3339"); ok && !d.IsNewResource() {
		now := time.Now().UTC()
		rotationTimestamp, err := time.Parse(time.RFC3339, v.(string))

		if err != nil {
			return diag.Errorf("error parsing rotation_rfc3339 (%s): %s", v.(string), err)
		}

		if now.After(rotationTimestamp) {
			log.Printf("[INFO] Expiration timestamp (%s) is after current timestamp (%s), removing from state", v.(string), now.Format(time.RFC3339))
			d.SetId("")
			return nil
		}
	}

	if err := d.Set("day", timestamp.Day()); err != nil {
		return diag.Errorf("error setting day: %s", err)
	}

	if err := d.Set("hour", timestamp.Hour()); err != nil {
		return diag.Errorf("error setting hour: %s", err)
	}

	if err := d.Set("minute", timestamp.Minute()); err != nil {
		return diag.Errorf("error setting minute: %s", err)
	}

	if err := d.Set("month", int(timestamp.Month())); err != nil {
		return diag.Errorf("error setting month: %s", err)
	}

	if err := d.Set("rfc3339", timestamp.Format(time.RFC3339)); err != nil {
		return diag.Errorf("error setting rfc3339: %s", err)
	}

	if err := d.Set("second", timestamp.Second()); err != nil {
		return diag.Errorf("error setting second: %s", err)
	}

	if err := d.Set("unix", timestamp.Unix()); err != nil {
		return diag.Errorf("error setting unix: %s", err)
	}

	if err := d.Set("year", timestamp.Year()); err != nil {
		return diag.Errorf("error setting year: %s", err)
	}

	return nil
}

func resourceTimeRotatingUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	if d.HasChanges("rotation_days", "rotation_hours", "rotation_minutes", "rotation_months", "rotation_years") {
		timestamp, err := time.Parse(time.RFC3339, d.Id())

		if err != nil {
			return diag.Errorf("error parsing timestamp (%s): %s", d.Id(), err)
		}

		var rotationTimestamp time.Time

		if v, ok := d.GetOk("rotation_days"); ok {
			rotationTimestamp = timestamp.AddDate(0, 0, v.(int))
		}

		if v, ok := d.GetOk("rotation_hours"); ok {
			rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Hour)
		}

		if v, ok := d.GetOk("rotation_minutes"); ok {
			rotationTimestamp = timestamp.Add(time.Duration(v.(int)) * time.Minute)
		}

		if v, ok := d.GetOk("rotation_months"); ok {
			rotationTimestamp = timestamp.AddDate(0, v.(int), 0)
		}

		if v, ok := d.GetOk("rotation_years"); ok {
			rotationTimestamp = timestamp.AddDate(v.(int), 0, 0)
		}

		if err := d.Set("rotation_rfc3339", rotationTimestamp.Format(time.RFC3339)); err != nil {
			return diag.Errorf("error setting rotation_rfc3339: %s", err)
		}
	}

	return resourceTimeRotatingRead(ctx, d, m)
}

func resourceTimeRotatingConditionExpirationChange(ctx context.Context, diff *schema.ResourceDiff, meta interface{}) bool {
	return diff.HasChange("rotation_days") ||
		diff.HasChange("rotation_hours") ||
		diff.HasChange("rotation_minutes") ||
		diff.HasChange("rotation_months") ||
		diff.HasChange("rotation_rfc3339") ||
		diff.HasChange("rotation_years")
}
