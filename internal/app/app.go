package app

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jegutierrez/atlas-dns/internal/databank"
	"github.com/jegutierrez/atlas-dns/internal/ping"
	log "github.com/sirupsen/logrus"
)

// App holds the app general config and request handlers.
type App struct {
	config *Config

	pingController     ping.Controller
	databankController databank.Controller
}

// GetSectorID returns the sector ID where the DNS is deployed.
func (a *App) GetSectorID() int {
	return a.config.SectorID
}

// ErrInvalidLogLevel error thrown when an invalid log level is configured.
var ErrInvalidLogLevel = errors.New("invalid level provided")

// NewApp instantiate the appropiate config and load the corresponding app container.
func NewApp() (*App, error) {
	env := os.Getenv("DNS_ENVIRONMENT")
	log.Infof("app running on environment: %s", env)

	var config *Config
	var err error
	switch env {
	case "PRODUCTION":
		config, err = NewProductionConfig()
	default:
		config, err = NewLocalConfig()
	}
	if err != nil {
		return nil, err
	}

	logLevel, err := log.ParseLevel(config.LogLevel)
	if err != nil {
		return nil, ErrInvalidLogLevel
	}
	log.SetLevel(logLevel)

	locationService := databank.NewLocationService()
	return &App{
		config: config,

		pingController:     ping.NewPongController(),
		databankController: databank.NewController(config.SectorID, locationService),
	}, nil
}

// RouterSetup register the http handlers to its corresponding paths.
func (a *App) RouterSetup() *gin.Engine {
	r := gin.Default()
	r.GET("/ping", a.pingController.PingHandler)
	r.POST("/calculate-databank-location", a.databankController.LocatorHandler)
	return r
}

// Shutdown holds termination tasks.
// For example: close DB connections.
func (a *App) Shutdown() {
	log.Info("[event: app_shutdown] application has been shutdown successfuly")
}
