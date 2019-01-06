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
	"math"
	"net"
	"strconv"
	// "time"
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

func (*server) PrimeManyTimes(req *calcpb.PrimeManyTimesRequest, stream calcpb.CalcService_PrimeManyTimesServer) error {
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
			k++
		}
		// time.Sleep(10 * time.Millisecond)
	}
	return nil
}

func (*server) AverageLong(stream calcpb.CalcService_AverageLongServer) error {
	fmt.Printf("AverageLong function was invoked with stream request \n")
	var total int32
	var result float32
	var amount_of_numbers int64
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			// finished reading the client stream
			result = float32(total) / float32(amount_of_numbers)
			return stream.SendAndClose(&calcpb.AverageResponse{
				Result: result,
			})

		}
		if err != nil {
			log.Fatalf("Error while reading the client stream: %v", err)
		}

		number := req.GetIntegers().GetNumberOne()
		amount_of_numbers++
		total = total + number
	}
	return nil
}

func (*server) FindMaximum(stream calcpb.CalcService_FindMaximumServer) error {
	fmt.Printf("FindMaximum function was invoked with stream request \n")
	var highest_number int32
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("error while reading client stream: %v", err)
			return err
		}
		number := req.GetIntegers().GetNumberOne()
		if number > highest_number {
			highest_number = number
			err = stream.Send(&calcpb.FindMaximumResponse{
				Result: highest_number,
			})
			if err != nil {
				log.Fatalf("error while receiving client stream: %v", err)
				return err
			}
		}
	}
	return nil
}

func (*server) SquareRoot(ctx context.Context, req *calcpb.SquareRootRequest) (*calcpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invoked with %v\n", req)
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number: %v", number),
		)
	}
	return &calcpb.SquareRootResponse{
		Result: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Starting server...")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calcpb.RegisterCalcServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
