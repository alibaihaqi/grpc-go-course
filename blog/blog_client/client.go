package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/alibaihaqi/grpc-go-course/blog/blogpb"

	"google.golang.org/grpc"
)

func main() {

	fmt.Println("Blog Client")

	opts := grpc.WithInsecure()

	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := blogpb.NewBlogServiceClient(cc)

	// doCreateBlogCall(c)
	// doReadBlogCall(c)
	// doUpdateBlogCall(c)
	// doDeleteBlogCall(c)
	doListBlogCall(c)
}

func doCreateBlogCall(c blogpb.BlogServiceClient) {
	fmt.Println("Do Create Blog")
	blog := &blogpb.Blog{
		AuthorId: "Jacky",
		Title:    "My First Blog",
		Content:  "Content of the first Blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Blog has been created: %v\n", res)
}

func doReadBlogCall(c blogpb.BlogServiceClient) {
	fmt.Println("Do Read Blog")
	res, err := c.ReadBlog(context.Background(), &blogpb.ReadBlogRequest{
		BlogId: "fake-id",
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Got the blog: %v\n", res)
}

func doUpdateBlogCall(c blogpb.BlogServiceClient) {
	fmt.Println("Do Update Blog")

	updateBlog := &blogpb.Blog{
		Id:       "fake-id",
		AuthorId: "Changed Author",
		Title:    "My First Updated Title",
		Content:  "My First Updated Content",
	}

	res, err := c.UpdateBlog(context.Background(), &blogpb.UpdateBlogRequest{
		Blog: updateBlog,
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Got the updated blog: %v\n", res)
}

func doDeleteBlogCall(c blogpb.BlogServiceClient) {
	fmt.Println("Do Delete Blog")
	res, err := c.DeleteBlog(context.Background(), &blogpb.DeleteBlogRequest{
		BlogId: "fake-id",
	})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
	}
	fmt.Printf("Success delete blog with id: %v\n", res)
}

func doListBlogCall(c blogpb.BlogServiceClient) {
	fmt.Println("Do List Blog")
	resStream, err := c.ListBlog(context.Background(), &blogpb.ListBlogRequest{})
	if err != nil {
		log.Fatalf("Unexpected error: %v\n", err)
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
		log.Printf("Response from ListBlogCall: %v", msg.GetBlog())
	}
}
