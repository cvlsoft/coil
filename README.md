# Coil

Coil is a lightweight Go package for configuration management, built on top of Viper and Cobra. It simplifies configuration handling by allowing you to compose and layer settings using Go structs. Configuration values are automatically exposed as both CLI flags and config file options, providing flexible deployment options for your applications.

## ➕ Install

```bash
go get github.com/cvlstack/coil
```

## 📦 Prebuilt Configurations

Coil ships with a number of very basic and useful configurations. These out-of-the-box options cover basic functions like API keys, database connections, and authentication etc:

- `coil.Config`: Base Coil configuration used on all struct definitions.
- `coil.APIServiceConfig`: Defines fundamental configurations for an API service
- `coil.DatabaseConfig`: Helps define standard database connection details.

We hope to expand this list of predefined types with community contributions.

## ⚡️ Quickstart

```go
package config

import (
	configs "github.com/cvlstack/coil"
)

// Config represents your app's local config
type Config struct {
	coil.Config
	coil.APIServiceConfig
	coil.DatabaseConfig
	MyCustomConfig
}

// MyCustomConfig representsa custom configuration
type MyCustomConfig struct {
	FooBar string `type:"string" name:"foo_bar" default:"static" desc:"Foo bar value"`
}

// NewConfig is a factory generator for your configuration
func NewConfig() *Config {
	c := coil.NewConfig(&Config{})
	return c.(*Config)
}
```
This simple declaration will allow you to define your YAML config file like so:
```yaml
name: "prod"
dbhost: 192.168.1.1
dbname: my_name
dbuser: my_user
dbpass: my_pass
dbport: 25061
port: 8443
host: 0.0.0.0
api_url: https://api.myservice.com
foo_bar: override
```
You can also use any of the configuration format supported by Viper (JSON, TOML, YAML, ENV, etc.). You can also automatically use your CLI for defining config overrides:
```bash
go run main.go --foo_bar=dynamic
```

## 🔀 Using Prefixes for Multiple Instances

When you need to use the same configuration type multiple times (e.g., multiple database connections), use the `prefix` tag to avoid naming collisions:

```go
type Config struct {
    coil.Config
    PrimaryDB coil.DatabaseConfig `prefix:"primary"`
    ReplicaDB coil.DatabaseConfig `prefix:"replica"`
}
```

This allows you to configure each instance independently via environment variables:

```bash
# Primary database
export PRIMARY_DBHOST=primary.example.com
export PRIMARY_DBPORT=5432
export PRIMARY_DBUSER=primary_user

# Replica database
export REPLICA_DBHOST=replica.example.com
export REPLICA_DBPORT=5433
export REPLICA_DBUSER=replica_user
```

Or via CLI flags:

```bash
go run main.go --primary_dbhost=primary.example.com --replica_dbhost=replica.example.com
```

Or in your config file:

```yaml
primary_dbhost: primary.example.com
primary_dbport: 5432
replica_dbhost: replica.example.com
replica_dbport: 5433
```

## 🌐 Community Contributions

We welcome contributions from the community to expand the list of predefined types. If you have a configuration type that you think would be useful for others, please submit a pull request with your contribution.

## 📃 License

Coil is released under the Apache License 2.0. See [LICENSE](LICENSE) for details.
