package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/alibaihaqi/grpc-go-course/greet/greetpb"
	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Hello I'm a client")

	cc, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := greetpb.NewGreetServiceClient(cc)
	// fmt.Println("Created client: %f", c)

	// doUnaryCall(c)
	// doServerStreamingCall(c)
	// doClientStreamingCall(c)
	doBiDirectionalCall(c)
}

func doUnaryCall(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Unary RPC...")
	req := &greetpb.GreetRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jacky",
			LastName:  "Hugo",
		},
	}

	res, err := c.Greet(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Greet RPC: %v", err)
	}
	log.Printf("Response from Greet: %v", res.Result)
}

func doServerStreamingCall(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC...")

	req := &greetpb.GreetManyTimesRequest{
		Greeting: &greetpb.Greeting{
			FirstName: "Jacky",
			LastName:  "Hugo",
		},
	}

	resStream, err := c.GreetManyTimes(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling GreetManyTimes RPC:, %v", err)
	}
	for {
		msg, err := resStream.Recv()
		if err == io.EOF {
			// we've reached the end of the stream
			break
		}
		if err != nil {
			log.Fatalf("error while reading stream: %v", err)
		}
		log.Printf("Response from GreetManyTimes: %v", msg.GetResult())
	}
}

func doClientStreamingCall(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC...")

	requests := []*greetpb.LongGreetRequest{
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jacky",
				LastName:  "Hugo",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
				LastName:  "Chris",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Mark",
				LastName:  "Joe",
			},
		},
		&greetpb.LongGreetRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Piper",
				LastName:  "Hugo",
			},
		},
	}

	stream, err := c.LongGreet(context.Background())
	if err != nil {
		log.Fatalf("error while calling LongGreet: %v", err)
	}

	// we iterate over our slice and send each message individually
	for _, req := range requests {
		fmt.Printf("Sending request: %v\n", req)
		stream.Send(req)
		time.Sleep(1000 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from LongGreet: %v", err)
	}
	fmt.Printf("LongGreet response: %v", res)
}

func doBiDirectionalCall(c greetpb.GreetServiceClient) {
	fmt.Println("Starting to do a Bi-Directional Streaming RPC...")

	stream, err := c.GreetEveryone(context.Background())
	if err != nil {
		log.Fatalf("error while calling GreetEveryone: %v", err)
		return
	}

	requests := []*greetpb.GreetEveryoneRequest{
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Jacky",
				LastName:  "Hugo",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "John",
				LastName:  "Doe",
			},
		},
		&greetpb.GreetEveryoneRequest{
			Greeting: &greetpb.Greeting{
				FirstName: "Lucy",
				LastName:  "Chris",
			},
		},
	}

	waitc := make(chan struct{})

	// we send a buch of messages to the client (go routine)
	go func() {
		// function to send bunch of messages

		for _, req := range requests {
			fmt.Printf("Sending request: %v\n", req)
			stream.Send(req)
			time.Sleep(1000 * time.Millisecond)
		}

		stream.CloseSend()
	}()

	// we receive a bunch of messages from the client (go routine)
	go func() {
		// function to receive bunch of messages

		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}
			log.Printf("Response from GreetEveryone: %v\n", msg.GetResult())
		}
		close(waitc)
	}()

	// block until everything is done
	<-waitc
}
