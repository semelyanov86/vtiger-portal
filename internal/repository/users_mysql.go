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
				INSERT INTO users (crmid, first_name, last_name, description, account_id, account_name, title, department, email, password, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, imageattachmentids, version, phone, assigned_user_id) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	var args = []any{user.Crmid, user.FirstName, user.LastName, user.Description, user.AccountId, user.AccountName, user.Title, user.Department, user.Email, user.Password.Hash, user.CreatedAt, user.UpdatedAt, user.IsActive, user.MailingCity, user.MailingStreet, user.MailingCountry, user.OtherCountry, user.MailingState, user.MailingPoBox, user.OtherCity, user.OtherState, user.MailingZip, user.OtherZip, user.OtherStreet, user.OtherPoBox, user.Image, user.Imageattachmentids, user.Version, user.Phone, user.AssignedUserId}

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
	var query = `SELECT id, crmid, first_name, last_name, description, account_id, account_name, title, department, email, password, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, imageattachmentids, version, phone, assigned_user_id, otp_enabled, otp_verified FROM users WHERE email = ?`
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
		&user.IsActive, &user.MailingCity, &user.MailingStreet, &user.MailingCountry, &user.OtherCountry, &user.MailingState, &user.MailingPoBox, &user.OtherCity, &user.OtherState, &user.MailingZip, &user.OtherZip, &user.OtherStreet, &user.OtherPoBox, &user.Image, &user.Imageattachmentids, &user.Version, &user.Phone, &user.AssignedUserId, &user.Otp_enabled, &user.Otp_verified,
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

func (r *UsersRepo) GetById(ctx context.Context, id int64) (domain.User, error) {
	var query = `SELECT id, crmid, first_name, last_name, description, account_id, account_name, title, department, email, password, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, imageattachmentids, version, phone, assigned_user_id, otp_verified, otp_enabled, otp_secret, otp_auth_url FROM users WHERE id = ?`
	var user domain.User

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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
		&user.IsActive, &user.MailingCity, &user.MailingStreet, &user.MailingCountry, &user.OtherCountry, &user.MailingState, &user.MailingPoBox, &user.OtherCity, &user.OtherState, &user.MailingZip, &user.OtherZip, &user.OtherStreet, &user.OtherPoBox, &user.Image, &user.Imageattachmentids, &user.Version, &user.Phone, &user.AssignedUserId, &user.Otp_verified, &user.Otp_enabled, &user.Otp_secret, &user.Otp_auth_url,
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
	var query = `UPDATE users SET first_name = ?, last_name = ?, description = ?, account_id = ?, account_name = ?, title = ?, department = ?, email = ?, password = ?, updated_at = NOW(), is_active = ?, mailingcity = ?, mailingstreet = ?, mailingcountry = ?, othercountry = ?, mailingstate = ?, mailingpobox = ?, othercity = ?, otherstate = ?, mailingzip = ?, otherzip = ?, otherstreet = ?, otherpobox = ?, image = ?, imageattachmentids = ?, version = version + 1, phone = ?, assigned_user_id = ?, otp_verified = ?, otp_enabled = ?, otp_auth_url = ?, otp_secret = ? WHERE id = ?`
	var args = []any{user.FirstName, user.LastName, user.Description, user.AccountId, user.AccountName, user.Title, user.Department, user.Email, user.Password.Hash, user.IsActive, user.MailingCity, user.MailingStreet, user.MailingCountry, user.OtherCountry, user.MailingState, user.MailingPoBox, user.OtherCity, user.OtherState, user.MailingZip, user.OtherZip, user.OtherStreet, user.OtherPoBox, user.Image, user.Imageattachmentids, user.Phone, user.AssignedUserId, user.Otp_verified, user.Otp_enabled, user.Otp_auth_url, user.Otp_secret, user.Id}

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
        SELECT users.id, crmid, first_name, last_name, description, account_id, account_name, title, department, email, created_at, updated_at, is_active, mailingcity, mailingstreet, mailingcountry, othercountry, mailingstate, mailingpobox, othercity, otherstate, mailingzip, otherzip, otherstreet, otherpobox, image, imageattachmentids, version, phone, assigned_user_id, otp_enabled, otp_verified, otp_auth_url, otp_secret
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = ?
        AND tokens.scope = ? 
        AND tokens.expired_at > ?`

	args := []any{hex.EncodeToString(tokenHash[:]), tokenScope, time.Now()}

	var user domain.User

	err := r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.Id, &user.Crmid, &user.FirstName, &user.LastName, &user.Description, &user.AccountId, &user.AccountName, &user.Title, &user.Department, &user.Email, &user.CreatedAt, &user.UpdatedAt, &user.IsActive, &user.MailingCity, &user.MailingStreet, &user.MailingCountry, &user.OtherCountry, &user.MailingState, &user.MailingPoBox, &user.OtherCity, &user.OtherState, &user.MailingZip, &user.OtherZip, &user.OtherStreet, &user.OtherPoBox, &user.Image, &user.Imageattachmentids, &user.Version, &user.Phone, &user.AssignedUserId, &user.Otp_enabled, &user.Otp_verified, &user.Otp_auth_url, &user.Otp_secret,
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

func (r *UsersRepo) SaveOtp(ctx context.Context, otpSecret string, otpUrl string, userId int64) error {
	var query = `UPDATE users SET otp_secret = ?, otp_auth_url = ?, version = version + 1, updated_at = NOW() WHERE id = ?`
	var args = []any{otpSecret, otpUrl, userId}

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
	return nil
}

func (r *UsersRepo) EnableAndVerifyOtp(ctx context.Context, userId int64) error {
	var query = `UPDATE users SET otp_enabled = ?, otp_verified = ?, version = version + 1, updated_at = NOW() WHERE id = ?`
	var args = []any{1, 1, userId}

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
	return nil
}

func (r *UsersRepo) VerifyOrInvalidateOtp(ctx context.Context, userId int64, valid bool) error {
	value := 0
	if valid {
		value = 1
	}
	var query = `UPDATE users SET otp_verified = ?, version = version + 1, updated_at = NOW() WHERE id = ?`
	var args = []any{value, userId}

	_, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return nil
}
