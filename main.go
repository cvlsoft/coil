package main

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configer provides an identifier interface for all configuration types
type Configer interface {
	generate()
	getParser() *viper.Viper
}

// AuthConfig represents a composable struct for db connections
type AuthConfig struct {
	JWTPhrase string `type:"string" name:"jwt_phrase" default:"82uushdf8h2398ru09sduf" desc:"Phrase for generating tokens"`
}

// Config is a standard definition for config interfaces
type Config struct {
	viper *viper.Viper
}

// getParser returns the current parser instance
func (c *Config) getParser() *viper.Viper {
	return c.viper
}

// generate adds generators to the register
func (c *Config) generate() {
	pflag.String("config", "", "Path for a configuration file to load")
	c.viper = CreateViper()
}

// defineFlagsFromStruct performs a deep recurse into the specified object
// to find tags and declare them against viper
func defineFlagsFromStruct(t reflect.Type) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			defineFlagsFromStruct(field.Type)
			continue
		}
		flagType := field.Tag.Get("type")
		// Define flags based on their types
		switch flagType {
		case "string":
			pflag.String(field.Tag.Get("name"), field.Tag.Get("default"), field.Tag.Get("desc"))
		case "int":
			i, err := strconv.Atoi(field.Tag.Get("default"))
			if err == nil {
				pflag.Int64(field.Tag.Get("name"), int64(i), field.Tag.Get("desc"))
			}
		case "bool":
			var val bool = false
			if field.Tag.Get("default") == "true" {
				val = true
			}
			pflag.Bool(field.Tag.Get("name"), val, field.Tag.Get("desc"))
		case "float32":
			i, err := strconv.Atoi(field.Tag.Get("default"))
			if err == nil {
				pflag.Float32(field.Tag.Get("name"), float32(i), field.Tag.Get("desc"))
			}
		case "float64":
			i, err := strconv.Atoi(field.Tag.Get("default"))
			if err == nil {
				pflag.Float64(field.Tag.Get("name"), float64(i), field.Tag.Get("desc"))
			}
		}
	}
}

// setPropertiesFromFlags performs a deep recurse into the specified object
// to retrieve and bind them to the struct
func setPropertiesFromFlags(vp reflect.Value, viper *viper.Viper) {
	v := vp.Elem()
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		switch field.Type.Kind() {
		case reflect.Struct:
			setPropertiesFromFlags(v.Field(i).Addr(), viper)
		case reflect.String:
			v.Field(i).SetString(viper.GetString(field.Tag.Get("name")))
		case reflect.Bool:
			v.Field(i).SetBool(viper.GetBool(field.Tag.Get("name")))
		case reflect.Int:
			v.Field(i).SetInt(viper.GetInt64(field.Tag.Get("name")))
		case reflect.Float32:
		case reflect.Float64:
			v.Field(i).SetFloat(viper.GetFloat64(field.Tag.Get("name")))
		}
	}
	// Finally detect if a parse method exists and trigger it
	method := vp.MethodByName("Parse")
	if method.IsValid() {
		method.Call([]reflect.Value{reflect.ValueOf(viper)})
	}
}

// NewConfig generates a new configuration setup
func NewConfig(c Configer) Configer {
	defineFlagsFromStruct(reflect.TypeOf(c).Elem())
	c.generate()
	setPropertiesFromFlags(reflect.ValueOf(c), c.getParser())
	return c
}

// CreateViper creates a parser instance to configure CLI.
// It can be used for packages that re-implement the command line flags
func CreateViper() (v *viper.Viper) {
	// Read configurations and assign them
	v = viper.New()
	v.AutomaticEnv()
	pflag.Parse()
	v.BindPFlags(pflag.CommandLine)
	// Override values if they exist already
	if v.GetString("config") != "" {
		v.SetConfigFile(v.GetString("config"))
		if err := v.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				panic("Could not find configuration file")
			} else {
				fmt.Println(err)
				panic("Could not parse configuration file")
			}
		}
	}
	return
}
