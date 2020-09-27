package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/alibaihaqi/grpc-go-course/blog/blogpb"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var collection *mongo.Collection

type server struct{}

type blogItem struct {
	ID       bson.ObjectId `bson:"_id,omitempty"`
	AuthorID string        `bson:"author_ud"`
	Title    string        `bson:"title"`
	Content  string        `bson:"content"`
}

func main() {
	// if we crash the go code, we get the file name and line number
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Blog Service Started")

	fmt.Println("Connecting to MongoDB")
	client, mongoErr := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	clientErr := client.Connect(ctx)
	if clientErr != nil {
		log.Fatal(clientErr)
	}

	collection = client.Database("mydb").Collection("blog")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}

	s := grpc.NewServer(opts...)
	blogpb.RegisterBlogServiceServer(s, &server{})

	go func() {
		fmt.Println("Strarting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to server: %v", err)
		}
	}()

	// Wait for Ctrl C to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing MongoDB connection")
	client.Disconnect(ctx)
	fmt.Println("End of Program")
}
