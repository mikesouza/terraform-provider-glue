package glue

import (
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/config"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProvider *schema.Provider
var testAccProviders map[string]terraform.ResourceProvider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"glue": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func TestProvider_configure(t *testing.T) {
	c, err := config.NewRawConfig(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
	rc := terraform.NewResourceConfig(c)
	p := Provider()
	err = p.Configure(rc)
	if err != nil {
		t.Fatal(err)
	}
}

func testAccPreCheck(t *testing.T) {
	err := testAccProvider.Configure(terraform.NewResourceConfig(nil))
	if err != nil {
		t.Fatal(err)
	}

	if v := os.Getenv("GLUE_ENV_PATH"); v == "" {
		t.Fatal("GLUE_ENV_PATH must be set for acceptance tests")
	}
}

func testAccCheckDataSourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find data source: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Data source ID not set")
		}
		return nil
	}
}

func testAccCheckResourceID(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Can't find resource: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("Resource ID not set")
		}
		return nil
	}
}
