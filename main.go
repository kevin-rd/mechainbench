package main

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"mechainbench/app/config"
	"mechainbench/app/engine"
	"os"
)

func init() {
	log.Logger = log.Logger.With().Caller().Logger().Level(zerolog.InfoLevel)
	log.Logger = log.Logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	// parse config
	configFile := "config/config.toml"
	appConfig := config.ParseConfig(configFile)

	e := engine.NewDefaultEngine(appConfig)
	e.Run(context.Background())
	e.Close()
}
