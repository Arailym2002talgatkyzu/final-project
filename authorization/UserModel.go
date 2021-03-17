package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/pgxpool"
	"golang.org/x/crypto/bcrypt"
	"strings"
)

var user = &UserModel{}

type UserModel struct {
	DB *pgxpool.Pool
}

func (user *UserModel) Insert(name, username, password string) error {
	Password, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	stmt := "INSERT INTO users (name, username, password) VALUES($1, $2, $3)"

	_, err = user.DB.Exec(context.Background(), stmt, name, username, string(Password))
	if err != nil {
		if strings.Contains(err.Error(), "users_uc_email") {
			return ErrDuplicateData
		}
		return err
	}
	return nil
}

func (user *UserModel) Authenticate(username, password string) (int, error) {
	var id int
	var Password []byte
	stmt := "SELECT id, password FROM users WHERE email = $1 AND active = TRUE"
	row := user.DB.QueryRow(context.Background(), stmt, username)
	err := row.Scan(&id, &Password)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return 0, ErrEmpty
		} else {
			return 0, err
		}
	}
	err = bcrypt.CompareHashAndPassword(Password, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidData
		} else {
			return 0, err
		}
	}

	return id, nil
}

func (user *UserModel) Get(id int) (*User, error) {
	u := &User{}
	stmt := `SELECT id, name, username FROM users WHERE id = $1`
	err := user.DB.QueryRow(context.Background(), stmt, id).Scan(&u.ID, &u.Name, &u.Username)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrEmpty
		} else {
			return nil, err
		}
	}
	return u, nil
}
