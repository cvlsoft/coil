package coil

import "time"

// APIServiceConfig is a global struct passed to all services
type APIServiceConfig struct {
	Version string        `type:"string" name:"version" default:"1.0.0" desc:"API version (follows semver)"`
	Name    string        `type:"string" name:"name" default:"service-api" desc:"Default name of the service"`
	Build   string        `type:"string" name:"build" default:"UNSPECIFIED" desc:"Build version"`
	Host    string        `type:"string" name:"host" default:"localhost" desc:"Server hostname to bind to"`
	URL     string        `type:"string" name:"api_url" default:"" desc:"The URL to the API"`
	Port    int           `type:"int" name:"port" default:"80" desc:"Server port to bind to"`
	Timeout time.Duration `type:"duration" name:"timeout" default:"15s" desc:"Timeout for any connection i.e. 10s"`
}

// DatabaseConfig represents a composable struct for db connections
type DatabaseConfig struct {
	DBHost  string `type:"string" name:"dbhost" default:"localhost" desc:"Database hostname"`
	DBUser  string `type:"string" name:"dbuser" default:"" desc:"Database username"`
	DBName  string `type:"string" name:"dbname" default:"" desc:"Database name"`
	DBPass  string `type:"string" name:"dbpass" default:"" desc:"Database password"`
	DBSSL   string `type:"string" name:"dbssl" default:"disable" desc:"Database SSL mode"`
	DBDebug bool   `type:"string" name:"dbdebug" default:"" desc:"Enable database debug mode"`
	DBPort  int    `type:"int" name:"dbport" default:"5432" desc:"Database port number"`
}

// LogConfig represents a composable struct for logging
type LogConfig struct {
	// Core logging settings
	Level  string `type:"string" name:"log_level" default:"info" desc:"Log level (trace, debug, info, warn, error, fatal)"`
	Format string `type:"string" name:"log_format" default:"json" desc:"Log format (json, text, logfmt)"`

	// Output configuration
	Output     string `type:"string" name:"log_output" default:"stdout" desc:"Log output destination (stdout, stderr, file)"`
	FilePath   string `type:"string" name:"log_file_path" default:"./logs/app.log" desc:"Path to log file when output is 'file'"`
	MaxSize    int    `type:"int" name:"log_max_size" default:"100" desc:"Maximum size in megabytes before rotation"`
	MaxBackups int    `type:"int" name:"log_max_backups" default:"3" desc:"Maximum number of old log files to retain"`
	MaxAge     int    `type:"int" name:"log_max_age" default:"28" desc:"Maximum number of days to retain old log files"`
	Compress   bool   `type:"bool" name:"log_compress" default:"false" desc:"Whether to compress rotated log files"`

	// Field configuration
	StaticFields string `type:"string" name:"log_static_fields" default:"" desc:"Static fields to include in all logs (JSON format)"`
	ServiceName  string `type:"string" name:"log_service_name" default:"" desc:"Service name to include in logs"`
	Environment  string `type:"string" name:"log_environment" default:"" desc:"Environment name (dev, staging, prod)"`
	InstanceID   string `type:"string" name:"log_instance_id" default:"" desc:"Instance/container ID to include in logs"`
}
