package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
	// "io"
	"log"
	// "time"
)

func main() {

	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	conn, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := blogpb.NewBlogServiceClient(conn)

	fmt.Println("Creating the blog")
	blog := &blogpb.Blog{
		AuthorId: "Colby",
		Title:    "My first blog",
		Content:  "Content of blog",
	}
	createBlogRes, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Blog has been created: %v", createBlogRes)
}
