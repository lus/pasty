package main

import (
	"context"
	"errors"
	"github.com/lus/pasty/internal/cleanup"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/meta"
	"github.com/lus/pasty/internal/reports"
	"github.com/lus/pasty/internal/storage"
	"github.com/lus/pasty/internal/storage/postgres"
	"github.com/lus/pasty/internal/storage/sqlite"
	"github.com/lus/pasty/internal/web"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"os/signal"
	"strings"
)

func main() {
	// Set up the logger
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	if !meta.IsProdEnvironment() {
		log.Logger = log.Output(zerolog.ConsoleWriter{
			Out: os.Stderr,
		})
		log.Warn().Msg("This distribution was compiled for development mode and is thus not meant to be run in production!")
	}

	// Load the configuration
	config.Compatibility()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Could not load the configuration.")
	}

	// Adjust the log level
	if !meta.IsProdEnvironment() {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	} else {
		level, err := zerolog.ParseLevel(cfg.LogLevel)
		if err != nil {
			log.Warn().Msg("An invalid log level was configured. Falling back to 'info'.")
			level = zerolog.InfoLevel
		}
		zerolog.SetGlobalLevel(level)
	}

	// Determine the correct storage driver to use
	var driver storage.Driver
	switch strings.TrimSpace(strings.ToLower(cfg.StorageDriver)) {
	case "postgres":
		driver = postgres.New(cfg.Postgres.DSN)
		break
	case "sqlite":
		driver = sqlite.New(cfg.SQLite.File)
		break
	case "file":
		// TODO: Readme notice
		log.Fatal().Msg("You have configured the legacy 'file' storage driver. This storage driver has been removed in favor of PostgreSQL and SQLite, but the latter one may be a seamless alternative for you. Head over to the projects README for more information.")
		break
	case "mongodb":
	case "s3":
		// TODO: Readme notice
		log.Fatal().Msg("You have configured a legacy storage driver. This storage driver has been removed in favor of PostgreSQL and SQLite, but the migration process is well-documented. Head over to the projects README for more information.")
		break
	default:
		log.Fatal().Str("driver_name", cfg.StorageDriver).Msg("An invalid storage driver name was given.")
		return
	}

	// Initialize the configured storage driver
	log.Info().Str("driver_name", cfg.StorageDriver).Msg("Initializing the storage driver...")
	if err := driver.Initialize(context.Background()); err != nil {
		log.Fatal().Err(err).Str("driver_name", cfg.StorageDriver).Msg("The storage driver could not be initialized.")
		return
	}
	defer func() {
		log.Info().Msg("Shutting down the storage driver...")
		if err := driver.Close(); err != nil {
			log.Err(err).Str("driver_name", cfg.StorageDriver).Msg("Could not shut down the storage driver.")
		}
	}()

	// Schedule the cleanup task if configured
	if cfg.Cleanup.Enabled {
		task := &cleanup.Task{
			Interval:    cfg.Cleanup.TaskInterval,
			MaxPasteAge: cfg.Cleanup.PasteLifetime,
			Repository:  driver.Pastes(),
		}
		log.Info().Msg("Scheduling the cleanup task...")
		task.Start()
		defer func() {
			log.Info().Msg("Shutting down the cleanup task...")
			task.Stop()
		}()
	}

	// Start the web server
	log.Info().Str("address", cfg.Address).Msg("Starting the web server...")
	webServer := &web.Server{
		Address:                   cfg.Address,
		Storage:                   driver,
		PasteIDLength:             cfg.PasteIDLength,
		PasteIDCharset:            cfg.PasteIDCharset,
		PasteLengthCap:            cfg.PasteLengthCap,
		ModificationTokensEnabled: cfg.ModificationTokensEnabled,
		ModificationTokenLength:   cfg.ModificationTokenLength,
		ModificationTokenCharset:  cfg.ModificationTokenCharset,
	}
	if cfg.Reports.Enabled {
		webServer.ReportClient = &reports.Client{
			WebhookURL:   cfg.Reports.WebhookURL,
			WebhookToken: cfg.Reports.WebhookToken,
		}
	}
	if cfg.ModificationTokenMaster != "" {
		webServer.AdminTokens = []string{cfg.ModificationTokenMaster}
	}
	go func() {
		if err := webServer.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("Could not start the web server.")
		}
	}()
	defer func() {
		log.Info().Msg("Shutting down the web server...")
		if err := webServer.Shutdown(context.Background()); err != nil {
			log.Err(err).Msg("Could not shut down the web server.")
		}
	}()

	// Wait for an interrupt signal
	log.Info().Msg("The application has been started. Use Ctrl+C to shut it down.")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan
}
