#!/bin/bash

# Greet Service
protoc --go_out=plugins=grpc:. greet/greetpb/greet.proto

# Calculator Service
protoc --go_out=plugins=grpc:. calculator/calculatorpb/calculator.proto

# Blog Service
protoc --go_out=plugins=grpc:. blog/blogpb/blog.proto