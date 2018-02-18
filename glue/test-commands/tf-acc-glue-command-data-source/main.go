package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// This is a minimal implementation of the external data source protocol
// intended only for use in the provider acceptance tests.
//
// In practice it's likely not much harder to just write a real Terraform
// plugin if you're going to be writing your data source in Go anyway;
// this example is just in Go because we want to avoid introducing
// additional language runtimes into the test environment.
func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintf(os.Stderr, "Invalid argument count\n")
		os.Exit(1)
	}

	switch outputType := os.Args[1]; outputType {
	case "json":
		var jsonOutput = map[string]string{}
		if len(os.Args) >= 3 {
			for i, v := range os.Args[2:] {
				jsonOutput[fmt.Sprintf("parameter%d", i)] = v
			}
		}

		outputBytes, err := json.Marshal(jsonOutput)
		if err != nil {
			panic(err)
		}
		os.Stdout.Write(outputBytes)
	case "text":
		if len(os.Args) >= 3 {
			var buffer bytes.Buffer
			for _, v := range os.Args[2:] {
				buffer.WriteString(v)
			}
			os.Stdout.Write(buffer.Bytes())
		}
	}

	os.Exit(0)
}
