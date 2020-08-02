package main

import (
	"os"

	"github.com/jegutierrez/atlas-dns/internal/app"
	log "github.com/sirupsen/logrus"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("[event: container_fail_init][service: main] error running server, err: %s", err.Error())
	}
}

func run() error {
	server, err := app.NewApp()
	defer server.Shutdown()

	if err != nil {
		return err
	}

	r := server.RouterSetup()

	port := os.Getenv("DNS_PORT")
	if port == "" {
		port = "8080"
	}

	return r.Run(":" + port)
}
