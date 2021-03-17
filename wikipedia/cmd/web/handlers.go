package main

import (
	"context"
	"errors"
	"github.com/Arailym2002talgatkyzu/final-project/authorization/authpb"
	"github.com/Arailym2002talgatkyzu/final-project/wikipedia/pkg/forms"
	"github.com/Arailym2002talgatkyzu/final-project/wikipedia/pkg/models"

	"fmt"
	"github.com/Arailym2002talgatkyzu/final-project/post_db/postpb"

	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

func StringToTime(value string) time.Time {
	result, _ := time.Parse("2006-01-02 15:04:05.999999999 -0700 MST", value)
	return result
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

	req := &postpb.GetPostsRequest{}
	stream, err := app.posts.GetPosts(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetPosts RPC: %v", err)
	}
	defer stream.CloseSend()

	var posts []*models.Post

LOOP:
	for {
		res, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break LOOP
			}
			log.Fatalf("Error while receiving from GetPosts RPC: %v", err)
		}
		log.Printf("Response from GetPOsts RPC, PostID: %v \n", res.GetPost().GetId())

		reqUser := &authpb.GetUserRequest{
			Id: res.GetPost().GetAuthorid(),
		}
		resUser, err := app.auth.GetUser(context.Background(), reqUser)
		if err != nil {
			log.Fatalf("Error while calling GetUser RPC: %v", err)
		}

		tempPost := &models.Post{
			ID:         int(res.GetPost().GetId()),
			AuthorID:   int(res.GetPost().GetAuthorid()),
			AuthorName: resUser.GetUser().GetName(),
			Title:      res.GetPost().GetTitle(),
			Article:    res.GetPost().GetArticle(),
			Published:    StringToTime(res.GetPost().GetPublished()),
		}
		posts = append(posts, tempPost)
	}

	// Web Design
	app.render(w, r, "home.page.tmpl", &templateData{
		Posts: posts,
	})
}

func (app *application) showArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	req := &postpb.GetPostRequest{
		Id: int32(id),
	}
	res, err := app.posts.GetPost(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GetPost RPC: %v", err)
	}
	log.Printf("Response from GetPost RPC: %s, ArticleID: %v", res.GetResult(), res.GetPost().GetId())

	reqUser := &authpb.GetUserRequest{
		Id: res.GetPost().GetAuthorid(),
	}

	resUser, err := app.auth.GetUser(context.Background(), reqUser)
	if err != nil {
		log.Fatalf("Error while calling GetUser RPC: %v", err)
	}

	post := &models.Post{
		ID:         int(res.GetPost().GetId()),
		AuthorID:   int(res.GetPost().GetAuthorid()),
		AuthorName: resUser.GetUser().GetName(),
		Title:      res.GetPost().GetTitle(),
		Article:    res.GetPost().GetArticle(),
		Published:    StringToTime(res.GetPost().GetPublished()),

	}

	// Web Design
	app.render(w, r, "show.page.tmpl", &templateData{
		Post: post,
		UserID:  app.session.GetInt(r, "authenticatedUserID"),
	})
}

func (app *application) createArticleForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "create.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) createArticle(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content", "expires")
	form.MaxLength("title", 100)
	form.PermittedValues("expires", "365", "7", "1")
	if !form.Valid() {
		app.render(w, r, "create.page.tmpl", &templateData{Form: form})
		return
	}

	req := &postpb.InsertPostRequest{
		Post: &postpb.Post{
			Title:    form.Get("title"),
			Article:  form.Get("content"),
			Authorname: app.session.GetString(r, "authenticatedUserName"),
			Authorid: int32(app.session.GetInt(r, "authenticatedUserID")),
		},
	}
	res, err := app.posts.InsertPost(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling InsertPost RPC: %v", err)
	}
	log.Printf("Response from InsertPost RPC: %v", res.GetResult())

	app.session.Put(r, "flash", "Article successfully created!")
	http.Redirect(w, r, fmt.Sprintf("/article/%d", res.GetId()), http.StatusSeeOther)
}

func (app *application) deleteArticle(w http.ResponseWriter, r *http.Request) {

	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	req := &postpb.DeletePostRequest{
		Id: int32(id),
	}
	res, err := app.posts.DeletePost(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling DeleteArticle RPC: %v", err)
	}
	log.Printf("Response from DeleteArticle RPC: %v", res.GetResult())

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) signupUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) signupUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("name", "email", "password")
	form.MaxLength("name", 255)
	form.MaxLength("email", 255)
	form.MatchesPattern("email", forms.EmailRX)
	form.MinLength("password", 10)
	if !form.Valid() {
		app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		return
	}

	req := &authpb.CreateUserRequest{
		User: &authpb.User{
			Name:     form.Get("name"),
			Username:    form.Get("email"),
			Password: form.Get("password"),
		},
	}
	res, _ := app.auth.CreateUser(context.Background(), req)
	if !res.GetStatus() {
		if res.GetResult() == models.ErrDuplicateData.Error() {
			form.Errors.Add("email", "Address is already in use")
			app.render(w, r, "signup.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, errors.New(res.GetResult()))
		}
		return
	}


	app.session.Put(r, "flash", "Your signup was successful. Please log in.")

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) loginUserForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "login.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}

func (app *application) loginUser(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	if !form.Valid() {
		form.Errors.Add("generic", "Email and Password are required")
		app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		return
	}

	req := &authpb.LoginUserRequest{
		User: &authpb.User{
			Username:    form.Get("email"),
			Password: form.Get("password"),
		},
	}

	res, _ := app.auth.AuthUser(context.Background(), req)

	if !res.GetStatus() {
		if res.GetResult() == models.ErrInvalidData.Error() {
			form.Errors.Add("generic", "Email or Password is incorrect")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else if res.GetResult() == models.ErrEmpty.Error() {
			form.Errors.Add("generic", "No user with given Email")
			app.render(w, r, "login.page.tmpl", &templateData{Form: form})
		} else {
			app.serverError(w, errors.New(res.GetResult()))
		}
		return
	}

	app.session.Put(r, "authenticatedUserID", int(res.GetId()))
	app.session.Put(r, "authenticatedUserName",res.GetName() )

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) logoutUser(w http.ResponseWriter, r *http.Request) {
	app.session.Remove(r, "authenticatedUserID")

	app.session.Put(r, "flash", "You've been logged out successfully!")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (app *application) searchForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "search.page.tmpl", &templateData{
		Form: forms.New(nil),
	})
}


func (app *application) editArticleForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}
	req2 := &postpb.GetPostRequest{
		Id: int32(id),
	}
	res2, err := app.posts.GetPost(context.Background(), req2)
	if err != nil {
		log.Fatalf("Error while calling GetArticle RPC: %v", err)
	}
	log.Printf("Response from GetArticle RPC: %s, ArticleID: %v", res2.GetResult(), res2.GetPost().GetId())

	r.PostFormValue("title")
	r.PostFormValue("content")
	form := forms.New(r.PostForm)
	article := &models.Post{
		ID:      id,
		Title:   res2.GetPost().GetTitle(),
		Article: res2.GetPost().GetArticle(),
	}

	form.Set("title", article.Title)
	form.Set("content", article.Article)

	app.render(w, r, "edit.page.tmpl", &templateData{
		Form: form, Post: article,
	})
}

func (app *application) editArticle(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get(":id"))
	if err != nil || id < 1 {
		app.notFound(w)
		return
	}

	err = r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("title", "content")
	form.MaxLength("title", 100)
	if !form.Valid() {
		app.render(w, r, "edit.page.tmpl", &templateData{Form: form, Post: &models.Post{ID: id}})
		return
	}

	req := &postpb.UpdatePostRequest{
		Post: &postpb.Post{
			Title:   form.Get("title"),
			Article: form.Get("content"),
			Id:      int32(id),
		},
	}
	res, err := app.posts.UpdatePost(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling EditArticle RPC: %v", err)
	}
	log.Printf("Response from EditArticle RPC: %v", res.GetResult())
	http.Redirect(w, r, "/", http.StatusSeeOther)

}
