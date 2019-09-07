package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	"github.com/amazingandyyy/go-grpc-start/calculator/calculatorpb"
	"google.golang.org/grpc"
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
