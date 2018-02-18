package glue

import (
	"fmt"
	"os"
	"testing"
)

func TestAccGlueConfigEnvPath(t *testing.T) {
	oldPathVar := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPathVar)

	mockExistingPathVar := "/mock;existing/path"
	err := os.Setenv("PATH", mockExistingPathVar)
	if err != nil {
		t.Fatalf("Expected no error, received: %s", err)
	}

	addedPathVar := "added/path"
	config := Config{
		EnvPath: addedPathVar,
	}

	_, err = config.Client()
	if err != nil {
		t.Fatalf("Expected no error, received: %s", err)
	}

	expectedPath := fmt.Sprintf("%s;%s", addedPathVar, mockExistingPathVar)
	if newPathVar := os.Getenv("PATH"); newPathVar != expectedPath {
		t.Fatalf("Expected PATH = %s; actual PATH = %s", expectedPath, newPathVar)
	}
}
