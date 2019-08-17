package server

import (
	"html/template"
	"net/http"

	"github.com/jmoiron/sqlx"

	"github.com/mgeri/snippetbox/store"
	"github.com/mgeri/snippetbox/store/mysql"

	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

// Define an application struct to hold the application-wide dependencies for the
// web application. For now we'll only include fields for the two custom loggers, but
// we'll add more to it as the build progresses.
type application struct {
	logger        *zerolog.Logger
	db            *sqlx.DB
	snippetStore  store.SnippetStore
	templateCache map[string]*template.Template
}

// ListenAndServe run Snippetbox server
func ListenAndServe(logger *zerolog.Logger) {

	var err error
	var db *sqlx.DB
	var snippetStore store.SnippetStore

	switch viper.GetString("storage.driver") {
	case "mysql":
		db, err = mysql.New(logger)
		snippetStore = mysql.NewMysqlSnippetStore(logger, db)
	default:
		db, err = mysql.New(logger)
		snippetStore = mysql.NewMysqlSnippetStore(logger, db)
	}
	if err != nil {
		logger.Fatal().Msgf("Database Error %s", err)
	}

	defer db.Close()

	// Initialize a new template cache...
	templateCache, err := newTemplateCache("./ui/html/")
	if err != nil {
		logger.Fatal().Msgf("Template cache Error %s", err)
	}

	// Initialize a new instance of application containing the dependencies.
	app := &application{logger, db, snippetStore, templateCache}

	// Initialize a new http.Server struct. We set the Addr and Handler fields so
	// that the server uses the same network address and routes as before, and set
	// the ErrorLog field so that the server now uses the custom errorLog logger in
	// the event of any problems.
	srv := &http.Server{
		Addr:    viper.GetString("server.address"),
		Handler: app.routes(),
	}

	app.logger.Info().Msgf("Starting server on %s", viper.GetString("server.address"))
	err = srv.ListenAndServe()
	if err != nil {
		app.logger.Fatal().Err(err).Msg("Startup failed")
	}
}
