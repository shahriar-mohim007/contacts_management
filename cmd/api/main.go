package main

import (
	"contacts/cmd/httpserver"
	"contacts/state"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	_ = godotenv.Load()
	cfg, err := state.NewConfig()
	if err != nil {
		log.Fatal().Err(err).Msg("Config parsing failed")
	}
	logLevel, err := zerolog.ParseLevel(cfg.LogLevel)

	if err != nil {
		logLevel = zerolog.DebugLevel
	}
	zerolog.SetGlobalLevel(logLevel)
	appState := state.NewState(cfg)
	httpserver.Serve(appState)

}
