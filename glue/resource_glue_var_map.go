package glue

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func resourceGlueVarMap() *schema.Resource {
	return &schema.Resource{
		Create: resourceGlueVarMapCreate,
		Read:   resourceGlueVarMapRead,
		Update: resourceGlueVarMapUpdate,
		Delete: resourceGlueVarMapDelete,

		Schema: map[string]*schema.Schema{
			"identifier": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"entries": &schema.Schema{
				Type:     schema.TypeMap,
				Required: true,
			},
			"keys": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"values": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"entry_count": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func resourceGlueVarMapCreate(d *schema.ResourceData, meta interface{}) error {
	entries := d.Get("entries").(map[string]interface{})
	identifier := d.Get("identifier").(string)

	entryCount := len(entries)
	keys := make([]string, entryCount)
	values := make([]string, entryCount)

	i := 0
	for k, v := range entries {
		keys[i] = k
		values[i] = v.(string)
		i++
	}

	if err := d.Set("keys", keys); err != nil {
		return err
	}

	if err := d.Set("values", values); err != nil {
		return err
	}

	if err := d.Set("entry_count", entryCount); err != nil {
		return err
	}

	d.SetId(identifier)

	return nil
}

func resourceGlueVarMapRead(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGlueVarMapUpdate(d *schema.ResourceData, meta interface{}) error {
	return nil
}

func resourceGlueVarMapDelete(d *schema.ResourceData, meta interface{}) error {
	d.SetId("")
	return nil
}
