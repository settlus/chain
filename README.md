# Settlus

**Settlus** is a purpose-built blockchain designed to provide a transparent settlement system for the creator economy.

## Get started
### Install necessary dependencies
```shell
curl https://get.ignite.com/cli! | bash
brew install golangci-lint
```

### Run the chain
```shell
ignite chain serve --skip-proto
```

`serve` command installs dependencies, builds, initializes, and starts your blockchain in development.

## Development
### Build
```shell
make
```

### Lint
```shell
make lint
```

### Test
```shell
make test
```

### Generate protobuf definition
```shell
make proto-gen
```

### Generate Swagger
```shell
make proto-swagger-gen
```

### Local network test
```shell
make localnet-build
make localnet-start

# stop local network test
make localnet-stop
```

## License
This project is licensed under the [LGPL-3.0 license](LICENSE). Specifically, the contents within the **evmos** folder are dervied from remarkable work originally pioneered by the [Evmos Foundation](https://evmos.org/), also under the LGPL-3.0 license. We have utilized code from [Evmos v12](https://github.com/evmos/evmos/commits/release/v12.x.x/) (**`b43ee16`**) and have made several modifications. All changes are documented [here](/evmos/CHANGES.diff), as well as in the commit logs of this repository.