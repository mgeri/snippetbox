package main

import (
	"flag"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config is an exported type that
// contains configuration info
type Config struct {
	Addr      string
	StaticDir string
	LogFile   string
}

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	log *zerolog.Logger
	cfg *Config
}

func main() {

	cfg := new(Config)

	// read params
	flag.StringVar(&cfg.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&cfg.StaticDir, "static-dir", "../../ui/static", "Path to static assets")
	flag.StringVar(&cfg.LogFile, "log-file", "", "Log file path (default is Stdout)")
	flag.Parse()

	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs

	var logWriter io.Writer

	if cfg.LogFile != "" {
		logWriter = &lumberjack.Logger{
			Filename:   cfg.LogFile,
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

	log := zerolog.New(logWriter).With().Timestamp().Logger()

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		log: &log,
		cfg: cfg,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:    cfg.Addr,
		Handler: app.routes(),
	}

	app.log.Info().Msgf("Starting server on %s", cfg.Addr)
	err := srv.ListenAndServe()
	if err != nil {
		app.log.Fatal().Err(err).Msg("Startup failed")
	}
}
