package main

import (
	"context"
	"fmt"
	"log"
	"net"

	proto "Grpc/proto"

	"github.com/google/uuid"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"google.golang.org/grpc"
)

type Todo struct {
	ID    string `gorm:"primaryKey"`
	Title string
}

type server struct {
	proto.UnimplementedExampleServer
	db *gorm.DB
}

func main() {

	db, err := gorm.Open(sqlite.Open("todos.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	err = db.AutoMigrate(&Todo{})
	if err != nil {
		log.Fatalf("Failed to migrate database schema: %v", err)
	}

	listener, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatal(err)
	}

	srv := grpc.NewServer()

	proto.RegisterExampleServer(srv, &server{db: db})

	log.Println("Server is running on port 9000")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) ServerReply(ctx context.Context, req *proto.HelloRequest) (*proto.HelloResponse, error) {

	todo := Todo{
		ID:    uuid.New().String(),
		Title: req.Somestring,
	}

	if err := s.db.Create(&todo).Error; err != nil {
		return nil, fmt.Errorf("failed to create todo: %v", err)
	}

	fmt.Println("Received from client:", req.Somestring)
	return &proto.HelloResponse{Reply: "Todo created with title: " + req.Somestring}, nil
}
