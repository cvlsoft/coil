package coil

// APIServiceConfig is a global struct passed to all services
type APIServiceConfig struct {
	Version string `type:"string" name:"version" default:"1.0.0" desc:"API version (follows semver)"`
	Name    string `type:"string" name:"name" default:"service-api" desc:"Default name of the service"`
	Build   string `type:"string" name:"build" default:"UNSPECIFIED" desc:"Build version"`
	Host    string `type:"string" name:"host" default:"localhost" desc:"Server hostname to bind to"`
	Port    int    `type:"int" name:"port" default:"80" desc:"Server port to bind to"`
	APIURL  string `type:"string" name:"api_url" default:"www" desc:"The URL to the API"`
}

// DatabaseConfig represents a composable struct for db connections
type DatabaseConfig struct {
	DBHost  string `type:"string" name:"dbhost" default:"localhost" desc:"Database hostname"`
	DBUser  string `type:"string" name:"dbuser" default:"" desc:"Database username"`
	DBName  string `type:"string" name:"dbname" default:"" desc:"Database name"`
	DBPass  string `type:"string" name:"dbpass" default:"" desc:"Database password"`
	DBSSL   string `type:"string" name:"dbssl" default:"disable" desc:"Database SSL mode"`
	DBPort  int    `type:"int" name:"dbport" default:"5432" desc:"Database port number"`
	DBDebug bool   `type:"string" name:"dbdebug" default:"" desc:"Enable database debug mode"`
}
