package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/calculator/calcpb"
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

	c := calcpb.NewSumServiceClient(conn)
	// fmt.Printf("Created client: %f", c)
	doUnary(c)

}

func doUnary(c calcpb.SumServiceClient) {
	fmt.Println("Starting unary rpc...")
	req := &calcpb.SumRequest{
		Integers: &calcpb.Integers{
			NumberOne: 3,
			NumberTwo: 10,
		},
	}
	res, err := c.Integers(context.Background(), req)
	if err != nil {
		log.Fatalf("error while call Integers RPC: %v", err)
	}
	log.Printf("Response from Integers: %v", res)
}
