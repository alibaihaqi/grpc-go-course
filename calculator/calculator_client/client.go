package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/alibaihaqi/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Hi, I'm a calculator client")

	cc, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := calculatorpb.NewCalculatorServiceClient(cc)

	// doUnaryCall(c)           // RPC unary call
	// doServerStreamingCall(c) // RPC server streaming call
	doClientStreamingCall(c) // RPC client streaming call
}

func doUnaryCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC Sum service")

	req := &calculatorpb.SumRequest{
		FirstNumber:  3,
		SecondNumber: 10,
	}

	res, err := c.Sum(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling Sum RPC: %v", err)
	}
	log.Printf("Response from Sum: %v", res.SumNumber)
}

func doServerStreamingCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Server Streaming RPC PrimeNumberDecomposition service")

	req := &calculatorpb.PrimeNumberDecompositionRequest{
		Number: 120,
	}

	resStream, err := c.PrimeNumberDecomposition(context.Background(), req)
	if err != nil {
		log.Fatalf("error while calling PrimeNumberDecomposition RPC: %v", err)
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
		log.Printf("Response from PrimeNumberDecomposition: %v", msg.GetResultNumber())
	}
}

func doClientStreamingCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Client Streaming RPC ComputeAverage service")

	stream, err := c.ComputeAverage(context.Background())
	if err != nil {
		log.Fatalf("error while calling ComputeAverage: %v", err)
	}

	numbers := []int32{3, 5, 9, 54, 23}

	// we iterate over our slice and send each message individually
	for _, number := range numbers {
		fmt.Printf("Sending request: %v\n", number)
		stream.Send(&calculatorpb.ComputeAverageRequest{
			Number: number,
		})
		time.Sleep(100 * time.Millisecond)
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("error while receiving response from ComputeAverage: %v", err)
	}
	fmt.Printf("ComputeAverage response: %v\n", res.GetAverageNumber())
}
