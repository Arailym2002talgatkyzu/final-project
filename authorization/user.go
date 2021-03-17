package main

import "errors"

var (
	ErrEmpty           = errors.New("models: no record found")
	ErrInvalidData = errors.New("models: incorrect username or password")
	ErrDuplicateData     = errors.New("models: duplicate username")
)
type User struct {
	ID             int
	Name           string
	Username          string
	Password       []byte
}
