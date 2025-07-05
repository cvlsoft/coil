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
	// Test that we can access MyCustomConfig fields
	if cfg.FooBar == "static" {
		t.Error("FooBar field is not accessible or not initialized")
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
		t.Errorf("FooBar should be set from environment variable: got %q, want %q", cfg.FooBar, testValue)
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
