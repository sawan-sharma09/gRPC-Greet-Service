// package greet_server
package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"real_grpc/greet/pb"
	"strconv"
	"time"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Greet(ctx context.Context, req *pb.GreetRequest) (*pb.GreetResponse, error) {
	fmt.Println("Greet function was invoked with request-->", req)

	first_name := req.GetGreeting().GetFirstName()
	last_name := req.GetGreeting().GetLastName()
	result := "Hello" + first_name + last_name
	res := &pb.GreetResponse{
		Result: result,
	}
	return res, nil
}

// Server side streaming
func (*server) GreetManyTimes(req *pb.GreetManyTimesRequest, stream pb.GreetService_GreetManyTimesServer) error {
	fmt.Println("GreetManyTimes function was invoked with request-->", req)
	first_name := req.GetGreeting().GetFirstName()
	for i := 0; i < 10; i++ {
		result := "Hello " + first_name + " with number-" + strconv.Itoa(i)
		res := &pb.GreetManyTimesResponse{
			Result: result,
		}
		stream.Send(res)
		time.Sleep(1000 * time.Millisecond)
	}
	return nil
}

// Client side streaming
func (*server) LongGreet(stream pb.GreetService_LongGreetServer) error {
	fmt.Println("LongGreet function was invoked with a streaming request-->", stream)
	result := ""
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(&pb.LongGreetResponse{
				Result: result,
			})
		}
		first_name := req.GetGreeting().GetFirstName()
		result += "Hello " + first_name + " !  "

	}
}

func (*server) BiDiGreet(stream pb.GreetService_BiDiGreetServer) error {
	fmt.Println("Greet BiDi was invoked with a streaming request-->", stream)

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error while receiving data from stream..")
			return err
		}

		first_name := req.GetGreeting().FirstName
		result := "Hello " + first_name + " !  "
		fmt.Println("result->", result)
		sendErr := stream.Send(&pb.GreetBiDiReponse{
			Result: result,
		})
		if sendErr != nil {
			fmt.Println("Error in sending data to stream from server..")
			return err
		}
	}
	return nil

}

func main() {
	fmt.Println("In our server...")
	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatal("Error in listening to port", err)
	}
	fmt.Println("Waiting for data...1")

	s := grpc.NewServer()
	pb.RegisterGreetServiceServer(s, &server{})

	if err2 := s.Serve(lis); err2 != nil {
		log.Fatal("Failed to serve ", err2)
	}
	fmt.Println("Waiting for data...2")
}
