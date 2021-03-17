package models

import (
	"errors"
	"time"
)

var (
	ErrEmpty           = errors.New("models: no record found")
	ErrInvalidData = errors.New("models: incorrect username or password")
	ErrDuplicateData     = errors.New("models: duplicate username")
)

type Post struct {
	ID         int
	AuthorName string
	Title      string
	Article    string
	Published  time.Time
	AuthorID   int
}
type User struct {
	ID             int
	Name           string
	Username          string
	Password       []byte
}