package repository

import (
	"context"
	uuid "github.com/jackc/pgtype/ext/gofrs-uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
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

type UserDTO struct {
	ID    uuid.UUID `db:"id"`
	Name  string    `db:"name"`
	Email string    `db:"email"`
}

func NewPgRepository(databaseUrl string) (*PgRepository, error) {
	var err error
	once.Do(func() {
		db, dbErr := sqlx.Connect("pgx", databaseUrl)
		if dbErr != nil {
			err = dbErr
			return
		}
		if pingErr := db.Ping(); pingErr != nil {
			err = pingErr
			return
		}

		repository = &PgRepository{db: db}
	})
	return repository, err
}

func (repo *PgRepository) GetDB() *sqlx.DB {
	return repo.db
}

func (repo *PgRepository) GetUserByEmail(ctx context.Context, email string) (*UserDTO, error) {
	var user UserDTO
	query := `SELECT id, name, email FROM users WHERE email = $1`
	err := repo.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PgRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email, password, is_active, created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
	err := repo.db.QueryRowContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.IsActive).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}
