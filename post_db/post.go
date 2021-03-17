package main

import (
	"errors"
	"time"
)
var ErrNoRecord = errors.New("models: no matching record found")
type Post struct {
	ID       int
	AuthorId int
	AuthorName string
	Title    string
	Article  string
	Published time.Time
}
