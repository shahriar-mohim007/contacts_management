package state

import (
	"contacts/repository"
)

type State struct {
	Cfg        *Config
	Repository repository.Repository
}

func NewState(cfg *Config, db repository.Repository) *State {
	return &State{
		Cfg:        cfg,
		Repository: db,
	}
}
