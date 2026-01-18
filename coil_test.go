package coil

import (
	"os"
	"testing"
)

// Config represents your app's local config
type ConfigTest1 struct {
	Config
	MyCustomConfig
}

// MyCustomConfig represents a custom configuration
type MyCustomConfig struct {
	FooBar string `type:"string" name:"foo_bar" default:"static" desc:"Foo bar value"`
}

// NewConfig is a factory generator for your configuration
func NewConfigTest() *ConfigTest1 {
	cfg := NewConfig(&ConfigTest1{}, false)
	return cfg.(*ConfigTest1)
}

func TestNewConfig(t *testing.T) {
	// Test that NewConfig returns a non-nil Config pointer
	cfg := NewConfigTest()
	if cfg == nil {
		t.Fatal("NewConfig() returned nil")
	}
	// The cfg is already *Config, no need for type assertion
	// Just verify it's the correct type
	var _ *ConfigTest1 = cfg // This will fail to compile if cfg is not *Config
}

func TestConfigStructure(t *testing.T) {
	// Test that Config embeds all required structs
	cfg := NewConfigTest()
	// Debug: Print the type and value
	t.Logf("Config type: %T", cfg)
	t.Logf("Config value: %+v", cfg)
	// Test that we can access MyCustomConfig fields and it has the default
	// value
	if cfg.FooBar != "static" {
		t.Errorf(
			"FooBar field should have default value 'static', got %q",
			cfg.FooBar,
		)
	}
}

func TestConfigWithEnvironmentVariable(t *testing.T) {
	// Save original env var if it exists
	originalValue := os.Getenv("FOO_BAR")
	defer func() {
		if originalValue != "" {
			os.Setenv("FOO_BAR", originalValue)
		} else {
			os.Unsetenv("FOO_BAR")
		}
	}()
	// Set test environment variable
	testValue := "from_env"
	os.Setenv("FOO_BAR", testValue)
	// Create new config
	cfg := NewConfigTest()
	// Check if environment variable overrides default
	if cfg.FooBar != testValue {
		t.Errorf(
			"FooBar should be set from environment variable: got %q, want %q",
			cfg.FooBar,
			testValue,
		)
	}
}

func TestConfigTypeAssertion(t *testing.T) {
	cfg := NewConfigTest()
	// Test that all embedded types are accessible
	tests := []struct {
		name string
		test func() bool
	}{
		{
			name: "Config has MyCustomConfig",
			test: func() bool {
				_ = cfg.MyCustomConfig
				return true
			},
		},
		{
			name: "Config has FooBar field",
			test: func() bool {
				_ = cfg.FooBar
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer func() {
				if r := recover(); r != nil {
					t.Errorf("Accessing field caused panic: %v", r)
				}
			}()

			if !tt.test() {
				t.Errorf("Test %s failed", tt.name)
			}
		})
	}
}

func TestMultipleConfigInstances(t *testing.T) {
	// Test that multiple instances can be created independently
	cfg1 := NewConfigTest()
	cfg2 := NewConfigTest()
	if cfg1 == cfg2 {
		t.Error("NewConfig() should create independent instances")
	}
	// Both should have the same default values
	if cfg1.FooBar != cfg2.FooBar {
		t.Error("Independent configs should have same default values")
	}
}

func BenchmarkNewConfig(b *testing.B) {
	// Benchmark config creation
	for b.Loop() {
		_ = NewConfigTest()
	}
}

// ConfigWithPrefix tests the prefix functionality for avoiding collisions
type ConfigWithPrefix struct {
	Config
	PrimaryDB DatabaseConfig `prefix:"primary"`
	ReplicaDB DatabaseConfig `prefix:"replica"`
}

// NewConfigWithPrefix is a factory generator for prefix testing
func NewConfigWithPrefix() *ConfigWithPrefix {
	cfg := NewConfig(&ConfigWithPrefix{}, false)
	return cfg.(*ConfigWithPrefix)
}

func TestConfigWithPrefix(t *testing.T) {
	// Save original env vars
	origPrimaryHost := os.Getenv("PRIMARY_DBHOST")
	origReplicaHost := os.Getenv("REPLICA_DBHOST")
	origPrimaryPort := os.Getenv("PRIMARY_DBPORT")
	origReplicaPort := os.Getenv("REPLICA_DBPORT")
	defer func() {
		restoreEnv("PRIMARY_DBHOST", origPrimaryHost)
		restoreEnv("REPLICA_DBHOST", origReplicaHost)
		restoreEnv("PRIMARY_DBPORT", origPrimaryPort)
		restoreEnv("REPLICA_DBPORT", origReplicaPort)
	}()

	// Set different values for primary and replica
	os.Setenv("PRIMARY_DBHOST", "primary-host.example.com")
	os.Setenv("REPLICA_DBHOST", "replica-host.example.com")
	os.Setenv("PRIMARY_DBPORT", "5433")
	os.Setenv("REPLICA_DBPORT", "5434")

	cfg := NewConfigWithPrefix()

	// Verify primary DB config
	if cfg.PrimaryDB.DBHost != "primary-host.example.com" {
		t.Errorf(
			"PrimaryDB.DBHost = %q, want %q",
			cfg.PrimaryDB.DBHost,
			"primary-host.example.com",
		)
	}
	if cfg.PrimaryDB.DBPort != 5433 {
		t.Errorf("PrimaryDB.DBPort = %d, want %d", cfg.PrimaryDB.DBPort, 5433)
	}

	// Verify replica DB config
	if cfg.ReplicaDB.DBHost != "replica-host.example.com" {
		t.Errorf(
			"ReplicaDB.DBHost = %q, want %q",
			cfg.ReplicaDB.DBHost,
			"replica-host.example.com",
		)
	}
	if cfg.ReplicaDB.DBPort != 5434 {
		t.Errorf("ReplicaDB.DBPort = %d, want %d", cfg.ReplicaDB.DBPort, 5434)
	}
}

func TestConfigWithPrefixDefaults(t *testing.T) {
	// Clear any existing env vars that might interfere
	origPrimaryHost := os.Getenv("PRIMARY_DBHOST")
	origReplicaHost := os.Getenv("REPLICA_DBHOST")
	defer func() {
		restoreEnv("PRIMARY_DBHOST", origPrimaryHost)
		restoreEnv("REPLICA_DBHOST", origReplicaHost)
	}()
	os.Unsetenv("PRIMARY_DBHOST")
	os.Unsetenv("REPLICA_DBHOST")

	cfg := NewConfigWithPrefix()

	// Both should have the default value when no env var is set
	if cfg.PrimaryDB.DBHost != "localhost" {
		t.Errorf(
			"PrimaryDB.DBHost = %q, want default %q",
			cfg.PrimaryDB.DBHost,
			"localhost",
		)
	}
	if cfg.ReplicaDB.DBHost != "localhost" {
		t.Errorf(
			"ReplicaDB.DBHost = %q, want default %q",
			cfg.ReplicaDB.DBHost,
			"localhost",
		)
	}
}

func restoreEnv(key, value string) {
	if value != "" {
		os.Setenv(key, value)
	} else {
		os.Unsetenv(key)
	}
}

// AllTypesConfig tests all supported field types
type AllTypesConfig struct {
	Config
	TypesStruct AllTypesStruct
}

type AllTypesStruct struct {
	StringField  string  `type:"string"  name:"test_string"     default:"default_string" desc:"A string field"`
	IntField     int     `type:"int"     name:"test_int"        default:"42"             desc:"An int field"`
	BoolFieldT   bool    `type:"bool"    name:"test_bool_true"  default:"true"           desc:"A bool field defaulting to true"`
	BoolFieldF   bool    `type:"bool"    name:"test_bool_false" default:"false"          desc:"A bool field defaulting to false"`
	Float32Field float32 `type:"float32" name:"test_float32"    default:"3.14"           desc:"A float32 field"`
	Float64Field float64 `type:"float64" name:"test_float64"    default:"2.718281828"    desc:"A float64 field"`
}

func NewAllTypesConfig() *AllTypesConfig {
	cfg := NewConfig(&AllTypesConfig{}, false)
	return cfg.(*AllTypesConfig)
}

func TestAllFieldTypeDefaults(t *testing.T) {
	// Clear any env vars that might interfere
	envVars := []string{
		"TEST_STRING",
		"TEST_INT",
		"TEST_BOOL_TRUE",
		"TEST_BOOL_FALSE",
		"TEST_FLOAT32",
		"TEST_FLOAT64",
	}
	origVals := make(map[string]string)
	for _, env := range envVars {
		origVals[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for _, env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewAllTypesConfig()

	// Test string default
	if cfg.TypesStruct.StringField != "default_string" {
		t.Errorf(
			"StringField = %q, want %q",
			cfg.TypesStruct.StringField,
			"default_string",
		)
	}

	// Test int default
	if cfg.TypesStruct.IntField != 42 {
		t.Errorf("IntField = %d, want %d", cfg.TypesStruct.IntField, 42)
	}

	// Test bool defaults
	if cfg.TypesStruct.BoolFieldT != true {
		t.Errorf("BoolFieldT = %v, want %v", cfg.TypesStruct.BoolFieldT, true)
	}
	if cfg.TypesStruct.BoolFieldF != false {
		t.Errorf("BoolFieldF = %v, want %v", cfg.TypesStruct.BoolFieldF, false)
	}

	// Test float32 default
	if cfg.TypesStruct.Float32Field < 3.13 ||
		cfg.TypesStruct.Float32Field > 3.15 {
		t.Errorf(
			"Float32Field = %f, want approximately %f",
			cfg.TypesStruct.Float32Field,
			3.14,
		)
	}

	// Test float64 default
	if cfg.TypesStruct.Float64Field < 2.71 ||
		cfg.TypesStruct.Float64Field > 2.72 {
		t.Errorf(
			"Float64Field = %f, want approximately %f",
			cfg.TypesStruct.Float64Field,
			2.718281828,
		)
	}
}

func TestAllFieldTypesFromEnv(t *testing.T) {
	// Save and set env vars
	envVars := map[string]string{
		"TEST_STRING":     "env_string",
		"TEST_INT":        "100",
		"TEST_BOOL_TRUE":  "false",
		"TEST_BOOL_FALSE": "true",
		"TEST_FLOAT32":    "1.5",
		"TEST_FLOAT64":    "9.99",
	}
	origVals := make(map[string]string)
	for env, val := range envVars {
		origVals[env] = os.Getenv(env)
		os.Setenv(env, val)
	}
	defer func() {
		for env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewAllTypesConfig()

	if cfg.TypesStruct.StringField != "env_string" {
		t.Errorf(
			"StringField = %q, want %q",
			cfg.TypesStruct.StringField,
			"env_string",
		)
	}
	if cfg.TypesStruct.IntField != 100 {
		t.Errorf("IntField = %d, want %d", cfg.TypesStruct.IntField, 100)
	}
	if cfg.TypesStruct.BoolFieldT != false {
		t.Errorf("BoolFieldT = %v, want %v", cfg.TypesStruct.BoolFieldT, false)
	}
	if cfg.TypesStruct.BoolFieldF != true {
		t.Errorf("BoolFieldF = %v, want %v", cfg.TypesStruct.BoolFieldF, true)
	}
	if cfg.TypesStruct.Float32Field < 1.4 ||
		cfg.TypesStruct.Float32Field > 1.6 {
		t.Errorf(
			"Float32Field = %f, want approximately %f",
			cfg.TypesStruct.Float32Field,
			1.5,
		)
	}
	if cfg.TypesStruct.Float64Field < 9.98 ||
		cfg.TypesStruct.Float64Field > 10.0 {
		t.Errorf(
			"Float64Field = %f, want approximately %f",
			cfg.TypesStruct.Float64Field,
			9.99,
		)
	}
}

// NestedConfig tests deeply nested struct support
type NestedConfig struct {
	Config
	Level1 Level1Struct
}

type Level1Struct struct {
	L1Field string `type:"string" name:"l1_field" default:"level1" desc:"Level 1 field"`
	Level2  Level2Struct
}

type Level2Struct struct {
	L2Field string `type:"string" name:"l2_field" default:"level2" desc:"Level 2 field"`
	Level3  Level3Struct
}

type Level3Struct struct {
	L3Field string `type:"string" name:"l3_field" default:"level3" desc:"Level 3 field"`
}

func NewNestedConfig() *NestedConfig {
	cfg := NewConfig(&NestedConfig{}, false)
	return cfg.(*NestedConfig)
}

func TestNestedStructDefaults(t *testing.T) {
	envVars := []string{"L1_FIELD", "L2_FIELD", "L3_FIELD"}
	origVals := make(map[string]string)
	for _, env := range envVars {
		origVals[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for _, env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewNestedConfig()

	if cfg.Level1.L1Field != "level1" {
		t.Errorf("L1Field = %q, want %q", cfg.Level1.L1Field, "level1")
	}
	if cfg.Level1.Level2.L2Field != "level2" {
		t.Errorf("L2Field = %q, want %q", cfg.Level1.Level2.L2Field, "level2")
	}
	if cfg.Level1.Level2.Level3.L3Field != "level3" {
		t.Errorf(
			"L3Field = %q, want %q",
			cfg.Level1.Level2.Level3.L3Field,
			"level3",
		)
	}
}

func TestNestedStructFromEnv(t *testing.T) {
	envVars := map[string]string{
		"L1_FIELD": "env_level1",
		"L2_FIELD": "env_level2",
		"L3_FIELD": "env_level3",
	}
	origVals := make(map[string]string)
	for env, val := range envVars {
		origVals[env] = os.Getenv(env)
		os.Setenv(env, val)
	}
	defer func() {
		for env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewNestedConfig()

	if cfg.Level1.L1Field != "env_level1" {
		t.Errorf("L1Field = %q, want %q", cfg.Level1.L1Field, "env_level1")
	}
	if cfg.Level1.Level2.L2Field != "env_level2" {
		t.Errorf(
			"L2Field = %q, want %q",
			cfg.Level1.Level2.L2Field,
			"env_level2",
		)
	}
	if cfg.Level1.Level2.Level3.L3Field != "env_level3" {
		t.Errorf(
			"L3Field = %q, want %q",
			cfg.Level1.Level2.Level3.L3Field,
			"env_level3",
		)
	}
}

// NestedPrefixConfig tests nested prefixes
type NestedPrefixConfig struct {
	Config
	Outer OuterStruct `prefix:"outer"`
}

type OuterStruct struct {
	OuterField string      `type:"string" name:"field" default:"outer_default" desc:"Outer field"`
	Inner      InnerStruct `                                                                      prefix:"inner"`
}

type InnerStruct struct {
	InnerField string `type:"string" name:"field" default:"inner_default" desc:"Inner field"`
}

func NewNestedPrefixConfig() *NestedPrefixConfig {
	cfg := NewConfig(&NestedPrefixConfig{}, false)
	return cfg.(*NestedPrefixConfig)
}

func TestNestedPrefixes(t *testing.T) {
	// Test that nested prefixes combine correctly: outer_field and
	// outer_inner_field
	envVars := map[string]string{
		"OUTER_FIELD":       "outer_from_env",
		"OUTER_INNER_FIELD": "inner_from_env",
	}
	origVals := make(map[string]string)
	for env, val := range envVars {
		origVals[env] = os.Getenv(env)
		os.Setenv(env, val)
	}
	defer func() {
		for env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewNestedPrefixConfig()

	if cfg.Outer.OuterField != "outer_from_env" {
		t.Errorf(
			"OuterField = %q, want %q",
			cfg.Outer.OuterField,
			"outer_from_env",
		)
	}
	if cfg.Outer.Inner.InnerField != "inner_from_env" {
		t.Errorf(
			"InnerField = %q, want %q",
			cfg.Outer.Inner.InnerField,
			"inner_from_env",
		)
	}
}

func TestNestedPrefixDefaults(t *testing.T) {
	envVars := []string{"OUTER_FIELD", "OUTER_INNER_FIELD"}
	origVals := make(map[string]string)
	for _, env := range envVars {
		origVals[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for _, env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewNestedPrefixConfig()

	if cfg.Outer.OuterField != "outer_default" {
		t.Errorf(
			"OuterField = %q, want default %q",
			cfg.Outer.OuterField,
			"outer_default",
		)
	}
	if cfg.Outer.Inner.InnerField != "inner_default" {
		t.Errorf(
			"InnerField = %q, want default %q",
			cfg.Outer.Inner.InnerField,
			"inner_default",
		)
	}
}

// NoTagConfig tests fields without tags are ignored
type NoTagConfig struct {
	Config
	NoTagStruct NoTagStruct
}

type NoTagStruct struct {
	WithTag    string `type:"string" name:"with_tag" default:"tagged" desc:"Has tags"`
	WithoutTag string // No tags - should remain zero value
}

func NewNoTagConfig() *NoTagConfig {
	cfg := NewConfig(&NoTagConfig{}, false)
	return cfg.(*NoTagConfig)
}

func TestFieldsWithoutTagsIgnored(t *testing.T) {
	origVal := os.Getenv("WITH_TAG")
	os.Unsetenv("WITH_TAG")
	defer restoreEnv("WITH_TAG", origVal)

	cfg := NewNoTagConfig()

	if cfg.NoTagStruct.WithTag != "tagged" {
		t.Errorf("WithTag = %q, want %q", cfg.NoTagStruct.WithTag, "tagged")
	}
	if cfg.NoTagStruct.WithoutTag != "" {
		t.Errorf(
			"WithoutTag = %q, want empty string (zero value)",
			cfg.NoTagStruct.WithoutTag,
		)
	}
}

// MixedPrefixConfig tests mixing prefixed and non-prefixed structs
type MixedPrefixConfig struct {
	Config
	Regular  RegularStruct
	Prefixed RegularStruct `prefix:"prefixed"`
}

type RegularStruct struct {
	Value string `type:"string" name:"value" default:"default_val" desc:"A value"`
}

func NewMixedPrefixConfig() *MixedPrefixConfig {
	cfg := NewConfig(&MixedPrefixConfig{}, false)
	return cfg.(*MixedPrefixConfig)
}

func TestMixedPrefixAndNonPrefix(t *testing.T) {
	envVars := map[string]string{
		"VALUE":          "regular_env",
		"PREFIXED_VALUE": "prefixed_env",
	}
	origVals := make(map[string]string)
	for env, val := range envVars {
		origVals[env] = os.Getenv(env)
		os.Setenv(env, val)
	}
	defer func() {
		for env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewMixedPrefixConfig()

	if cfg.Regular.Value != "regular_env" {
		t.Errorf(
			"Regular.Value = %q, want %q",
			cfg.Regular.Value,
			"regular_env",
		)
	}
	if cfg.Prefixed.Value != "prefixed_env" {
		t.Errorf(
			"Prefixed.Value = %q, want %q",
			cfg.Prefixed.Value,
			"prefixed_env",
		)
	}
}

// EmptyDefaultConfig tests empty default values
type EmptyDefaultConfig struct {
	Config
	EmptyDefaults EmptyDefaultsStruct
}

type EmptyDefaultsStruct struct {
	EmptyString string `type:"string" name:"empty_string" default:""  desc:"Empty string default"`
	EmptyInt    int    `type:"int"    name:"empty_int"    default:"0" desc:"Zero int default"`
}

func NewEmptyDefaultConfig() *EmptyDefaultConfig {
	cfg := NewConfig(&EmptyDefaultConfig{}, false)
	return cfg.(*EmptyDefaultConfig)
}

func TestEmptyDefaults(t *testing.T) {
	envVars := []string{"EMPTY_STRING", "EMPTY_INT"}
	origVals := make(map[string]string)
	for _, env := range envVars {
		origVals[env] = os.Getenv(env)
		os.Unsetenv(env)
	}
	defer func() {
		for _, env := range envVars {
			restoreEnv(env, origVals[env])
		}
	}()

	cfg := NewEmptyDefaultConfig()

	if cfg.EmptyDefaults.EmptyString != "" {
		t.Errorf(
			"EmptyString = %q, want empty string",
			cfg.EmptyDefaults.EmptyString,
		)
	}
	if cfg.EmptyDefaults.EmptyInt != 0 {
		t.Errorf("EmptyInt = %d, want 0", cfg.EmptyDefaults.EmptyInt)
	}
}

// Test getParser returns the viper instance
func TestGetParser(t *testing.T) {
	cfg := NewConfigTest()
	parser := cfg.Config.getParser()
	if parser == nil {
		t.Error("getParser() returned nil, expected viper instance")
	}
}

// SimpleCfg for merge testing
type SimpleCfg struct {
	Config
	Simple SimpleStruct
}

type SimpleStruct struct {
	Field string `type:"string" name:"merge_test_field" default:"merge_default" desc:"Merge test"`
}

// Test NewConfig with merge=true (default behavior)
func TestNewConfigWithMerge(t *testing.T) {
	origVal := os.Getenv("MERGE_TEST_FIELD")
	os.Unsetenv("MERGE_TEST_FIELD")
	defer restoreEnv("MERGE_TEST_FIELD", origVal)

	// Test with explicit merge=true
	cfg := NewConfig(&SimpleCfg{}, true)
	simpleCfg := cfg.(*SimpleCfg)
	if simpleCfg.Simple.Field != "merge_default" {
		t.Errorf("Field = %q, want %q", simpleCfg.Simple.Field, "merge_default")
	}
}

// CaseCfg for case sensitivity testing
type CaseCfg struct {
	Config
	Case CaseStruct
}

type CaseStruct struct {
	MixedCase string `type:"string" name:"mixed_case_field" default:"default" desc:"Mixed case field"`
}

// Test case sensitivity of environment variables
func TestEnvVarCaseSensitivity(t *testing.T) {
	// Viper uses uppercase env vars by default
	origVal := os.Getenv("MIXED_CASE_FIELD")
	os.Setenv("MIXED_CASE_FIELD", "uppercase_env")
	defer restoreEnv("MIXED_CASE_FIELD", origVal)

	cfg := NewConfig(&CaseCfg{}, false)
	caseCfg := cfg.(*CaseCfg)
	if caseCfg.Case.MixedCase != "uppercase_env" {
		t.Errorf(
			"MixedCase = %q, want %q",
			caseCfg.Case.MixedCase,
			"uppercase_env",
		)
	}
}

// Benchmark for prefix config creation
func BenchmarkNewConfigWithPrefix(b *testing.B) {
	for b.Loop() {
		_ = NewConfigWithPrefix()
	}
}

// Benchmark for nested config creation
func BenchmarkNewNestedConfig(b *testing.B) {
	for b.Loop() {
		_ = NewNestedConfig()
	}
}

// Benchmark for all types config creation
func BenchmarkNewAllTypesConfig(b *testing.B) {
	for b.Loop() {
		_ = NewAllTypesConfig()
	}
}
