package repository

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"encoding/hex"
	"github.com/semelyanov86/vtiger-portal/internal/domain"
	"time"
)

type TokensRepo struct {
	db *sql.DB
}

func NewTokensRepo(db *sql.DB) *TokensRepo {
	return &TokensRepo{
		db: db,
	}
}

func (t TokensRepo) New(ctx context.Context, userId int64, ttl time.Duration, scope string) (*domain.Token, error) {
	token, err := generateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	err = t.Insert(ctx, token)
	return token, err
}

func (t TokensRepo) Insert(ctx context.Context, token *domain.Token) error {
	var query = `INSERT INTO tokens (hash, user_id, expired_at, scope) VALUES (?,?,?,?)`
	var args = []any{token.Hash, token.UserId, token.Expiry, token.Scope}

	result, err := t.db.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	token.ID = id
	return err
}

func (t TokensRepo) DeleteAllForUser(ctx context.Context, scope string, userId int64) error {
	var query = `DELETE FROM tokens WHERE scope = ? AND user_id = ?`

	_, err := t.db.ExecContext(ctx, query, scope, userId)
	return err
}

func generateToken(userId int64, ttl time.Duration, scope string) (*domain.Token, error) {
	var token = &domain.Token{
		UserId: userId,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}
	var randomBytes = make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)
	var hash = sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hex.EncodeToString(hash[:])
	return token, nil
}
