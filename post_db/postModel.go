package main

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

var post = &PostModel{}

type PostModel struct {
	DB *pgxpool.Pool
}

func (post *PostModel) Add(title, article, authorname string, authorid int32 ) (int, error) {

	stmt := "INSERT INTO post (title, article,published, authorname, authorid) VALUES ($1, $2, $3, $4) RETURNING id"

	var id uint64
	row := post.DB.QueryRow(context.Background(), stmt, title, article, time.Now(),  authorname, authorid)
	var err error
	err = row.Scan(&id)
	if err != nil {
		return 0, err
	}
	return int(id), nil
}


func (post *PostModel) Get(id int) (*Post, error) {

	stmt := "SELECT id, title, article, published,  authorname, authorid FROM post WHERE  id = $1"

	row := post.DB.QueryRow(context.Background(), stmt,  id)
	s := &Post{}

	err := row.Scan(&s.ID, &s.Title, &s.Article, &s.Published, &s.AuthorName, &s.AuthorId)
	if err != nil {
		if err.Error() == "no rows in result set" {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}
	return s, nil
}

func (post *PostModel) Delete(id int32) error {

	stmt := "DELETE FROM post WHERE id = $1"

	comtag, err := post.DB.Exec(context.Background(), stmt, id)
	if err != nil {
		log.Printf("Can not delete: %v\n", err)
		return err
	}

	if comtag.RowsAffected() == 0 {
		return errors.New("no changes")
	}

	return nil
}

func (post *PostModel) Update(title,article string, id int) error {
	stmt := "UPDATE articles SET title = $1, content = $2 WHERE id = $3"
	ct, err := post.DB.Exec(context.Background(), stmt, title, article, id)
	if err != nil {
		log.Printf("Can not Update: %v\n", err)
		return err
	}
	if ct.RowsAffected() == 0 {
		return errors.New("no changes")
	}

	return nil
}

func (post *PostModel) GetAll() ([]*Post, error) {

	stmt := "SELECT id, title, article, published, authorname, authorid FROM post"

	rows, err := post.DB.Query(context.Background(), stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post

	for rows.Next() {
		s := &Post{}
		err = rows.Scan(&s.ID, &s.Title, &s.Article, &s.Published,  &s.AuthorName, &s.AuthorId)
		if err != nil {
			return nil, err
		}
		posts = append(posts, s)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}
