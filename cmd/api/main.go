package main

import (
	"contacts/cmd/httpserver"
	"contacts/repository"
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
	db, err := repository.NewPgRepository(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("pg repository error")
	}
	appState := state.NewState(cfg, db)
	httpserver.Serve(appState)

}
