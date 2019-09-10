package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"

	"google.golang.org/grpc/codes"

	"github.com/amazingandyyy/go-grpc-start/calculator/calculatorpb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

type server struct{}

func (*server) Calculate(ctx context.Context, req *calculatorpb.CalculatorRequest) (*calculatorpb.CalculatorResponse, error) {
	fmt.Printf("Calculator function was invoked with %v\n", req)
	firstNum := req.GetCalculating().GetNumberOne()
	secondNumber := req.GetCalculating().GetNumberTwo()
	result := firstNum + secondNumber
	res := &calculatorpb.CalculatorResponse{
		Result: result,
	}
	return res, nil
}

func (*server) PrimeNumberDecompose(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecomposeServer) error {
	fmt.Printf("PrimeNumberDecompose function was invoked with %v\n", req)
	number := req.GetCalculating().GetNumberOne()
	k := int32(2)
	for number > 1 {
		if number%k == 0 {
			result := k
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				Result: result,
			}
			stream.Send(res)
			number = number / k
		} else {
			k = k + 1
		}
	}
	return nil
}

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	sum := int32(0)
	total := int32(0)
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			stream.SendAndClose(&calculatorpb.ComputeAverageResponse{
				Result: sum / total,
			})
			break
		}
		if err != nil {
			log.Fatalf("Error when reading from client %v", err)
		}
		sum += res.GetNumber()
		total += int32(1)
	}
	return nil
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum was invokded!")
	max := int32(0)
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while receiving from client: %v", err)
		}
		num := req.GetNumber()
		if num > max {
			max = num
			err := stream.Send(&calculatorpb.FindMaximumResponse{
				Max: max,
			})
			if err != nil {
				log.Fatalf("Error while sending to client: %v", err)
			}
		}
	}
	return nil
}

func (*server) SquareRoot(c context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Println("Received SquareRoot RPC")
	number := req.GetNumber()
	if number < 0 {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Received a negative number %v", number),
		)
	}
	return &calculatorpb.SquareRootResponse{
		NumberRoot: math.Sqrt(float64(number)),
	}, nil
}

func main() {
	fmt.Println("Calculator init")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
