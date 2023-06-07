package main

import (
	"context"
	"github.com/lus/pasty/internal/config"
	"github.com/lus/pasty/internal/meta"
	"github.com/lus/pasty/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"os/signal"
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

	// Initialize the configured storage driver
	driver, ok := storage.ResolveDriver(cfg.StorageDriver)
	if !ok {
		log.Fatal().Str("driver_name", cfg.StorageDriver).Msg("An invalid storage driver name was given.")
		return
	}
	log.Info().Str("driver_name", cfg.StorageDriver).Msg("Initializing the storage driver...")
	if err := driver.Initialize(context.Background(), cfg); err != nil {
		log.Fatal().Err(err).Str("driver_name", cfg.StorageDriver).Msg("The storage driver could not be initialized.")
		return
	}
	defer func() {
		log.Info().Msg("Shutting down the storage driver...")
		if err := driver.Close(); err != nil {
			log.Err(err).Str("driver_name", cfg.StorageDriver).Msg("Could not shut down the storage driver.")
		}
	}()

	// Wait for an interrupt signal
	log.Info().Msg("The application has been started. Use Ctrl+C to shut it down.")
	shutdownChan := make(chan os.Signal, 1)
	signal.Notify(shutdownChan, os.Interrupt)
	<-shutdownChan
}
