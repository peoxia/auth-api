package main

import (
	"context"
	"log"

	"github.com/peoxia/auth-api/config"
	"github.com/peoxia/auth-api/server"
)

func main() {

	log.Printf("Starting ...")

	ctx := context.Background()
	config, err := config.LoadConfig()

	if err != nil {
		log.Fatal(err.Error())
	}

	var s server.Server

	if err := s.Create(ctx, config); err != nil {
		log.Fatal(err.Error())
	}

	if err := s.Serve(ctx); err != nil {
		log.Fatal(err.Error())
	}
}
