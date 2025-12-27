# Coil Architecture

## Overview

Coil is a lightweight Go configuration management library built on top of [Viper](https://github.com/spf13/viper) and [pflag](https://github.com/spf13/pflag). It provides a struct-based approach to configuration management, automatically generating CLI flags and environment variable bindings from Go struct tags.

## Core Design Principles

1. **Struct-First Configuration**: Define configurations as Go structs with tags, eliminating manual flag registration
2. **Composability**: Layer multiple configuration structs together through embedding
3. **Multiple Input Sources**: Support CLI flags, environment variables, and config files with automatic precedence handling
4. **Type Safety**: Leverage Go's type system for configuration validation at compile time
5. **Zero Boilerplate**: Minimize repetitive code through reflection-based automation

## Architecture Components

### 1. Core Interface: `Configer`

```go
type Configer interface {
    generate()
    getParser() *viper.Viper
}
```

The `Configer` interface defines the contract for all configuration types. It ensures:
- Configuration initialization via `generate()`
- Access to the underlying Viper instance via `getParser()`

### 2. Base Type: `Config`

The `Config` struct is the foundation that all configuration types must embed:

```go
type Config struct {
    viper *viper.Viper
}
```

**Key Features:**
- Holds the Viper instance for configuration parsing
- Implements the `Configer` interface
- Provides the `HasConfig()` method for type introspection at runtime

**Location**: `coil.go:16-50`

### 3. Configuration Factory: `NewConfig()`

**Signature**: `func NewConfig(c Configer, merge ...bool) Configer`

**Purpose**: Creates and initializes configuration instances through reflection

**Process Flow**:
1. Creates a pflag.FlagSet for CLI flags
2. Recursively discovers struct fields via `defineFlagsFromStruct()`
3. Optionally merges flags into global CommandLine
4. Calls `generate()` to initialize Viper
5. Binds configuration values via `setPropertiesFromFlags()`
6. Returns the initialized configuration

**Location**: `coil.go:187-201`

### 4. Struct Tag System

Coil uses struct tags to define configuration metadata:

```go
type Example struct {
    Field string `type:"string" name:"field_name" default:"value" desc:"Description"`
}
```

**Supported Tags**:
- `type`: Data type (string, int, bool, float32, float64, duration, []string)
- `name`: CLI flag and config file key name
- `default`: Default value when not provided
- `desc`: Human-readable description for help text
- `prefix`: Namespace prefix for nested configurations

**Location**: `coil.go:69-134` (defineFlagsFromStruct)

### 5. Prefix System

The prefix system enables multiple instances of the same configuration type without naming conflicts:

```go
type Config struct {
    PrimaryDB DatabaseConfig `prefix:"primary"`
    ReplicaDB DatabaseConfig `prefix:"replica"`
}
```

**How It Works**:
- Prefixes are applied recursively during flag definition
- Environment variables: `PRIMARY_DBHOST`, `REPLICA_DBHOST`
- CLI flags: `--primary_dbhost`, `--replica_dbhost`
- Nested prefixes combine: `outer_inner_field`

**Location**: `coil.go:80-93`, `coil.go:152-164`

### 6. Reflection Engine

Coil uses Go's `reflect` package for two main operations:

#### a. Flag Definition (`defineFlagsFromStructWithPrefix`)
- Recursively traverses struct fields
- Detects prefix tags and applies them
- Creates appropriate pflag types based on field types
- Handles nested structs automatically

**Location**: `coil.go:69-134`

#### b. Property Binding (`setPropertiesFromFlagsWithPrefix`)
- Recursively traverses struct instances
- Reads values from Viper
- Sets struct field values using reflection
- Applies defaults when values are not set
- Calls `Parse()` method if it exists

**Location**: `coil.go:137-184`

### 7. Viper Integration

**Creation**: `CreateViper()` function
**Process**:
1. Creates new Viper instance
2. Enables automatic environment variable binding
3. Parses command-line flags
4. Binds flags to Viper
5. Loads config file if `--config` flag is provided

**Location**: `coil.go:210-227`

### 8. Pre-built Configurations

Coil provides ready-to-use configuration types in `configs.go`:

#### `APIServiceConfig`
Common API service settings:
- Version, Name, Build
- Host, Port, URL
- Timeout duration

#### `DatabaseConfig`
Standard database connection parameters:
- Host, Port, Name
- User, Password
- SSL mode, Debug flag

#### `LogConfig`
Comprehensive logging configuration:
- Level, Format, Output
- File rotation (MaxSize, MaxBackups, MaxAge)
- Static fields, Service metadata

**Location**: `configs.go`

## Configuration Precedence

Coil follows Viper's precedence order (highest to lowest):

1. **CLI Flags**: `--flag=value`
2. **Environment Variables**: `VARIABLE_NAME=value`
3. **Config File**: YAML/JSON/TOML files
4. **Default Values**: From struct tags

## Data Flow

```
User Input (CLI/Env/File)
         ↓
    Viper Parser
         ↓
   Reflection Engine
         ↓
  Struct Population
         ↓
   Application Code
```

### Initialization Flow

1. User calls `NewConfig(&Config{})`
2. Struct type is analyzed via reflection
3. Flags are defined from struct tags
4. Viper instance is created and initialized
5. Environment variables are bound automatically
6. Config file is loaded if specified
7. Values are set on struct fields via reflection
8. Configured struct is returned

## Type Support

### Supported Types
- `string`: Text values
- `[]string`: String slices (comma-separated)
- `int`: Integer values (stored as int64 internally)
- `bool`: Boolean flags
- `float32`: 32-bit floating point
- `float64`: 64-bit floating point
- `duration`: Time durations (e.g., "10s", "5m")

### Type Conversion
Automatic conversion happens at the Viper level:
- Environment variables are parsed based on flag type
- Config file values are unmarshaled to appropriate types
- Default values are parsed from string tags

## Advanced Features

### 1. Config Type Introspection

The `HasConfig()` method allows runtime checking:

```go
if cfg.HasConfig(DatabaseConfig{}) {
    // Database config is embedded
}
```

**Location**: `coil.go:29-42`

### 2. Custom FlagSet Support

For testing or custom scenarios:

```go
func NewConfigWithFlagSet(c Configer, fs *pflag.FlagSet) Configer
```

Allows using a specific FlagSet instead of the global one.

**Location**: `coil.go:203-208`

### 3. Merge Control

Control whether flags merge into global CommandLine:

```go
NewConfig(&Config{}, false) // Don't merge
NewConfig(&Config{}, true)  // Merge (default)
```

### 4. Parse Method Hook

Structs can implement a `Parse(v *viper.Viper)` method for custom post-processing:

```go
func (c *MyConfig) Parse(v *viper.Viper) {
    // Custom validation or transformation
}
```

**Location**: `coil.go:175-178`

## Testing Strategy

The test suite (`coil_test.go`) validates:

1. **Basic Functionality**: Config creation, struct access
2. **Environment Variables**: Override defaults, case sensitivity
3. **Prefix System**: Single and nested prefixes, collision avoidance
4. **Type Support**: All supported types with defaults and env vars
5. **Nested Structs**: Deep nesting, recursive processing
6. **Edge Cases**: Empty defaults, missing tags, mixed scenarios

**Benchmark Tests**: Measure performance of config creation for different complexity levels

## Extension Points

### Adding New Pre-built Configs

Create new structs in `configs.go`:

```go
type NewConfig struct {
    Field string `type:"string" name:"field" default:"value" desc:"Description"`
}
```

### Adding New Type Support

Extend the switch statements in:
1. `defineFlagsFromStructWithPrefix()` for flag definition
2. `setPropertiesFromFlagsWithPrefix()` for value binding

## Dependencies

- **github.com/spf13/viper**: Configuration parsing and management
- **github.com/spf13/pflag**: POSIX/GNU-style command-line flags

**Indirect Dependencies**:
- File format parsers (YAML, TOML, JSON)
- File system abstraction (afero)
- Type casting utilities

## Performance Considerations

1. **Reflection Overhead**: Occurs only during initialization, not runtime
2. **Caching**: Viper caches parsed values
3. **Memory**: Minimal overhead per configuration instance
4. **Benchmarks**: See `coil_test.go` for performance metrics

## Security Considerations

1. **Sensitive Data**: Config file paths should be validated
2. **Environment Variables**: Automatically bound, consider sensitive data exposure
3. **Default Values**: Visible in help text and code
4. **File Permissions**: Config files may contain credentials

**Best Practices**:
- Use environment variables for secrets in production
- Restrict config file permissions
- Don't commit config files with secrets to version control
- Use `default:""` for sensitive fields

## Future Extensibility

The architecture supports:
- Custom validators per field
- Conditional configuration based on environment
- Dynamic reloading of configuration
- Configuration inheritance and overrides
- Plugin-based configuration sources

## Usage Patterns

### Simple Application

```go
type Config struct {
    coil.Config
    coil.APIServiceConfig
}

cfg := coil.NewConfig(&Config{})
```

### Multi-Database Application

```go
type Config struct {
    coil.Config
    PrimaryDB coil.DatabaseConfig `prefix:"primary"`
    ReplicaDB coil.DatabaseConfig `prefix:"replica"`
}
```

### Complex Application

```go
type Config struct {
    coil.Config
    coil.APIServiceConfig
    coil.LogConfig
    DB coil.DatabaseConfig `prefix:"db"`
    Cache CacheConfig `prefix:"cache"`
    CustomSettings
}
```

## Error Handling

- **Missing Config File**: Panics with "Could not find configuration file"
- **Invalid Config File**: Panics with "Could not parse configuration file"
- **Invalid Flags**: Handled by pflag (prints error and exits)
- **Type Mismatches**: Viper attempts conversion, may return zero values

## Conclusion

Coil's architecture leverages Go's type system and reflection to provide a declarative, composable approach to configuration management. By building on Viper and pflag, it inherits their robustness while reducing boilerplate and improving developer experience.
