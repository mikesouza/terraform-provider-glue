package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceGlueFilterMap() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlueFilterMapRead,

		Schema: map[string]*schema.Schema{
			"input": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
			"key": &schema.Schema{
				Type:     schema.TypeSet,
				Optional: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"equals": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"contains": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"prefix": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
						"suffix": {
							Type:     schema.TypeList,
							Optional: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"output": &schema.Schema{
				Type:     schema.TypeMap,
				Computed: true,
			},
		},
	}
}

func dataSourceGlueFilterMapRead(d *schema.ResourceData, meta interface{}) error {
	input := d.Get("input").(map[string]interface{})
	if input == nil {
		return fmt.Errorf("missing input")
	}

	m := MapFilter{
		Input: input,
	}

	if keyFilterConfig, ok := d.GetOk("key"); ok {
		if keyFilterSet, ok := keyFilterConfig.(*schema.Set); ok && keyFilterSet.Len() > 0 {
			for _, keyFilters := range keyFilterSet.List() {
				for filterName, filterInput := range keyFilters.(map[string]interface{}) {
					m.SetKeyFilter(filterName, filterInput.([]interface{}))
				}
			}
		}
	}

	m.Apply()

	if err := d.Set("output", m.Output); err != nil {
		return err
	}

	d.SetId("glue_filter_map")

	return nil
}
