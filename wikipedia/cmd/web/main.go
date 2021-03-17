// In terminal: go run ./cmd/web -addr=":4000"
package main

import (
	"flag"
	"github.com/Arailym2002talgatkyzu/final-project/authorization/authpb"
	"github.com/Arailym2002talgatkyzu/final-project/post_db/postpb"
	"github.com/golangcollege/sessions"
	"google.golang.org/grpc"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

type contextKey string

const contextKeyIsAuthenticated = contextKey("isAuthenticated")

type application struct {
	errorLog      *log.Logger
	infoLog       *log.Logger
	session       *sessions.Session
	posts         postpb.PostServiceClient
	auth          authpb.AuthServiceClient
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	secret := flag.String("secret", "s6Ndh+pPbnzHbS*+9Pk8qGWhTzbpa@ge", "Secret key")
	flag.Parse()

	// Loggers init
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	// Post_DB Microservice
	port := "60051"
	conn1, err := grpc.Dial("127.0.0.1:"+port, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn1.Close()
	postsDBService := postpb.NewPostServiceClient(conn1)


	// Auth Service
	port3 := "60059"
	conn3, err := grpc.Dial("127.0.0.1:"+port3, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn3.Close()
	authService := authpb.NewAuthServiceClient(conn3)

	// Template Cache init
	templateCache, err := newTemplateCache("app/wikipedia/ui/html/")
	if err != nil {
		errorLog.Fatal(err)
	}

	// Session init
	session := sessions.New([]byte(*secret))
	session.Lifetime = 12 * time.Hour

	app := &application{
		errorLog:      errorLog,
		infoLog:       infoLog,
		session:       session,
		posts:         postsDBService,
		auth:          authService,
		templateCache: templateCache,
	}

	srv := &http.Server{
		Addr:     *addr,
		ErrorLog: errorLog,
		Handler:  app.routes(),
	}

	infoLog.Printf("Starting a server on port%v", *addr)
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
}
