package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"os/exec"
	"path"
)

func dataSourceGlueCommand() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceGlueCommandRead,

		Schema: map[string]*schema.Schema{
			"command": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
			},
			"parameters": &schema.Schema{
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"output": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func validateCommand(command string) error {
	if len(command) < 1 {
		return fmt.Errorf("command must not be empty")
	}

	// first element is assumed to be an executable command, possibly found
	// using the PATH environment variable.
	_, err := exec.LookPath(command)
	if err != nil {
		return fmt.Errorf("%s: command not found", command)
	}

	return nil
}

func dataSourceGlueCommandRead(d *schema.ResourceData, meta interface{}) error {
	command := d.Get("command").(string)
	parameters := ArrayToStrings(d.Get("parameters").([]interface{}))

	if err := validateCommand(path.Join(command)); err != nil {
		return err
	}

	cmd := exec.Command(command, parameters...)

	result, err := cmd.Output()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if exitErr.Stderr != nil && len(exitErr.Stderr) > 0 {
				return fmt.Errorf("failed to execute %q: %s", command, string(exitErr.Stderr))
			}
			return fmt.Errorf("command %s failed with no error message", command)
		}

		return fmt.Errorf("failed to execute %s: %s", command, err)
	}

	if err := d.Set("output", fmt.Sprintf("%s", result)); err != nil {
		return err
	}

	d.SetId("glue_command")

	return nil
}
