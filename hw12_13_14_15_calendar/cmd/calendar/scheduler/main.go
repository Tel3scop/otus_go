package main

import (
	"context"
	"flag"
	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/app"
	"log"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	ctx := context.Background()
	flag.Parse()

	a, err := app.NewApp(ctx, configFile)
	if err != nil {
		log.Fatalf("failed to init app: %s", err.Error())
	}

	err = a.Run()
	if err != nil {
		log.Fatalf("failed to run app: %s", err.Error())
	}
}
