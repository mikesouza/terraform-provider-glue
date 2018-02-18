package glue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"regexp"
)

func dataSourceGlueFilterRegExp() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlueFilterRegExpRead,

		Schema: map[string]*schema.Schema{
			"input": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"expression": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"output": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceGlueFilterRegExpRead(d *schema.ResourceData, meta interface{}) error {
	input := d.Get("input").(string)

	expression := d.Get("expression").(string)
	re, err := regexp.Compile(expression)
	if err != nil {
		return err
	}

	output := re.FindStringSubmatch(input)
	if err := d.Set("output", output); err != nil {
		return err
	}

	d.SetId("glue_filter_regexp")

	return nil
}
