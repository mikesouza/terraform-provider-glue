package glue

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

// Provider returns a terraform.ResourceProvider.
func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"env_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("GLUE_ENV_PATH", nil),
				Description: "The environment PATH variable to use for Glue commands.",
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"glue_var_map": resourceGlueVarMap(),
		},

		DataSourcesMap: map[string]*schema.Resource{
			"glue_command":         dataSourceGlueCommand(),
			"glue_filter_jmespath": dataSourceGlueFilterJMESPath(),
			"glue_filter_regexp":   dataSourceGlueFilterRegExp(),
			"glue_filter_map":      dataSourceGlueFilterMap(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		EnvPath: d.Get("env_path").(string),
	}

	return config.Client()
}
