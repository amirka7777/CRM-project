package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"
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
