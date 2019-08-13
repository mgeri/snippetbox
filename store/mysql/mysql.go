package mysql

import (
	"database/sql"
	"fmt"

	"github.com/rs/zerolog"

	// import mysql driver
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"

	// migrate using files
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/spf13/viper"
)

// New returns a new db pool from config
func New(logger *zerolog.Logger) (*sql.DB, error) {

	// Initialize a new connection pool
	db, err := sql.Open(viper.GetString("storage.driver"), viper.GetString("storage.dsn"))
	if err != nil {
		logger.Error().Err(err).Msg("Could not connect to database")
		return nil, fmt.Errorf("could not connect to database")
	}

	// Ping the database
	if err = db.Ping(); err != nil {
		logger.Error().Err(err).Msg("Could not ping database")
		return nil, fmt.Errorf("could not ping database")
	}

	// Set the maximum number of concurrently open connections to 5. Setting this
	// to less than or equal to 0 will mean there is no maximum limit (which
	// is also the default setting).
	db.SetMaxOpenConns(viper.GetInt("storage.maxOpenConns"))
	// configure pool idle and tmax life
	db.SetMaxIdleConns(viper.GetInt("storage.maxIdleConns"))
	db.SetConnMaxLifetime(viper.GetDuration("storage.connMaxLifetime"))

	logger.Info().Msg("Connected to database server")

	// Migration not set
	if !viper.IsSet("storage.migrationDir") || viper.GetString("storage.migrationDir") == "" {
		return db, nil
	}

	// Run migrations
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		logger.Error().Err(err).Msg("Could not start sql migration")
		return nil, fmt.Errorf("migration failed")
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", viper.GetString("storage.migrationDir")),
		"mysql", driver)

	if err != nil {
		logger.Error().Err(err).Msg("Migration init failed")
		return nil, fmt.Errorf("migration failed")
	}

	// Do we wipe the database
	if viper.GetBool("storage.wipe") {
		err = m.Down()
		switch {
		case err == migrate.ErrNoChange:
			logger.Info().Msgf("Database Down schema current")
		case err != nil:
			logger.Error().Err(err).Msg("Migrate Database Down Error")
			return nil, fmt.Errorf("migration failed")
		default:
			logger.Warn().Msgf("Database wipe completed")
		}
	}

	// Perform the migration up
	err = m.Up()
	switch {
	case err == migrate.ErrNoChange:
		logger.Info().Msgf("Database Up schema current")
	case err != nil:
		logger.Error().Err(err).Msg("Migrate Database Up Error")
		return nil, fmt.Errorf("migration failed")
	default:
		logger.Warn().Msgf("Database migration completed")
	}

	return db, nil

}
