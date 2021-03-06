package main

import (
	"github.com/bmizerany/pat"
	"github.com/justinas/alice"
	"net/http"
)

func (app *application) routes() http.Handler {
	standardMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	dynamicMiddleware := alice.New(app.session.Enable, app.authenticate)

	mux := pat.New()
	mux.Get("/", dynamicMiddleware.ThenFunc(app.home))
	mux.Get("/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createArticleForm))
	mux.Post("/post/create", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.createArticle))
	mux.Get("/post/:id", dynamicMiddleware.ThenFunc(app.showArticle))

	// Custom //ADD Check As for Create with additional parameter author
	mux.Post("/post/delete/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.deleteArticle))
	mux.Get("/post/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.editArticleForm))
	mux.Post("/post/update/:id", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.editArticle))




	// User
	mux.Get("/user/signup", dynamicMiddleware.ThenFunc(app.signupUserForm))
	mux.Post("/user/signup", dynamicMiddleware.ThenFunc(app.signupUser))
	mux.Get("/user/login", dynamicMiddleware.ThenFunc(app.loginUserForm))
	mux.Post("/user/login", dynamicMiddleware.ThenFunc(app.loginUser))
	mux.Post("/user/logout", dynamicMiddleware.Append(app.requireAuthentication).ThenFunc(app.logoutUser))

	fileServer := http.FileServer(http.Dir("app/wikipedia/ui/static/"))
	mux.Get("/static/", http.StripPrefix("/static", fileServer))

	return standardMiddleware.Then(mux)
}
