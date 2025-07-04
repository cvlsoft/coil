# Coil Config

A small and handy Go configuration composition package built on Viper and Cobra. Coil makes it easy to define and stack your configurations with composed structs. Your config settings are instantly available as CLI flags or config files.

## ➕ Install

```bash
go get github.com/cvlsoft/coil
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
	configs "github.com/cvlsoft/coil"
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

## 🌐 Community Contributions

We welcome contributions from the community to expand the list of predefined types. If you have a configuration type that you think would be useful for others, please submit a pull request with your contribution.

## 📃 License

Coil is released under the MIT. See [LICENSE.txt](https://github.com/cvlsoft/coil/blob/master/LICENSE.txt)
