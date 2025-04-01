package main

import (
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
)

const storageDir = "../images"

type server struct {
	pb.UnimplementedImageServiceServer
}

func main() {

	if err := os.MkdirAll(storageDir, os.ModePerm); err != nil {
		log.Fatalf("Ошибка создания папки: %v", err)
	}

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterImageServiceServer(grpcServer, &server{})

	fmt.Println("Сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}
