package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccGlueFilterJMESPathDataSourceConfig = `
provider "glue" {
	# TODO
}

data "glue_command" "test" {
	command = "%s"
	parameters = ["json", "foo", "bar"]
}

data "glue_filter_jmespath" "test" {
	input = "${data.glue_command.test.output}"
	expression = "parameter1"
}

output "output" {
	value = "${data.glue_filter_jmespath.test.output}"
}
`

func TestAccGlueFilterJMESPathDataSource(t *testing.T) {
	commandPath, err := buildGlueCommandTestProgram()
	if err != nil {
		t.Fatal(err)
		return
	}

	n := "data.glue_filter_jmespath.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccGlueFilterJMESPathDataSourceConfig, commandPath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(n),
					testAccCheckGlueFilterJMESPathDataSourceOutput(n),
				),
			},
		},
	})
}

func testAccCheckGlueFilterJMESPathDataSourceOutput(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		if outputs["output"] == nil {
			return fmt.Errorf("missing output")
		}

		output := outputs["output"].Value
		if output != "bar" {
			return fmt.Errorf("output is %q; expect 'bar'", output)
		}

		return nil
	}
}
