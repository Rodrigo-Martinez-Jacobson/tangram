// Package config models the application configuration.
// The configuration values will be loaded from command arguments,
// config file, environment variables and default values.
// The prevalence order is (more to less prevalence):
// - command arguments
// - config file
// - environment
// - default values
package conf

import (
	"os"
	"strconv"
	"time"
)

const (
	defaultConfigFile = "tangram.toml"
)

// Config contains application configuration
type Config struct {
	addr            string
	shutdownTimeout time.Duration
}

// Address returns the HTTP server address.
// This format have the format "[inet-address]:port"
func (c Config) Address() string {
	return c.addr
}

// ShutdownTimeout returns the application shutdown timeout to wait to
// shutdown the HTTP server for graceful shutdown
func (c Config) ShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}

// Load applcation configuration from default config file
func Load() (c Config, err error) {
	return loadConfig(defaultConfigFile)
}

func loadConfig(file string) (c Config, err error) {
	c = Config{}
	c.addr = confValue(confDefs["address"])
	c.shutdownTimeout = asDuration(confValue(confDefs["shutdownTimeout"]))
	return
}

// confValue returns a configuration value from best-fit source.
// Parameters:
// - arg: the name of configuration value in the command line argument
// - conf: the name of configuration value in configuration file
// - env: the name of configuration value in the environment
// - def: the default value
func confValue(conf configDef) string {
	if value, exist := os.LookupEnv(conf.env); exist {
		return value
	}
	return conf.def
}

func asDuration(val string) time.Duration {
	duration, _ := strconv.Atoi(val)
	return time.Duration(duration) * time.Second
}
