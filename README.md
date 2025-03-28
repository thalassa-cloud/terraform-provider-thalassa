# Terraform Provider Thalassa Cloud

Thalassa Cloud Terraform Provider

## Documentation

- [registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs](https://registry.terraform.io/providers/thalassa-cloud/thalassa/latest/docs)

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/install)

## Examples

TODO

## License

[Apache 2.0 License 2.0](lICENSE)

## Contributing

Set up the provider locally and make sure to change the variables if needed:
```bash
make install NAMESPACE=local HOSTNAME=terraform.local OS_ARCH=darwin_arm64
```

Use the locally installed provider in your Terraform configuration:
```hcl
terraform {
  required_providers {
    thalassa = {
      source = "thalassa.cloud/thalassa/thalassa"
    }
  }
}
```
