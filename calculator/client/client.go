package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/amazingandyyy/go-grpc-start/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Calculator Client init")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect %v", err)
	}

	defer cc.Close()
	c := calculatorpb.NewCalculatorServiceClient(cc)
	// doUnary(c)
	doServerStream(c)
}

func doUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &calculatorpb.CalculatorRequest{
		Calculating: &calculatorpb.Calculating{
			NumberOne: 3,
			NumberTwo: 10,
		},
	}
	res, err := c.Calculate(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator RPC %v", err)
	}
	log.Printf("Response from Calculator: %v", res)
}

func doServerStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Start to do a server stream RPC...")
	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Calculating: &calculatorpb.Calculating{
			NumberOne: 120,
		},
	}
	resStream, err := c.PrimeNumberDecompose(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Calculator server stream RPC %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream %v", err)
		}
		fmt.Printf("Response from server %v\n", msg)
	}
}
