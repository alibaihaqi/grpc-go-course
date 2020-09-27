package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/status"

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
	// doClientStreamingCall(c) // RPC client streaming call
	// doBiDiStreamingCall(c) // RPC bi-directional streaming call
	doErrorUnary(c) // Error unary call RPC
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

func doBiDiStreamingCall(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a BiDirectional Streaming RPC Find Maximum service")

	stream, err := c.FindMaximum(context.Background())
	if err != nil {
		log.Fatalf("error while calling FindMaximum: %v", err)
	}

	numbers := []int32{1, 5, 3, 6, 2, 20}

	waitc := make(chan struct{})

	go func() {
		for _, number := range numbers {
			fmt.Printf("Sending request: %v\n", number)
			stream.Send(&calculatorpb.FindMaximumRequest{
				Number: number,
			})
			time.Sleep(100 * time.Millisecond)
		}

		stream.CloseSend()
	}()

	go func() {
		for {
			msg, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("error while reading stream: %v", err)
				break
			}
			fmt.Printf("FindMaximum Response: %v\n", msg.GetMaxNumber())
		}
		close(waitc)
	}()

	<-waitc
}

func doErrorUnary(c calculatorpb.CalculatorServiceClient) {
	fmt.Println("Starting to do a Unary RPC SquareRoot service")

	// correct call
	doErrorCall(c, 10)

	// incorrect call
	doErrorCall(c, -9)
}

func doErrorCall(c calculatorpb.CalculatorServiceClient, n int32) {
	res, err := c.SquareRoot(context.Background(), &calculatorpb.SquareRootRequest{Number: n})
	if err != nil {
		respErr, ok := status.FromError(err)
		if ok {
			// actual error from gRPC (user error)
			fmt.Printf("Error message from server: %v\n", respErr.Message())
			fmt.Println(respErr.Code())
			if respErr.Code() == codes.InvalidArgument {
				fmt.Println("We probably sent a negative number")
				return
			}
		} else {
			log.Fatalf("Big Error calling SquareRoot RPC: %v", err)
			return
		}
	}
	log.Printf("Response from SquareRoot: %v\n", res.GetNumberRoot())
}
