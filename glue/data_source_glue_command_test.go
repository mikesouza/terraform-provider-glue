package glue

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"testing"
)

const testAccGlueCommandDataSourceConfig = `
provider "glue" {
	# TODO
}

data "glue_command" "test" {
	command = "%s"
	parameters = ["json", "foo", "bar"]
}

output "output" {
	value = "${data.glue_command.test.output}"
}
`

const testAccGlueCommandDataSourceConfigError = `
data "glue_command" "test" {
	command = "%s"
}
`

func TestAccGlueCommandDataSource(t *testing.T) {
	commandPath, err := buildGlueCommandTestProgram()
	if err != nil {
		t.Fatal(err)
		return
	}

	n := "data.glue_command.test"
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(testAccGlueCommandDataSourceConfig, commandPath),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(n),
					testAccCheckGlueCommandDataSourceOutput(n),
				),
			},
		},
	})
}

func testAccCheckGlueCommandDataSourceOutput(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		if outputs["output"] == nil {
			return fmt.Errorf("missing output")
		}

		output := []byte(outputs["output"].Value.(string))
		var result map[string]string
		err := json.Unmarshal(output, &result)
		if err != nil {
			return err
		}

		parameters := []string{"foo", "bar"}

		for i := 0; i < 2; i++ {
			k := fmt.Sprintf("parameter%d", i)
			v := result[k]

			if v != parameters[i] {
				return fmt.Errorf("'%s' output is %q; expect '%s'", k, v, parameters[i])
			}
		}

		return nil
	}
}

func TestAccGlueCommandDataSourceError(t *testing.T) {
	commandPath, err := buildGlueCommandTestProgram()
	if err != nil {
		t.Fatal(err)
		return
	}

	resource.UnitTest(t, resource.TestCase{
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config:      fmt.Sprintf(testAccGlueCommandDataSourceConfigError, commandPath),
				ExpectError: regexp.MustCompile("Invalid argument count"),
			},
		},
	})
}

func buildGlueCommandTestProgram() (string, error) {
	// We have a simple Go program that we use as a stub for testing.
	cmd := exec.Command(
		"go", "install",
		"github.com/MikeSouza/terraform-provider-glue/glue/test-commands/tf-acc-glue-command-data-source",
	)
	err := cmd.Run()

	if err != nil {
		return "", fmt.Errorf("failed to build test stub program: %s", err)
	}

	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = filepath.Join(os.Getenv("HOME") + "/go")
	}

	commandPath := path.Join(
		filepath.SplitList(gopath)[0], "bin", "tf-acc-glue-command-data-source",
	)

	return commandPath, nil
}
