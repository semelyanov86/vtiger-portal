package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"errors"
	"github.com/go-sql-driver/mysql"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"strings"
	"time"
)

var ErrDuplicateEmail = errors.New("duplicate email")

type UsersRepo struct {
	db *sql.DB
}

func NewUsersRepo(db *sql.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (r *UsersRepo) Insert(ctx context.Context, user *domain.User) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.Version = 1

	var query = `
				INSERT INTO users (crmid, first_name, last_name, description, account_id, account_name, title, department, email, password, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, version) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var args = []any{user.Crmid, user.FirstName, user.LastName, user.Description, user.AccountId, user.AccountName, user.Title, user.Department, user.Email, user.Password.Hash, user.CreatedAt, user.UpdatedAt, user.IsActive, user.MailingCity, user.MailingStreet, user.MailingCountry, user.OtherCountry, user.MailingState, user.MailingPoBox, user.OtherCity, user.OtherState, user.MailingZip, user.OtherZip, user.OtherStreet, user.OtherPoBox, user.Image, user.Version}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.email") {
				return ErrDuplicateEmail
			}
		}
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.Id = id

	return nil
}

func (r *UsersRepo) GetByEmail(ctx context.Context, email string) (domain.User, error) {
	var query = `SELECT id, crmid, first_name, last_name, description, account_id, account_name, title, department, email, password, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, version FROM users WHERE email = ?`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.Id,
		&user.Crmid,
		&user.FirstName,
		&user.LastName,
		&user.Description,
		&user.AccountId,
		&user.AccountName,
		&user.Title,
		&user.Department,
		&user.Email, &user.Password.Hash,
		&user.CreatedAt, &user.UpdatedAt,
		&user.IsActive, &user.MailingCity, &user.MailingStreet, &user.MailingCountry, &user.OtherCountry, &user.MailingState, &user.MailingPoBox, &user.OtherCity, &user.OtherState, &user.MailingZip, &user.OtherZip, &user.OtherStreet, &user.OtherPoBox, &user.Image, &user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return user, ErrRecordNotFound
		default:
			return user, err
		}
	}
	return user, nil
}

func (r *UsersRepo) Update(ctx context.Context, user *domain.User) error {
	var query = `UPDATE users SET first_name = ?, last_name = ?, description = ?, account_id = ?, account_name = ?, title = ?, department = ?, email = ?, password = ?, updated_at = NOW(), is_active = ?, mailingcity = ?, mailingstreet = ?, mailingcountry = ?, othercountry = ?, mailingstate = ?, mailingpobox = ?, othercity = ?, otherstate = ?, mailingzip = ?, otherzip = ?, otherstreet = ?, otherpobox = ?, image = ?, version = version + 1 WHERE id = ? AND version = ?`
	var args = []any{user.FirstName, user.LastName, user.Description, user.AccountId, user.AccountName, user.Title, user.Department, user.Email, user.Password.Hash, user.IsActive, user.MailingCity, user.MailingStreet, user.MailingCountry, user.OtherCountry, user.MailingState, user.MailingPoBox, user.OtherCity, user.OtherState, user.MailingZip, user.OtherZip, user.OtherStreet, user.OtherPoBox, user.Image, user.Id, user.Version}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		var mySQLError *mysql.MySQLError
		if errors.As(err, &mySQLError) {
			if mySQLError.Number == 1062 && strings.Contains(mySQLError.Message, "users.email") {
				return ErrDuplicateEmail
			}
		}
		if errors.Is(err, sql.ErrNoRows) {
			return ErrEditConflict
		}
		return err
	}
	user.Version++
	user.UpdatedAt = time.Now()
	return nil
}

func (r *UsersRepo) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*domain.User, error) {
	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT users.id, crmid, first_name, last_name, description, account_id, account_name, title, department, email, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = ?
        AND tokens.scope = ? 
        AND tokens.expired_at > ?`

	args := []any{hex.EncodeToString(tokenHash[:]), tokenScope, time.Now()}

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id, &user.Crmid, &user.FirstName, &user.LastName, &user.Description, &user.AccountId, &user.AccountName, &user.Title, &user.Department, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.MailingCity, &user.MailingStreet, &user.MailingCountry, &user.OtherCountry, &user.MailingState, &user.MailingPoBox, &user.OtherCity, &user.OtherState, &user.MailingZip, &user.OtherZip, &user.OtherStreet, &user.OtherPoBox, &user.Image, &user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}
