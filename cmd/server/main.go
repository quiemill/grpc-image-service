package main

import (
	"context"
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"log"
	"net"
	"os"
	"path/filepath"

	"google.golang.org/grpc"
)

const storageDir = "../../storage"

type ImageServer struct {
	pb.UnimplementedImageServiceServer
}

func (s *ImageServer) UploadImage(ctx context.Context, req *pb.ImageBatch) (*pb.UploadResponse, error) {
	if len(req.Images) == 0 {
		return &pb.UploadResponse{Success: false, Info: "Нет файлов для загрузки"}, nil
	}
	for _, img := range req.Images {
		path := filepath.Join(storageDir, img.Filename)
		path = filepath.Clean(path)
		err := os.WriteFile(path, img.Data, 0644)
		if err != nil {
			return nil, err
		}
	}
	return &pb.UploadResponse{Success: true, Info: "Файлы загружены"}, nil
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
	pb.RegisterImageServiceServer(grpcServer, &ImageServer{})

	fmt.Println("Сервер запущен на порту 50051")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Ошибка работы сервера: %v", err)
	}
}
