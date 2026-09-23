package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/yukay/CRM/internal/models"
)

type ReposirotySupport struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *ReposirotySupport {
	return &ReposirotySupport{db: db}
}

func (r *ReposirotySupport) CreateUser(e, p string, c time.Time) error {

	query := `INSERT INTO users (email, password_hash, created_at) VALUES (?, ?, ?)`
	_, err := r.db.Exec(query, e, p, c)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return fmt.Errorf("Email уже занят")
		}
		return err
	}

	return nil

}

func (r *ReposirotySupport) GetUserByEmail(email string) (*models.UserLoginCheck, error) {

	query := `SELECT id, email, password_hash FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	user := &models.UserLoginCheck{}
	err := row.Scan(&user.Id, &user.Email, &user.PasswordHash)
	if err != nil {
		return nil, err
	}

	return user, nil

}
