package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"testing"
)

const testAccGlueFilterMapDataSourceName = "data.glue_filter_map.test"

const testAccGlueFilterMapDataSourceConfig = `
provider "glue" {
	# TODO
}

resource "glue_var_map" "test" {
	identifier = "someUniqueId"
	entries = {
%s
	}
}

data "glue_filter_map" "test" {
	input = "${glue_var_map.test.entries}"
	key {
%s
	}
}

output "output" {
	value = "${data.glue_filter_map.test.output}"
}
`

const testAccGlueFilterMapDataSourceConfigEntries = `
		aKey1 = "value1"
		bKey2 = "value2"
		cK3y3 = "value3"
		dKey4 = "value4"
`

const testAccGlueFilterMapDataSourceConfigKeyEquals = `
		equals = ["%s"]
`

const testAccGlueFilterMapDataSourceConfigKeyContains = `
		contains = ["K3y"]
`

const testAccGlueFilterMapDataSourceConfigKeyPrefix = `
		prefix = ["b", "c"]
`

const testAccGlueFilterMapDataSourceConfigKeySuffix = `
		suffix = ["3"]
`

const testAccGlueFilterMapDataSourceConfigMultipleKeyFilters = `
		equals   = ["dKey4"]
		contains = ["Key"]
		prefix   = ["a"]
`

func TestAccGlueFilterMapDataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					""),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceCountOutput(4),
				),
			},
		},
	})
}

func TestAccGlueFilterMapDataSourceKeyEquals(t *testing.T) {
	key := "bKey2"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					fmt.Sprintf(testAccGlueFilterMapDataSourceConfigKeyEquals, key)),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceKeysOutput(key),
				),
			},
		},
	})
}

func TestAccGlueFilterMapDataSourceKeyContains(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					testAccGlueFilterMapDataSourceConfigKeyContains),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceCountOutput(1),
				),
			},
		},
	})
}

func TestAccGlueFilterMapDataSourceKeyPrefix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					testAccGlueFilterMapDataSourceConfigKeyPrefix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceCountOutput(2),
				),
			},
		},
	})
}

func TestAccGlueFilterMapDataSourceKeySuffix(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					testAccGlueFilterMapDataSourceConfigKeySuffix),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceCountOutput(1),
				),
			},
		},
	})
}

func TestAccGlueFilterMapDataSourceMultipleKeyFilters(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: fmt.Sprintf(
					testAccGlueFilterMapDataSourceConfig,
					testAccGlueFilterMapDataSourceConfigEntries,
					testAccGlueFilterMapDataSourceConfigMultipleKeyFilters),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDataSourceID(testAccGlueFilterMapDataSourceName),
					testAccCheckGlueFilterMapDataSourceCountOutput(3),
				),
			},
		},
	})
}

func testAccCheckGlueFilterMapDataSourceKeysOutput(keys ...string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		output, err := validateKeyFilterOutputCount(outputs, len(keys))
		if err != nil {
			return err
		}

		for _, key := range keys {
			_, err = validateKeyFilterOutputValue(output, key, nil)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func testAccCheckGlueFilterMapDataSourceCountOutput(expectedCount int) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		outputs := s.RootModule().Outputs

		_, err := validateKeyFilterOutputCount(outputs, expectedCount)
		return err
	}
}

func validateKeyFilterOutputCount(outputs map[string]*terraform.OutputState, expectedCount int) (map[string]interface{}, error) {
	if outputs["output"] == nil {
		return nil, fmt.Errorf("missing output")
	}

	output := outputs["output"].Value.(map[string]interface{})
	if len(output) != expectedCount {
		return nil, fmt.Errorf("output is length %d; expect length %d", len(output), expectedCount)
	}

	return output, nil
}

func validateKeyFilterOutputValue(output map[string]interface{}, key string, expectedValue interface{}) (interface{}, error) {
	value, ok := output[key]
	if !ok {
		return nil, fmt.Errorf("output is missing key '%q'", key)
	}

	if expectedValue != nil {
		if value != expectedValue {
			return value, fmt.Errorf("output key '%q' value is '%q'; expect '%q'", key, value, expectedValue)
		}
	}

	return value, nil
}
