package main

import (
	"context"
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
	// fmt.Printf("Created client: %f", c)
	doUnary(c)

}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting unary rpc...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Colby",
			LastName:  "Taperts",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while call Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res)
}
