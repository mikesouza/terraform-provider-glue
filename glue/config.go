package glue

import (
	"fmt"
	"log"
	"os"
)

// Config represents the provider configuration
type Config struct {
	EnvPath string
}

// Client returns the configured provider client.
func (c *Config) Client() (interface{}, error) {
	var err error

	if c.EnvPath != "" {
		var pathVar string
		if pathVar = os.Getenv("PATH"); pathVar != "" {
			pathVar = fmt.Sprintf("%s;%s", c.EnvPath, pathVar)
		} else {
			pathVar = c.EnvPath
		}

		err = os.Setenv("PATH", pathVar)
		log.Printf("[INFO] Glue environment PATH = %s", pathVar)
	}

	return nil, err
}
