package repository

import (
	"github.com/jmoiron/sqlx"
	"sync"
)

type PgRepository struct {
	db *sqlx.DB
}

var (
	once       sync.Once
	repository *PgRepository
)

func NewPgRepository(databaseUrl string) (*PgRepository, error) {
	var err error
	once.Do(func() {
		db, dbErr := sqlx.Connect("pgx", databaseUrl)
		if dbErr != nil {
			err = dbErr
			return
		}
		repository = &PgRepository{db: db}
	})
	return repository, err
}

func (repo *PgRepository) GetDB() *sqlx.DB {
	return repo.db
}
