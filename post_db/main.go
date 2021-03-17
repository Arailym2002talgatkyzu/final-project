package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/Arailym2002talgatkyzu/final-project/post_db/postpb"
	_ "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type Server struct {
	postpb.UnimplementedPostServiceServer
}

func main() {
	// Connection to DB

	dns := flag.String("dns", "postgres://postgres:650464@localhost:5432/wikipedia", "Postgre data source name")
	db, err := connectDB(*dns)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	post = &PostModel{
		DB: db,
	}

	// Server Creation
	port := "60051" // Port of post_DB Microservice
	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v\n", err)
	}
	s := grpc.NewServer()
	postpb.RegisterPostServiceServer(s, &Server{})
	log.Println("Server is running on port: " + port)
	if err := s.Serve(l); err != nil {
		log.Fatalf("Failed to serve: %v\n", err)
	}
}

func connectDB(dns string) (*pgxpool.Pool, error) {
	conn, err := pgxpool.Connect(context.Background(), dns)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable connect to database: %v\n", err)
		os.Exit(1)
		return nil, err
	}
	return conn, nil
}

func (s *Server) GetPost(ctx context.Context, req *postpb.GetPostRequest) (*postpb.GetPostResponse, error) {
	log.Printf("GetPost function was invoked with %v \n", req)

	postRes, err := post.Get(int(req.GetId()))

	res := &postpb.GetPostResponse{
		Post: &postpb.Post{
			Id:       int32(postRes.ID),
			Authorname: postRes.AuthorName,
			Title:    postRes.Title,
			Article:  postRes.Article,
			Published:  postRes.Published.String(),

		},
		Result: "Success",
	}

	return res, err
}


func (s *Server) GetPosts(req *postpb.GetPostsRequest, stream postpb.PostService_GetPostsServer) error {
	log.Printf("GetPosts function was invoked\n")

	arr, err := post.GetAll()
	if err != nil {
		return err
	}

	for _, articleRes := range arr {
		res := &postpb.GetPostsResponse{
			Post: &postpb.Post{
				Id:       int32(articleRes.ID),
				Authorname: articleRes.AuthorName,
				Authorid: int32(articleRes.AuthorId),
				Title:    articleRes.Title,
				Article:  articleRes.Article,
				Published:  articleRes.Published.String(),
			},
		}

		if err := stream.Send(res); err != nil {
			log.Fatalf("Error while sending GetPosts responses: %v", err.Error())
		}
	}

	return nil
}

func (s *Server) DeletePost(ctx context.Context, req *postpb.DeletePostRequest) (*postpb.DeletePostResponse, error) {
	log.Printf("DeletePostfunction was invoked with %v \n", req)

	res := &postpb.DeletePostResponse{
		Result: "Success",
	}

	err := post.Delete(int32(req.GetId()))
	if err != nil {
		res.Result = "Fail"
		log.Fatalf("Error DeletePost: %v", err.Error())
	}

	return res, nil
}



func (s *Server) InsertPost(ctx context.Context, req *postpb.InsertPostRequest) (*postpb.InsertPostResponse, error) {
	log.Printf("InsertPost function was invoked with %v \n", req)

	id, err := post.Add(req.GetPost().GetTitle(), req.GetPost().GetArticle(),  req.GetPost().GetAuthorname(), req.GetPost().GetAuthorid())

	res := &postpb.InsertPostResponse{
		Id:     int32(id),
		Result: "Success",
	}

	if err != nil {
		res.Result = "Fail"
		log.Fatalf("Error InsertPost: %v", err.Error())
	}

	return res, nil
}

func (s *Server) UpdatePost(ctx context.Context, req *postpb.UpdatePostRequest) (*postpb.UpdatePostResponse, error) {
	log.Printf("EditArticle function was invoked with %v \n", req)
	res := &postpb.UpdatePostResponse{
		Result: "Success",
	}

	err := post.Update(req.GetPost().GetTitle(), req.GetPost().GetArticle(), int(req.GetPost().GetId()))
	if err != nil {
		res.Result = "Fail"
		log.Fatalf("Error Update Article: %v", err.Error())
	}

	return res, nil
}
