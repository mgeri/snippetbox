package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mgeri/snippetbox/conf"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config and global logger
var configFile string
var pidFile string
var logger zerolog.Logger

// The Root CCobraorba Handler
var rootCmd = &cobra.Command{
	Version: conf.Version,
	Use:     conf.Executable,
	// Run: func(cmd *cobra.Command, args []string) {
	// 	server.ListenAndServe(&logger)
	// },
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Create Pid File
		pidFile = viper.GetString("pidfile")
		if pidFile != "" {
			file, err := os.OpenFile(pidFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0666)
			if err != nil {
				return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
			}
			defer file.Close()
			_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
			if err != nil {
				return fmt.Errorf("Could not create pid file: %s Error:%v", pidFile, err)
			}
		}
		return nil
	},
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		// Remove Pid file
		if pidFile != "" {
			os.Remove(pidFile)
		}
	},
}

// Execute starts the program
func Execute() {
	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

// This is the main initializer handling cli, config and log
func init() {
	// Initialize configuration
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Sets up the config file, environment etc
	viper.SetEnvPrefix(strings.ToUpper(conf.Executable))
	viper.SetTypeByDefaultValue(true)                      // If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.AutomaticEnv()                                   // Automatically use environment variables where available
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_")) // Environement variables use underscores instead of periods

	// If a config file is found, read it in.
	if configFile != "" {
		viper.SetConfigFile(configFile)
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read config file: %s ERROR: %s\n", configFile, err.Error())
			os.Exit(1)
		}
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath("./")
		viper.AddConfigPath("$HOME/." + conf.Executable)
		viper.ReadInConfig()
	}

}

func initLog() {

	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	logLevel := zerolog.DebugLevel
	// log level
	logLevel, err := zerolog.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse log level: %s ERROR: %s\n", viper.GetString("logger.level"), err.Error())
	}

	zerolog.SetGlobalLevel(logLevel)

	var logWriter io.Writer

	if viper.GetString("logger.file") != "" {
		logWriter = &lumberjack.Logger{
			Filename:   viper.GetString("logger.file"),
			MaxSize:    100, // megabytes
			MaxBackups: 3,
			MaxAge:     28, // days
		}
	} else {
		// log on stdout

		// pretty console logger
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		// log.Logger = log.Output(output)

		// output.FormatLevel = func(i interface{}) string {
		// 	return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		// }
		// output.FormatMessage = func(i interface{}) string {
		// 	return fmt.Sprintf("***%s***", i)
		// }
		// output.FormatFieldName = func(i interface{}) string {
		// 	return fmt.Sprintf("%s:", i)
		// }
		// output.FormatFieldValue = func(i interface{}) string {
		// 	return strings.ToUpper(fmt.Sprintf("%s", i))
		// }

		logWriter = output
	}

	logger = zerolog.New(logWriter).With().Timestamp().Logger()

}
