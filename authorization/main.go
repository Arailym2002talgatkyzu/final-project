package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"github.com/Arailym2002talgatkyzu/final-project/authorization/authpb"
	"github.com/jackc/pgx/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type Server struct {
	authpb.UnimplementedAuthServiceServer
}

func main() {
	// Connection to DB
	dns := flag.String("dns", "postgres://postgres:650464@localhost:5432/wikipedia", "data and source name")
	db, err :=ConnectDB(*dns)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	user = &UserModel{
		DB: db,
	}

	// Server Creation
	port := "60059" // Port of authorization Microservice
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	authpb.RegisterAuthServiceServer(s, &Server{})
	log.Println("Server is running on port: " + port)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func ConnectDB(dns string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable connect to database: %v\n", err)
		os.Exit(1)
		return nil, err
	}
	return conn, nil
}


func (s *Server) CreateUser(ctx context.Context, req *authpb.CreateUserRequest) (*authpb.CreateUserResponse, error) {
	log.Printf("CreateUser function was invoked with %v \n", req.GetUser().GetUsername())

	errSql := user.Insert(req.GetUser().GetName(), req.GetUser().GetUsername(), req.GetUser().GetPassword())

	res := &authpb.CreateUserResponse{
		Result: "Success",
		Status: true,
	}

	if errSql != nil {
		res.Status = false
		if errors.Is(errSql, ErrDuplicateData) {
			res.Result = ErrDuplicateData.Error()
		} else {
			res.Result = errSql.Error()
		}
	}

	return res, nil
}


func (s *Server) LoginUser(ctx context.Context, req *authpb.LoginUserRequest) (*authpb.LoginUserResponse, error) {
	log.Printf("LoginUser function was invoked with %s \n", req.GetUser().GetUsername())

	id, errSql := user.Authenticate(req.GetUser().GetUsername(), req.GetUser().GetPassword())

	res := &authpb.LoginUserResponse{
		Result: "Success",
		Status: true,
	}

	if errSql != nil {
		res.Id = 0
		res.Status = false
		if errors.Is(errSql, ErrInvalidData) {
			res.Result = ErrInvalidData.Error()
		} else if errors.Is(errSql, ErrEmpty) {
			res.Result = ErrEmpty.Error()
		}
	} else {
		res.Id = int32(id)
	}

	log.Printf("Result: %s", res.GetResult())

	return res, nil
}

func (s *Server) GetUser(ctx context.Context, req *authpb.GetUserRequest) (*authpb.GetUserResponse, error) {
	log.Printf("GetUser function was invoked with id: %s \n", req.GetId())

	resSql, errSql := user.Get(int(req.GetId()))

	res := &authpb.GetUserResponse{
		Result: "Success",
		Status: true,
		User:   &authpb.User{},
	}

	if errSql != nil {
		res.Status = false
		if errors.Is(errSql, ErrEmpty) {
			res.Result = ErrEmpty.Error()
		}
	} else {
		res.GetUser().Id = int32(resSql.ID)
		res.GetUser().Name = resSql.Name
		res.GetUser().Username = resSql.Username
	}

	return res, nil
}