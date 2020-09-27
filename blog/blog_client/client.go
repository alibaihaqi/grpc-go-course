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

	doCreateBlogCall(c)
}

func doCreateBlogCall(c blogpb.BlogServiceClient) {
	blog := &blogpb.Blog{
		AuthorId: "Jacky",
		Title:    "My First Blog",
		Content:  "Content of the first Blog",
	}

	res, err := c.CreateBlog(context.Background(), &blogpb.CreateBlogRequest{Blog: blog})
	if err != nil {
		log.Fatalf("Unexpected error:", err)
	}
	fmt.Printf("Blog has been created: %v", res)
}
