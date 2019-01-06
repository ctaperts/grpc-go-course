package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
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

	// Server Streaming RPC call
	doServerStreaming(c)

	// Client Streaming RPC call
	doClientStreaming(c)

	// Bi-directional streaming RPC call
	doBiDirectionalStreaming(c)

	// Error codes
	doSquareUnary(c)

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

func doClientStreaming(c calcpb.CalcServiceClient) {
	fmt.Println("Starting client streaming")

	requests := []*calcpb.AverageRequest{
		&calcpb.AverageRequest{
			Integers: &calcpb.Integers{
				NumberOne: 1,
			},
		},
		&calcpb.AverageRequest{
			Integers: &calcpb.Integers{
				NumberOne: 2,
			},
		},
		&calcpb.AverageRequest{
			Integers: &calcpb.Integers{
				NumberOne: 3,
			},
		},
		&calcpb.AverageRequest{
			Integers: &calcpb.Integers{
				NumberOne: 4,
			},
		},
	}

	stream, err := c.AverageLong(context.Background())
	if err != nil {
		log.Fatalf("error while client streaming: %v", err)
	}
	// we iterate over our slive and send each individually
	for _, req := range requests {
		fmt.Printf("Sending req: %v\n", req)
		stream.Send(req)
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while recieving response from AverageLong %v", err)
	}
	fmt.Printf("AverageLong response: %v\n", res)

}

func doBiDirectionalStreaming(c calcpb.CalcServiceClient) {
	fmt.Println("Starting bi-directional streaming")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error while creating stream: %v", err)
	}

	numbers := []int32{4, 6, 8, 1, 3, 32}
	waitc := make(chan struct{})
	go func() {
		for _, req := range numbers {
			fmt.Printf("Sending integer: %v\n", req)
			stream.Send(&calcpb.FindMaximumRequest{Integers: &calcpb.Integers{NumberOne: req}})
			time.Sleep(100 * time.Millisecond)
		}
		stream.CloseSend()
	}()
	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error while recieving: %v", err)
				break
			}
			fmt.Printf("Response: %v\n", res)
		}
		close(waitc)
	}()

	<-waitc
}

func doSquareUnary(c calcpb.CalcServiceClient) {
	fmt.Println("Starting Square root unary directional streaming")
	// correct call
	var number int32 = 20
	doSquareRootCall(c, number)
	number = -10
	doSquareRootCall(c, number)
}

func doSquareRootCall(c calcpb.CalcServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calcpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from grpc(user error)
			fmt.Println(respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("Error: cannot square a negative error")
			}
		} else {
			log.Fatalf("Error call SquareRoot: %v", err)
		}

		// error call
	}
	fmt.Printf("Result of square root of %v: %v", n, res.GetResult())
}
