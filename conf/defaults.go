package conf

import (
	"github.com/spf13/viper"
)

func init() {

	// Logger Defaults
	viper.SetDefault("logger.level", "debug")
	// if no file is specified, log on standard output
	viper.SetDefault("logger.file", "")

	// Pidfile
	viper.SetDefault("pidfile", "")

	// Server Configuration
	viper.SetDefault("server.address", ":4000")
	viper.SetDefault("server.staticDir", "./ui/static/")

	// Database Settings
	viper.SetDefault("storage.driver", "mysql")
	viper.SetDefault("storage.dsn", "snippetbox:snippetbox@tcp(localhost:3306)/snippetbox?tls=skip-verify&timeout=90s&multiStatements=true")
	viper.SetDefault("storage.wipe", "false")
	viper.SetDefault("storage.maxOpenConns", 3)
	viper.SetDefault("storage.maxIdleConns", 3)
	viper.SetDefault("storage.connMaxLifetime", 0)
	viper.SetDefault("storage.migrationDir", "./store/mysql/migrations")

}
