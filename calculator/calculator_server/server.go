package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/alibaihaqi/grpc-go-course/calculator/calculatorpb"

	"google.golang.org/grpc"
)

type server struct{}

func (*server) Sum(ctx context.Context, req *calculatorpb.SumRequest) (*calculatorpb.SumResponse, error) {
	fmt.Printf("Sum function was invoked with %v\n", req)

	firstNumber := req.GetFirstNumber()
	secondNumber := req.GetSecondNumber()

	result := firstNumber + secondNumber

	res := &calculatorpb.SumResponse{
		SumNumber: result,
	}

	return res, nil
}

func (*server) PrimeNumberDecomposition(req *calculatorpb.PrimeNumberDecompositionRequest, stream calculatorpb.CalculatorService_PrimeNumberDecompositionServer) error {
	fmt.Printf("PrimeNumberDecomposition function was invoked with %v\n", req)

	num := req.GetNumber()
	var pnum int32 = 2
	for num > 1 {
		if num%pnum == 0 {
			fmt.Printf("current factor: %v\n", pnum)
			res := &calculatorpb.PrimeNumberDecompositionResponse{
				ResultNumber: fmt.Sprint(pnum),
			}
			num = num / pnum
			stream.Send(res)
			time.Sleep(1000 * time.Millisecond)
		} else {
			pnum++
		}
	}

	return nil
}

func main() {
	fmt.Println("Hello, I'm calculator server")

	lis, err := net.Listen("tcp", "localhost:50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	s := grpc.NewServer()
	calculatorpb.RegisterCalculatorServiceServer(s, &server{})

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to server: %v", err)
	}
}
