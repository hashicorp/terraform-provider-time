package tftime

import (
	"fmt"
	"regexp"
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
					return nil, fmt.Errorf("Unexpected format of ID (%q), expected CREATEDURATION,DESTROYDURATION where at least one value is non-empty", d.Id())
				}

				if idParts[0] != "" {
					if _, err := time.ParseDuration(idParts[0]); err != nil {
						return nil, fmt.Errorf("error parsing create_duration (%s): %w", idParts[0], err)
					}

					if err := d.Set("create_duration", idParts[0]); err != nil {
						return nil, fmt.Errorf("error setting create_duration: %s", err)
					}
				}

				if idParts[1] != "" {
					if _, err := time.ParseDuration(idParts[1]); err != nil {
						return nil, fmt.Errorf("error parsing destroy_duration (%s): %w", idParts[1], err)
					}

					if err := d.Set("destroy_duration", idParts[1]); err != nil {
						return nil, fmt.Errorf("error setting destroy_duration: %s", err)
					}
				}

				d.SetId(time.Now().UTC().Format(time.RFC3339))

				return []*schema.ResourceData{d}, nil
			},
		},

		Schema: map[string]*schema.Schema{
			"create_duration": {
				Type:     schema.TypeString,
				Optional: true,
				AtLeastOneOf: []string{
					"create_duration",
					"destroy_duration",
				},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`), "must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
			},
			"destroy_duration": {
				Type:     schema.TypeString,
				Optional: true,
				AtLeastOneOf: []string{
					"create_duration",
					"destroy_duration",
				},
				ValidateFunc: validation.StringMatch(regexp.MustCompile(`^[0-9]+(\.[0-9]+)?(ms|s|m|h)$`), "must be a number immediately followed by ms (milliseconds), s (seconds), m (minutes), or h (hours). For example, \"30s\" for 30 seconds."),
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
	if v, ok := d.GetOk("create_duration"); ok {
		duration, err := time.ParseDuration(v.(string))

		if err != nil {
			return fmt.Errorf("error parsing create_duration (%s): %w", v.(string), err)
		}

		time.Sleep(duration)
	}

	d.SetId(time.Now().UTC().Format(time.RFC3339))

	return nil
}

func resourceTimeSleepDelete(d *schema.ResourceData, m interface{}) error {
	if v, ok := d.GetOk("destroy_duration"); ok {
		duration, err := time.ParseDuration(v.(string))

		if err != nil {
			return fmt.Errorf("error parsing destroy_duration (%s): %w", v.(string), err)
		}

		time.Sleep(duration)
	}

	return nil
}
