package app

import (
	"errors"
	"io"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
)

var (
	// ErrMissingSectorID error returned when a sector ID is not provided.
	ErrMissingSectorID = errors.New("missing sector ID")
	// ErrInvalidSectorID error returned when an invalid sector ID is given.
	ErrInvalidSectorID = errors.New("invalid sector ID, must be a valid integer")
)

// Config holds the application configurable params.
type Config struct {
	LogLevel  string
	LogOutput io.Writer
	SectorID  int
}

// NewLocalConfig creates the default config to run the app locally.
func NewLocalConfig() (*Config, error) {
	return &Config{
		LogLevel:  log.DebugLevel.String(),
		LogOutput: os.Stdout,
		SectorID:  1,
	}, nil
}

// NewProductionConfig loads the config to run the app into docker container.
func NewProductionConfig() (*Config, error) {
	logLevel := os.Getenv("DSN_LOG_LEVEL")
	if logLevel == "" {
		logLevel = log.WarnLevel.String()
	}

	sectorID := os.Getenv("DSN_SECTOR_ID")
	if sectorID == "" {
		return nil, ErrMissingSectorID
	}
	sectorIDInt, err := strconv.Atoi(sectorID)
	if err != nil {
		return nil, ErrInvalidSectorID
	}

	return &Config{
		LogLevel:  logLevel,
		LogOutput: os.Stdout,
		SectorID:  sectorIDInt,
	}, nil
}
