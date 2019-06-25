package server

import (
	"net/http"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	logger *zerolog.Logger
}

func ListenAndServe(logger *zerolog.Logger) {

	// Initialize a new instance of application containing the dependencies.
	app := &application{
		logger: logger,
	}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: app.routes(),
	}

	app.logger.Info().Msgf("Starting server on %s", viper.GetString("server.address"))
	err := srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal().Err(err).Msg("Startup failed")
	}
}
