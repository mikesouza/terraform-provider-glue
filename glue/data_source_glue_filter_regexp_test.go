package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"strings"
	"testing"
)

const testAccGlueFilterRegExpDataSourceConfig = `
provider "glue" {
	# TODO
}

data "glue_command" "test" {
	command = "%s"
	parameters = ["text", "\"Hello World\""]
}

data "glue_filter_regexp" "test" {
	input = "${data.glue_command.test.output}"
	expression = "Hello ([A-z]+)"
}

output "output" {
	value = "${data.glue_filter_regexp.test.output}"
}
`

func TestAccGlueFilterRegExpDataSource(t *testing.T) {
	commandPath, err := buildGlueCommandTestProgram()
	if err != nil {
		t.Fatal(err)
		return
	}

	n := "data.glue_filter_regexp.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccGlueFilterRegExpDataSourceConfig, commandPath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(n),
					testAccCheckGlueFilterRegExpDataSourceOutput(n),
				),
			},
		},
	})
}

func testAccCheckGlueFilterRegExpDataSourceOutput(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		if outputs["output"] == nil {
			return fmt.Errorf("missing output")
		}

		output := ArrayToStrings(outputs["output"].Value.([]interface{}))

		if len(output) < 2 {
			return fmt.Errorf("output is array [%v]; expect [Hello World, World]", strings.Join(output, ", "))
		}

		return nil
	}
}
