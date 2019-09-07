package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/amazingandyyy/go-grpc-start/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hello I'm a client")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer cc.Close()
	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client: %f", c)
	// doUnary(c)
	doServerStream(c)
}

func doUnary(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andy",
			LastName:  "Chen",
		},
	}
	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling Greet RPC %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStream(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Stream RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Andy",
			LastName:  "Chen",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("Error while calling GreetManyTime RPC %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reachedthe end of the stream
			break
		}
		if err != nil {
			log.Fatalf("Error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTime: %v", msg.GetResult())
	}
}
