package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofrs/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"strings"
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

func (repo *PgRepository) GetContactByID(ctx context.Context, contactID uuid.UUID) (*ContactWithUserResponse, error) {
	query := `
        SELECT 
            contacts.id AS contact_id,
            contacts.phone,
            contacts.street,
            contacts.city,
            contacts.state,
            contacts.zip_code,
            contacts.country,
            users.name AS user_name,
            users.email AS user_email
        FROM 
            contacts
        JOIN 
            users ON contacts.user_id = users.id
        WHERE 
            contacts.id = $1;
    `

	var response ContactWithUserResponse
	err := repo.db.QueryRowContext(ctx, query, contactID).Scan(
		&response.ContactID,
		&response.Phone,
		&response.Street,
		&response.City,
		&response.State,
		&response.ZipCode,
		&response.Country,
		&response.UserName,
		&response.UserEmail,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("no contact found with ID: %s", contactID)
		}
		return nil, err
	}
	return &response, nil
}

func (repo *PgRepository) PatchContact(ctx context.Context, contactID uuid.UUID, contact *Contact) error {

	var queryParts []string
	var args []interface{}
	argID := 1

	if contact.Phone != "" {
		queryParts = append(queryParts, fmt.Sprintf("phone = $%d", argID))
		args = append(args, contact.Phone)
		argID++
	}
	if contact.Street != "" {
		queryParts = append(queryParts, fmt.Sprintf("street = $%d", argID))
		args = append(args, contact.Street)
		argID++
	}
	if contact.City != "" {
		queryParts = append(queryParts, fmt.Sprintf("city = $%d", argID))
		args = append(args, contact.City)
		argID++
	}
	if contact.State != "" {
		queryParts = append(queryParts, fmt.Sprintf("state = $%d", argID))
		args = append(args, contact.State)
		argID++
	}
	if contact.ZipCode != "" {
		queryParts = append(queryParts, fmt.Sprintf("zip_code = $%d", argID))
		args = append(args, contact.ZipCode)
		argID++
	}
	if contact.Country != "" {
		queryParts = append(queryParts, fmt.Sprintf("country = $%d", argID))
		args = append(args, contact.Country)
		argID++
	}

	if len(queryParts) == 0 {
		return fmt.Errorf("no fields provided to update")
	}

	query := fmt.Sprintf("UPDATE contacts SET %s WHERE id = $%d", strings.Join(queryParts, ", "), argID)
	args = append(args, contactID)

	_, err := repo.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (repo *PgRepository) DeleteContactByID(ctx context.Context, contactID uuid.UUID) error {
	query := `
        DELETE FROM contacts
        WHERE id = $1;
    `

	result, err := repo.db.ExecContext(ctx, query, contactID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
