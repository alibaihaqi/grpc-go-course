package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math"
	"net"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

func (*server) ComputeAverage(stream calculatorpb.CalculatorService_ComputeAverageServer) error {
	fmt.Println("PrimeNumberDecomposition function was invoked")

	sum := int32(0)
	count := int32(0)

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			// we have finished reading the client streaming
			result := float32(sum) / float32(count)
			res := &calculatorpb.ComputeAverageResponse{
				AverageNumber: result,
			}
			fmt.Printf("Response: %v\n", res.GetAverageNumber())
			return stream.SendAndClose(res)
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
		}

		sum += req.GetNumber()
		count++
	}
}

func (*server) FindMaximum(stream calculatorpb.CalculatorService_FindMaximumServer) error {
	fmt.Println("FindMaximum function was invoked")

	var maxNum int32 = 0

	for {
		req, err := stream.Recv()

		if err == io.EOF {
			return nil
		}
		if err != nil {
			log.Fatalf("Error while reading client stream: %v", err)
			return err
		}

		num := req.GetNumber()
		if num > maxNum {
			maxNum = num
			sendErr := stream.Send(&calculatorpb.FindMaximumResponse{
				MaxNumber: maxNum,
			})
			if sendErr != nil {
				log.Fatalf("Error while sending data to client: %v", sendErr)
				return sendErr
			}
		}
	}
}

func (*server) SquareRoot(ctx context.Context, req *calculatorpb.SquareRootRequest) (*calculatorpb.SquareRootResponse, error) {
	fmt.Printf("SquareRoot function was invoked with %v\n", req)

	num := req.GetNumber()

	if num > 0 {
		result := math.Sqrt(float64(num))

		res := &calculatorpb.SquareRootResponse{
			NumberRoot: result,
		}

		return res, nil
	}
	return nil, status.Errorf(
		codes.InvalidArgument,
		fmt.Sprintf("Received a negative number: %v", num),
	)
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
