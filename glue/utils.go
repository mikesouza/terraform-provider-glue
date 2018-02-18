package glue

import (
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"math/rand"
	"time"
)

func SetSeed() {
	rand.Seed(time.Now().Unix())
}

func StringsToArray(s []string) []interface{} {
	o := make([]interface{}, len(s))
	for i, v := range s {
		o[i] = v
	}

	return o
}

func ArrayToStrings(o []interface{}) []string {
	s := make([]string, len(o))
	for i, v := range o {
		s[i] = v.(string)
	}

	return s
}

func StringsToSet(slice []string) map[string]bool {
	set := make(map[string]bool, len(slice))
	for _, s := range slice {
		set[s] = true
	}

	return set
}

func ArrayToSet(slice []interface{}) map[interface{}]bool {
	set := make(map[interface{}]bool, len(slice))
	for _, s := range slice {
		set[s] = true
	}

	return set
}

func PrintResourceAttributes(s *terraform.State, n string) error {
	rs, ok := s.RootModule().Resources[n]
	if !ok {
		return fmt.Errorf("Can't find resource: %s", n)
	}

	if rs.Primary.ID == "" {
		return fmt.Errorf("Resource ID not set")
	}

	for k, v := range rs.Primary.Attributes {
		fmt.Printf("%s = %s\n", k, v)
	}

	return nil
}
