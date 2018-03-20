# Terraform Glue Provider

Terraform Provider for "gluing" together the inputs and outputs from external sources and other providers.

It provides a flexible set of data sources for executing shell commands and filtering output by [JMESPath][4] or [RE2][5] regular expressions, as well as resources for persisting variables in the [State][3].

For general information about Terraform, visit the [official website][1] and [GitHub project page][2].

[1]: https://terraform.io/
[2]: https://github.com/hashicorp/terraform
[3]: https://www.terraform.io/docs/state/
[4]: http://jmespath.org/
[5]: https://github.com/google/re2/wiki/Syntax

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10.x (or higher)
- [Go](https://golang.org/doc/install) 1.9+ (to build the provider plugin)

## Using the Provider

If you're building the provider, follow the instructions to [install it as a plugin.](https://www.terraform.io/docs/plugins/basics.html#installing-a-plugin) After placing it into your plugins directory, run `terraform init` to initialize it.

### Example configuration

```hcl
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
```

See also `examples` folders for more details.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.9+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

Clone repository to: `$GOPATH/src/github.com/MikeSouza/terraform-provider-glue`

```sh
$ mkdir -p $GOPATH/src/github.com/MikeSouza; cd $GOPATH/src/github.com/MikeSouza
$ git clone git@github.com:MikeSouza/terraform-provider-glue
...
```

Enter the provider directory and build the provider:

```sh
$ cd $GOPATH/src/github.com/MikeSouza/terraform-provider-glue
$ make build
...
```

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-glue
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
...
```

In order to run the full suite of Acceptance tests, run `make testacc`.

```sh
$ make testacc
...
```
