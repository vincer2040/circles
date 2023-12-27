package db

import (
	"context"
	"database/sql"

	_ "github.com/tursodatabase/libsql-client-go/libsql"
	"github.com/vincer2040/circles/internal/user"
)

type CirclesDB struct {
	db  *sql.DB
	ctx context.Context
}

func New(url string) (*CirclesDB, error) {
	db, err := sql.Open("libsql", url)
	if err != nil {
		return nil, err
	}
	ctx := context.Background()
	return &CirclesDB{db, ctx}, nil
}

func (cdb *CirclesDB) CreateUserTable() error {
	_, err := cdb.exec(
		`CREATE TABLE IF NOT EXISTS
        users(
            first TEXT NOT NULL,
            last TEXT NOT NULL,
            email TEXT NOT NULL,
            password TEXT NOT NULL,
            PRIMARY KEY(email)
        )`,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) InsertUser(user *user.User) error {
	_, err := cdb.exec(
		`INSERT INTO
        users(first, last, email, password)
        VALUES (?, ?, ?, ?)`,
		user.First, user.Last, user.Email, user.Password,
	)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) GetUser(email string) (*user.User, error) {
	var user user.User
	query := "SELECT first, last, email, password FROM users WHERE email = ?"
	row := cdb.db.QueryRow(query, email)
	err := row.Scan(&user.First, &user.Last, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (cdb *CirclesDB) DeleteUser(email string) error {
	_, err := cdb.exec("DELETE FROM users WHERE email = ?", email)
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) DropUserTable() error {
	_, err := cdb.exec("DROP TABLE IF EXISTS users")
	if err != nil {
		return err
	}
	return nil
}

func (cdb *CirclesDB) Close() {
	cdb.db.Close()
}

func (cdb *CirclesDB) exec(stmt string, args ...any) (*sql.Result, error) {
	res, err := cdb.db.ExecContext(cdb.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (cdb *CirclesDB) query(stmt string, args ...any) (*sql.Rows, error) {
	res, err := cdb.db.QueryContext(cdb.ctx, stmt, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}
