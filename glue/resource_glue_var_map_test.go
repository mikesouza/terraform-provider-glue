package glue

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"sort"
	"strconv"
)

var testAccGlueVarMapResourceConfig = `
provider "glue" {
	# TODO
}

resource "glue_var_map" "test" {
	identifier = "someUniqueId"
	entries = {
		key1 = "value1",
		key2 = "value2"
	}
}

output "keys" {
	value = "${glue_var_map.test.keys}"
}

output "values" {
	value = "${glue_var_map.test.values}"
}

output "entry_count" {
	value = "${glue_var_map.test.entry_count}"
}
`

func TestAccGlueVarMapResource(t *testing.T) {
	n := "glue_var_map.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccGlueVarMapResourceConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckResourceID(n),
					testAccCheckGlueVarMapResourceAttributes(n),
				),
			},
		},
		CheckDestroy: testAccCheckGlueVarMapResourceDestroy,
	})
}

func testAccCheckGlueVarMapResourceAttributes(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		expected := map[string]interface{}{
			"entry_count": 2,
			"keys":        []string{"key1", "key2"},
			"values":      []string{"value1", "value2"},
		}

		errorMsgFormat := "output %q is %q; expect %q of type %q"

		for k, v := range expected {
			if outputs[k] == nil {
				return fmt.Errorf("missing output %q", k)
			}

			switch t := v.(type) {
			case int:
				expectedV := v.(int)
				outputV, _ := strconv.Atoi(outputs[k].Value.(string))
				if outputV != expectedV {
					return fmt.Errorf(errorMsgFormat, k, outputV, expectedV, t)
				}
			case []string:
				expectedV := v.([]string)
				outputV := ArrayToStrings(outputs[k].Value.([]interface{}))
				if len(outputV) != len(expectedV) {
					return fmt.Errorf(errorMsgFormat, k, outputV, expectedV, t)
				}

				sort.Strings(outputV)
				for i, v := range expectedV {
					if outputV[i] != v {
						return fmt.Errorf(errorMsgFormat, k, outputV, expectedV, t)
					}
				}
			default:
				return fmt.Errorf(errorMsgFormat, k, outputs[k].Value, v, t)
			}
		}

		return nil
	}
}

func testAccCheckGlueVarMapResourceDestroy(s *terraform.State) error {
	return nil
}
