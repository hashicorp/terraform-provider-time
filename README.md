# terraform-provider-time

Terraform Provider for time-based resources.

Please note: Issues on this repository are intended to be related to bugs or feature requests with this particular provider codebase. See [Terraform Community](https://www.terraform.io/community.html) for a list of resources to ask questions about Terraform or this provider and its usage.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 0.10+

## Using the Provider

This Terraform Provider is not available to install automatically via `terraform init` at this time. Instead, follow the instructions to [install it as a plugin](https://www.terraform.io/docs/plugins/basics.html#installing-plugins). Pre-built releases of this Terraform Provider are available for download in [GitHub Releases](https://github.com/bflad/terraform-provider-time/releases). After placing the custom provider into your plugins directory, run `terraform init` to initialize it.

### Resource Documentation

Until this Terraform Provider is brought under the Terraform Provider development program, resource documentation can be found within this repository.

- [`time_offset` Resource](./website/docs/r/offset.html.markdown)
- [`time_rotating` Resource](./website/docs/r/rotating.html.markdown)
- [`time_static` Resource](./website/docs/r/static.html.markdown)

## Developing the Provider

If you wish to work on the provider, you'll first need [Go 1.13 or later](http://www.golang.org) installed on your machine. This project uses [Go Modules](https://blog.golang.org/using-go-modules) making it safe to work with it outside of your existing [GOPATH](http://golang.org/doc/code.html#GOPATH).

### Building the Provider

From the top directory of the repository:

```console
$ go build
```

A `terraform-provider-time` binary will be left in the current directory.

### Testing the Provider

From the top directory of the repository:

```console
$ go test ./...
```
