package main

import (
	"fmt"
	"github.com/ctaperts/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
	"log"
	// "net"
)

func main() {

	fmt.Println("Client")
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}
	defer conn.Close()

	c := greetpb.NewGreetServiceClient(conn)
	fmt.Printf("Created client: %f", c)

}
