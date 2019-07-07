package conf

import (
	"github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	viper.SetDefault("logger.level", "debug")
	viper.SetDefault("logger.file", "")

	// Pidfile
	viper.SetDefault("pidfile", "")

	// Server Configuration
	viper.SetDefault("server.address", ":4000")
	viper.SetDefault("server.staticDir", "./ui/static/")

}
