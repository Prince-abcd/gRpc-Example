package main

import (
	"context"
	"fmt"
	"log"

	proto "Grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	conn, err := grpc.NewClient("localhost:9000", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	client := proto.NewExampleClient(conn)

	req := &proto.HelloRequest{Somestring: "hi there it me prince"}
	response, err := client.ServerReply(context.TODO(), req)
	if err != nil {
		log.Fatalf("Error calling ServerReply: %v", err)
	}

	// Print the server's response
	fmt.Printf("Server response: %s\n", response.Reply)
}
