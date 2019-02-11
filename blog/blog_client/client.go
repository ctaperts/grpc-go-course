package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/blog/blogpb"
	"google.golang.org/grpc"
	// "google.golang.org/grpc/codes"
	// "google.golang.org/grpc/credentials"
	// "google.golang.org/grpc/status"
	"io"
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
	fmt.Printf("Blog has been created: %v\n", createBlogRes)
	blogID := createBlogRes.GetBlog().GetId()

	// read Blog
	fmt.Println("Reading the blog")

	_, err = c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{BlogId: "1123123"})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	readBlogReq := &blogpb.ReadBlogRequest{BlogId: blogID}
	readBlogRes, err := c.ReadBlog(context.Background(), readBlogReq)
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was read: %v\n", readBlogRes)

	// update Blog
	newBlog := &blogpb.Blog{
		Id:       blogID,
		AuthorId: "Colby Taperts",
		Title:    "My first blog, ok (edited)",
		Content:  "Content of blog",
	}

	updateRes, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error happened while reading: %v\n", err)
	}

	fmt.Printf("Blog was updated: %v\n", updateRes)

	// delete Blog
	deleteRes, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("Error happened while deleting: %v\n", err)
	} else {
		fmt.Printf("Blog was deleted: %v\n", deleteRes)
	}

	// 	deleteRes, err = c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{BlogId: blogID})
	// 	if err != nil {
	// 		fmt.Printf("Error happened while deleting: %v\n", err)
	// 	} else {
	// 		fmt.Printf("Blog was deleted: %v\n", deleteRes)
	// 	}
	//
	// listBlogs

	stream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}
}
