package glue

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/jmespath/go-jmespath"
)

func dataSourceGlueFilterJMESPath() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlueFilterJMESPathRead,

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
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceGlueFilterJMESPathRead(d *schema.ResourceData, meta interface{}) error {
	input := d.Get("input").(string)
	inputBytes := []byte(input)
	var data map[string]interface{}
	err := json.Unmarshal(inputBytes, &data)
	if err != nil {
		return err
	}

	expression := d.Get("expression").(string)
	path, err := jmespath.Compile(expression)
	if err != nil {
		return err
	}

	output, err := path.Search(data)
	if err != nil {
		return err
	}

	if err := d.Set("output", fmt.Sprintf("%s", output)); err != nil {
		return err
	}

	d.SetId("glue_filter_jmespath")

	return nil
}
