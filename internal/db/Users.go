package db

import (
	"database/sql"
)

type UserModel struct {
	DB *sql.DB
}

func (usm *UserModel) Insert(user *User) error {
	query := `INSERT INTO public.users (email, password) VALUES ($1, $2) RETURNING id, created_at`
	args := []any{user.Email, user.Password}

	return usm.DB.QueryRow(query, args...).Scan(&user.ID, &user.CreatedAt)
}
