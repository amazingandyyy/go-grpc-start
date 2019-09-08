package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

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
	// doServerStream(c)
	// doClientStream(c)
	doBiDiStream(c)
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

func doClientStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Printf("Start to do a client stream RPC...")

	requests := []*calculatorpb.ComputeAverageRequest{
		&calculatorpb.ComputeAverageRequest{
			Number: 1,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 2,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 3,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 4,
		},
		&calculatorpb.ComputeAverageRequest{
			Number: 50,
		},
	}
	reqStream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("Error while calling ComputeAverage on the server %v\n", err)
	}
	for _, req := range requests {
		reqStream.Send(req)
	}
	res, err := reqStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Error while receiving result from ComputeAverage %v\n", err)
	}
	fmt.Printf("Average from server: %v\n", res.Result)
}

func doBiDiStream(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Start to do a BiDi stream RPC")
	requests := []*calculatorpb.FindMaximumRequest{
		&calculatorpb.FindMaximumRequest{
			Number: 1,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 5,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 3,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 6,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 2,
		},
		&calculatorpb.FindMaximumRequest{
			Number: 20,
		},
	}
	reqStream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("Error when sending to server %v", err)
	}
	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Printf("Sending message %v\n", req)
			reqStream.Send(req)
			time.Sleep(100 * time.Millisecond)
		}
		reqStream.CloseSend()
	}()

	go func() {
		for {
			res, err := reqStream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Error when receiving from server: %v", err)
			}
			max := res.Max
			fmt.Printf("New Max is: %v\n", max)
		}
		close(waitc)
	}()

	<-waitc
}
