package repository

import (
	"database/sql"
	"fmt"
	"log"
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

func (r *ReposirotySupport) CreateContact(user_id int, name, email, phone string, created_at time.Time) error {

	query := `INSERT INTO contacts (user_id, name, email, phone, created_at) VALUES (?, ?, ?, ?, ?)`

	_, err := r.db.Exec(query, user_id, name, email, phone, created_at)
	if err != nil {
		return err
	}

	return nil

}

func (r *ReposirotySupport) GetContactsByUserID(userID int) ([]models.Contact, error) {

	query := `SELECT * FROM contacts WHERE user_id = ?`
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}

	contacts := []models.Contact{}

	defer rows.Close()
	for rows.Next() {
		var c models.Contact
		err := rows.Scan(&c.Id, &c.User_id, &c.Name, &c.Email, &c.Phone, &c.Created_at)
		if err != nil {
			log.Println("ошибка при вызове скан при SELECT из бд: ", err)
			continue
		}

		contacts = append(contacts, c)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return contacts, nil

}

func (r *ReposirotySupport) GetContactByIDAndUserId(contactID, userID int) (*models.Contact, error) {

	query := `SELECT id, user_id, name, email, phone, created_at FROM contacts WHERE id = ? AND user_id = ?`
	contact := &models.Contact{}
	row := r.db.QueryRow(query, contactID, userID)
	err := row.Scan(&contact.Id, &contact.User_id, &contact.Name, &contact.Email, &contact.Phone, &contact.Created_at)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("нет контакта в бд")
		}
		return nil, err
	}

	return contact, nil

}

func (r *ReposirotySupport) DeleteContactByID(contactID, userID int) (int64, error) {

	query := `DELETE FROM contacts WHERE id = ? AND user_id = ?`

	result, err := r.db.Exec(query, contactID, userID)
	if err != nil {
		return 0, err
	}

	rowsDelete, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsDelete, nil

}

func (r *ReposirotySupport) UpdateContact(userID, contactID int, name, email, phone string) (int64, error) {

	query := `UPDATE contacts SET name = ?, email = ?, phone = ? WHERE id = ? AND user_id = ?`
	result, err := r.db.Exec(query, name, email, phone, contactID, userID)
	if err != nil {
		return 0, err
	}

	rowsAffec, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return rowsAffec, nil

}
