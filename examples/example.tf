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
    hello = "world"
    foo   = "bar"
  }
}

data "glue_filter_jmespath" "example" {
  input      = "${jsonencode(glue_var_map.example.entries)}"
  expression = "hello"
}

output "hello" {
  value = "${data.glue_filter_jmespath.example.output}"
}
