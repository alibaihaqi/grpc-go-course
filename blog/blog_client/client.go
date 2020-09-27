package main

import (
	"context"
	"fmt"
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
	doUpdateBlogCall(c)
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
		Id:       "5f70d349b2ce160f6c829f61",
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
