package state

import (
	"contacts/repository"
	"github.com/rs/zerolog/log"
)

type State struct {
	cfg        *Config
	Repository *repository.PgRepository
}

func NewState(cfg *Config) *State {
	db, err := repository.NewPgRepository(cfg.DatabaseUrl)
	if err != nil {
		log.Fatal().Err(err).Msg("pg repository error")
	}
	return &State{
		cfg:        cfg,
		Repository: db,
	}
}
