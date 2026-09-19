package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func InitDB(name string) (*sql.DB, error) {

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec("PRAGMA foreign_keys = ON;")
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil

}

func CreateUsersTable(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS users(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	email TEXT UNIQUE NOT NULL,
	password_hash TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}

func CreateContactsTable(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS contacts(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	name TEXT NOT NULL,
	email TEXT NOT NULL,
	phone TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,

	FOREIGN KEY (user_id) REFERENCES users (id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}

func CreateDealsTable(db *sql.DB) error {

	query := `
	CREATE TABLE IF NOT EXISTS deals(
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	user_id INTEGER NOT NULL,
	contact_id INTEGER NOT NULL,
	title TEXT NOT NULL,
	amount INTEGER NOT NULL,
	status TEXT NOT NULL,
	created_at TIMESTAMP NOT NULL,

	FOREIGN KEY (user_id) REFERENCES users (id)
		ON DELETE CASCADE
		ON UPDATE CASCADE,

	FOREIGN KEY (contact_id) REFERENCES contacts (id)
		ON DELETE CASCADE
		ON UPDATE CASCADE
	)`

	_, err := db.Exec(query)
	if err != nil {
		return err
	}

	return nil

}
