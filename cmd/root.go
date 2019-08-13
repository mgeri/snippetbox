package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/mgeri/snippetbox/conf"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config and global logger
var configFile string
var pidFile string
var logger zerolog.Logger

// The Root Cobra Handler
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
				return fmt.Errorf("could not create pid file: %s Error:%v", pidFile, err)
			}
			defer file.Close()
			_, err = fmt.Fprintf(file, "%d\n", os.Getpid())
			if err != nil {
				return fmt.Errorf("could not create pid file: %s Error:%v", pidFile, err)
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

// This is the main initializer handling cli, config and log
func init() { // nolint: gochecknoinits
	// Initialize configuration
	cobra.OnInitialize(initConfig, initLog)
	rootCmd.PersistentFlags().StringVarP(&configFile, "config", "c", "", "Config file")
}

// Execute starts the program
func Execute() {
	// Run the program
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	// Sets up the config file, environment etc
	viper.SetEnvPrefix(strings.ToUpper(conf.Executable))
	// If a default value is []string{"a"} an environment variable of "a b" will end up []string{"a","b"}
	viper.SetTypeByDefaultValue(true)
	// Automatically use environment variables where available
	viper.AutomaticEnv()
	// Environment variables use underscores instead of periods
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

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
		err := viper.ReadInConfig()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Could not read config file: .%s ERROR: %s\n", conf.Executable, err.Error())
			os.Exit(1)
		}
	}

}

func initLog() {

	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	// log level
	var logLevel zerolog.Level
	var err error
	logLevel, err = zerolog.ParseLevel(viper.GetString("logger.level"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse log level: %s ERROR: %s\n", viper.GetString("logger.level"), err.Error())
		logLevel = zerolog.DebugLevel
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

	// set global logger
	log.Logger = logger

}
