package main

import (
	"context"
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"log"
	"net"
	"os"
	"path/filepath"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
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

func (s *ImageServer) ListImages(ctx context.Context, _ *emptypb.Empty) (*pb.ImageList, error) {
	files, err := os.ReadDir(storageDir)
	if err != nil {
		return nil, err
	}
	var images []*pb.ImageInfo
	now := time.Now().Format("2006-01-02 15:04:05")
	for _, file := range files {
		images = append(images, &pb.ImageInfo{
			Filename:  file.Name(),
			CreatedAt: now,
			UpdatedAt: now,
		})
	}
	return &pb.ImageList{Images: images}, nil
}

func (s *ImageServer) DownloadImage(ctx context.Context, req *pb.ImageRequest) (*pb.ImageBatch, error) {
	if len(req.Filenames) == 0 {
		return &pb.ImageBatch{}, nil
	}
	images := []*pb.ImageData{}
	for _, filename := range req.Filenames {
		path := filepath.Join(storageDir, filename)
		data, err := os.ReadFile(path)
		if err == nil {
			images = append(images, &pb.ImageData{Filename: filename, Data: data})
		}
	}
	if len(images) == 0 {
		return &pb.ImageBatch{}, nil
	}
	return &pb.ImageBatch{Images: images}, nil
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
