// package greet_client
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"real_grpc/greet/pb"
	"time"

	"google.golang.org/grpc"
)

// func MyClient() {
func main() {
	fmt.Println("Hello I'm a client..")
	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("Couldn't connect to client", err)
	}

	defer cc.Close()

	c := pb.NewGreetServiceClient(cc)
	// fmt.Printf("Created client--> %f\n", c)
	fmt.Println("Created a client")

	// DoUnary(c)
	// DoServerStreaming(c)
	// DoClientStreaming(c)
	DoBiDiStreaming(c)
}

func DoUnary(c pb.GreetServiceClient) {
	fmt.Println("Stating to do an Unary RPC..")
	req := &pb.GreetRequest{
		Greeting: &pb.Greeting{
			FirstName: "Neal",
			LastName:  "Den Exter",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		fmt.Println("Error while calling Greet RPC: ", err, "response is-->", res)
	}
	fmt.Println("Response from Greet: ", res.Result)
}

func DoServerStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Server side streaming RPC..")
	req := &pb.GreetManyTimesRequest{
		Greeting: &pb.Greeting{
			FirstName: "Brad",
			LastName:  "Heritage",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		fmt.Println("Error in calling GreetManyTimes Rpc", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// it means we have reached the end of the steam and there are no more values to be received
			break
		}
		if err != nil {
			fmt.Println("Error while reading stream..!!")
		}
		fmt.Println("Response from GreetManyTimes", msg.GetResult())
	}
}

func DoClientStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting to do a Client side streaming RPC...")

	requests := []*pb.LongGreetRequest{
		{
			Greeting: &pb.Greeting{
				FirstName: "Mike",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Tyson",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "John",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Wick",
			},
		},
	}
	stream, err := c.LongGreet(context.Background())
	if err != nil {
		fmt.Println("Error while calling LongGreet..", err)
	}
	for _, req := range requests {
		fmt.Println("Sending request-->", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		fmt.Println("Error in closeAndRecv..", err)
	}

	fmt.Println("LongGreet Response..\n", res)
}

func DoBiDiStreaming(c pb.GreetServiceClient) {
	fmt.Println("Starting a Bi-Directional streaming rpc..")

	requests := []*pb.GreetBiDiRequest{
		{
			Greeting: &pb.Greeting{
				FirstName: "Finn",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Allen",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Neal",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Den",
			},
		},
		{
			Greeting: &pb.Greeting{
				FirstName: "Exter",
			},
		},
	}

	stream, err := c.BiDiGreet(context.Background())
	if err != nil {
		fmt.Println("Error in client side stream generation..")
	}

	waitc := make(chan struct{})
	go func() {
		for _, req := range requests {
			fmt.Println("Sending requests", req)
			stream.Send(req)
			// time.Sleep(1000 * time.Millisecond)
		}
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				fmt.Println("All the data in the stream has been printed..")
				// break
				close(waitc)
			}
			if err != nil {
				fmt.Println("Error in stream receive in client..")
			}
			fmt.Println("Received :", res.GetResult())
		}
		// close(waitc)
	}()
	<-waitc
}
