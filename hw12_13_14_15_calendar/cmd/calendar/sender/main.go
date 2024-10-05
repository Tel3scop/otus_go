package main

import (
	"context"
	"flag"
	"log"

	"github.com/Tel3scop/otus_go/hw12_13_14_15_calendar/internal/app"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "/etc/calendar/config.yaml", "Path to configuration file")
}

func main() {
	ctx := context.Background()
	flag.Parse()

	sender, err := app.NewSender(ctx, configFile)
	if err != nil {
		log.Fatalf("failed to init sender: %s", err.Error())
	}

	err = sender.Run()
	if err != nil {
		log.Fatalf("failed to run sender: %s", err.Error())
	}
}
