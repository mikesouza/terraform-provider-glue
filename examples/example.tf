provider "glue" {
  env_path = "/optionally/add/to;/env/path"
}

data "glue_command" "example" {
  command    = "aws"
  parameters = ["--version"]
}

data "glue_filter_regexp" "example" {
  input      = "${data.glue_command.example.output}"
  expression = "aws-cli/([0-9.]+)"
}

output "awscli_version" {
  value = "${data.glue_filter_regexp.example.output[1]}"
}

resource "glue_var_map" "example" {
  identifier = "my-map-id"

  entries = {
    foo   = "bar"
    hello = "world"
    key1  = "value1"
    key2  = "value2"
  }
}

data "glue_filter_map" "example" {
  depends_on = ["glue_var_map.example"]

  input = "${glue_var_map.example.entries}"

  key = {
    prefix = ["key"]
  }
}

output "filtered_keys" {
  value = "${data.glue_filter_map.example.output}"
}

data "glue_filter_jmespath" "example" {
  input      = "${jsonencode(glue_var_map.example.entries)}"
  expression = "hello"
}

output "filtered_json" {
  value = "${data.glue_filter_jmespath.example.output}"
}
