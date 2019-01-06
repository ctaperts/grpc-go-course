package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"io"
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

	c := calcpb.NewCalcServiceClient(conn)
	// Unary RPC call
	doUnary(c)

	cs := calcpb.NewCalcServiceClient(conn)
	// Server Streaming RPC call
	doServerStreaming(cs)

}

func doUnary(c calcpb.CalcServiceClient) {
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

func doServerStreaming(c calcpb.CalcServiceClient) {
	fmt.Println("Starting server streaming")

	req := &calcpb.PrimeManyTimesRequest{
		PrimeInteger: &calcpb.PrimeInteger{
			NumberOne: 120465,
		},
	}
	resStream, err := c.PrimeManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeInteger server streaming RPC: %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// End of stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from PrimeIntegerManyTimes: %v", msg.GetResult())
	}
}
