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
make chaind
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
Export personal github token to run dockerfile in local
```shell
GITHUB_TOKEN=`git config --global --list | grep 'url' | sed 's/.*git://' | awk -F '[@]' '{print $1}'` && export GITHUB_TOKEN
```

```shell
make localnet-build # first time only
make localnet-start

unset GITHUB_TOKEN # for safety
```
