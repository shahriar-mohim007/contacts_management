package repository

import (
	"context"
	"database/sql"
	"github.com/gofrs/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
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

func (repo *PgRepository) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	query := `SELECT * FROM users WHERE email = $1`
	err := repo.db.GetContext(ctx, &user, query, email)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (repo *PgRepository) CreateUser(ctx context.Context, user *User) error {
	query := `INSERT INTO users (id, name, email, password, is_active,created_at, updated_at) 
	          VALUES ($1, $2, $3, $4, $5, NOW(), NOW()) RETURNING id`
	err := repo.db.QueryRowContext(ctx, query, user.ID, user.Name, user.Email, user.Password, user.IsActive).Scan(&user.ID)
	if err != nil {
		return err
	}
	return nil
}

func (repo *PgRepository) ActivateUserByID(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE users SET is_active = TRUE WHERE id = $1`
	_, err := repo.db.ExecContext(ctx, query, userID)
	return err
}

func (repo *PgRepository) GetAllContacts(ctx context.Context, userID uuid.UUID) ([]Contact, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, phone, street, city, state, zip_code, country FROM contacts WHERE user_id = $1", userID)
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Error().Err(err).Msg("Error closing rows")
		}
	}(rows)

	var contacts []Contact
	for rows.Next() {
		var contact Contact
		if err := rows.Scan(&contact.ID, &contact.Phone, &contact.Street, &contact.City, &contact.State, &contact.ZipCode, &contact.Country); err != nil {
			return nil, err
		}
		contacts = append(contacts, contact)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}

	return contacts, nil
}

func (repo *PgRepository) CreateContact(ctx context.Context, contact *Contact) error {
	query := `
        INSERT INTO contacts 
        (id, user_id, phone, street, city, state, zip_code, country, created_at, updated_at) 
        VALUES 
        ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
    `

	_, err := repo.db.ExecContext(
		ctx, query,
		contact.ID, contact.UserID, contact.Phone, contact.Street, contact.City, contact.State, contact.ZipCode, contact.Country,
	)
	return err
}
