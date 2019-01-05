package main

import (
	"context"
	"fmt"
	"github.com/ctaperts/grpc-go-course/calculator/calcpb"
	"google.golang.org/grpc"
	"log"
	"net"
	"strconv"
	"time"
)

type server struct{}

func (*server) Integers(ctx context.Context, req *calcpb.SumRequest) (*calcpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)
	firstvalue := req.GetIntegers().GetNumberOne()
	secondvalue := req.GetIntegers().GetNumberTwo()
	result := firstvalue + secondvalue
	res := &calcpb.SumResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeManyTimes(req *calcpb.PrimeManyTimesRequest, stream calcpb.PrimeService_PrimeManyTimesServer) error {
	fmt.Printf("PrimeManyTimes function was invoked with %v\n", req)
	number := req.GetPrimeInteger().GetNumberOne()
	k := 2
	for n := int(number); n > 1; {
		if n%k == 0 {
			result := "Prime Number Decomposition: " + strconv.Itoa(k)
			res := &calcpb.PrimeManyTimesResponse{
				Result: result,
			}
			stream.Send(res)
			n = n / k
		} else {
			k = k + 1
		}
		time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func main() {
	fmt.Println("Hello world")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterSumServiceServer(s, &server{})
	calcpb.RegisterPrimeServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
