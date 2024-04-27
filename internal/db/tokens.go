package db

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/base32"
	"time"
)

const (
	ScopeActivation = "activation"
)

type Token struct {
	Plaintext string
	Hash      []byte
	UserID    int64
	Expiry    time.Time
	Scope     string
}

func generateToken(userID int64, ttl time.Duration, scope string) (*Token, error) {
	token := &Token{
		UserID: userID,
		Expiry: time.Now().Add(ttl),
		Scope:  scope,
	}

	randombytes := make([]byte, 16)
	_, err := rand.Read(randombytes)
	if err != nil {
		return nil, err
	}
	token.Plaintext = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randombytes)
	hash := sha256.Sum256([]byte(token.Plaintext))
	token.Hash = hash[:]

	return token, nil
}

type TokenModel struct {
	DB *sql.DB
}

func (tm TokenModel) New(userid int64, ttl time.Duration, scope string) (*Token, error) {
	token, err := generateToken(userid, ttl, scope)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (tm TokenModel) Insert(token *Token) error {
	query := `INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`
	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}
	err := tm.DB.QueryRow(query, args...)
	if err != nil {
		return err.Err()
	}
	return nil
}
