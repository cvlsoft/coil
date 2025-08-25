package coil

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Configer provides an identifier interface for all configuration types
type Configer interface {
	generate()
	getParser() *viper.Viper
}

// Config is a standard definition for config interfaces
type Config struct {
	viper *viper.Viper
}

// getParser returns the current parser instance
func (c *Config) getParser() *viper.Viper {
	return c.viper
}

// HasConfig checks if a specific config type is embedded in the Config struct
func (c *Config) HasConfig(checkType any) bool {
	// Get the type we're looking for
	targetType := reflect.TypeOf(checkType)
	if targetType.Kind() == reflect.Ptr {
		targetType = targetType.Elem()
	}
	// Check all fields in the Config struct
	configType := reflect.TypeOf(*c)
	for i := 0; i < configType.NumField(); i++ {
		field := configType.Field(i)
		if field.Type == targetType {
			return true
		}
	}
	return false
}

// generate adds generators to the register
func (c *Config) generate() {
	// Create a local flagset for the config flag
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	fs.String("config", "", "Path for a configuration file to load")
	// Add to global command line if not already defined
	if pflag.CommandLine.Lookup("config") == nil {
		pflag.CommandLine.AddFlagSet(fs)
	}
	c.viper = CreateViper()
}

// defineFlagsFromStruct performs a deep recurse into the specified object
// to find tags and declare them against a flagset
func defineFlagsFromStruct(t reflect.Type, fs *pflag.FlagSet) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if field.Type.Kind() == reflect.Struct {
			defineFlagsFromStruct(field.Type, fs)
			continue
		}
		flagName := field.Tag.Get("name")
		if flagName == "" {
			continue
		}
		flagType := field.Tag.Get("type")
		// Define flags based on their types
		switch flagType {
		case "string":
			fs.String(flagName, field.Tag.Get("default"), field.Tag.Get("desc"))
		case "[]string":
			fs.StringSlice(flagName, strings.Split(field.Tag.Get("default"), ","), field.Tag.Get("desc"))
		case "int":
			i, err := strconv.Atoi(field.Tag.Get("default"))
			if err == nil {
				fs.Int64(flagName, int64(i), field.Tag.Get("desc"))
			}
		case "bool":
			var val bool = false
			if field.Tag.Get("default") == "true" {
				val = true
			}
			fs.Bool(flagName, val, field.Tag.Get("desc"))
		case "float32":
			i, err := strconv.ParseFloat(field.Tag.Get("default"), 32)
			if err == nil {
				fs.Float32(flagName, float32(i), field.Tag.Get("desc"))
			}
		case "float64":
			i, err := strconv.ParseFloat(field.Tag.Get("default"), 64)
			if err == nil {
				fs.Float64(flagName, i, field.Tag.Get("desc"))
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
			v.Field(i).SetFloat(viper.GetFloat64(field.Tag.Get("name")))
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
func NewConfig(c Configer, merge ...bool) Configer {
	fs := pflag.NewFlagSet("config", pflag.ContinueOnError)
	defineFlagsFromStruct(reflect.TypeOf(c).Elem(), fs)
	// Only merge local flagset into global command line if requested
	shouldMerge := true // Default to true to maintain original behavior
	if len(merge) > 0 {
		shouldMerge = merge[0]
	}
	if shouldMerge {
		pflag.CommandLine.AddFlagSet(fs)
	}
	c.generate()
	setPropertiesFromFlags(reflect.ValueOf(c), c.getParser())
	return c
}

// NewConfigWithFlagSet generates a new configuration setup with a custom flagset
// This is useful for testing or when you want to use a specific flagset
func NewConfigWithFlagSet(c Configer, fs *pflag.FlagSet) Configer {
	defineFlagsFromStruct(reflect.TypeOf(c).Elem(), fs)
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

// CreateViperWithFlagSet creates a parser instance with a custom flagset
// This is useful for testing
func CreateViperWithFlagSet(fs *pflag.FlagSet) (v *viper.Viper) {
	v = viper.New()
	v.AutomaticEnv()
	fs.Parse([]string{}) // Parse with empty args for testing
	v.BindPFlags(fs)
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
