package main

import (
	"bufio"
	"context"
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

func uploadImage(client pb.ImageServiceClient) {
	fmt.Print("Введите пути к изображениям через запятую: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	paths := strings.Split(strings.TrimSpace(input), ",")
	var images []*pb.ImageData
	for _, path := range paths {
		path = strings.TrimSpace(path)
		filename := filepath.Base(path)
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Ошибка открытия файла %s: %v\n", path, err)
			continue
		}
		defer file.Close()
		data, err := io.ReadAll(file)
		if err != nil {
			fmt.Printf("Ошибка чтения файла %s: %v\n", path, err)
			continue
		}
		images = append(images, &pb.ImageData{
			Filename: filename,
			Data:     data,
		})
	}
	if len(images) == 0 {
		fmt.Println("Нет файлов для загрузки.")
		return
	}
	resp, err := client.UploadImage(context.Background(), &pb.ImageBatch{Images: images})
	if err != nil {
		fmt.Println("Ошибка загрузки изображений:", err)
		return
	}
	fmt.Println("Результат загрузки:", resp.Info)
}

func listImages(client pb.ImageServiceClient) {
	resp, err := client.ListImages(context.Background(), &emptypb.Empty{})
	if err != nil {
		fmt.Println("Ошибка получения списка изображений:", err)
		return
	}
	fmt.Println("\nСписок изображений на сервере:")
	for _, img := range resp.Images {
		fmt.Printf("Файл: %s | Создано: %s | Обновлено: %s\n", img.Filename, img.CreatedAt, img.UpdatedAt)
	}
}

func downloadImage(client pb.ImageServiceClient) {
	fmt.Print("Введите имя файла для скачивания: ")
	var filename string
	fmt.Scanln(&filename)
	resp, err := client.DownloadImage(context.Background(), &pb.ImageRequest{Filenames: []string{filename}})
	if err != nil || len(resp.Images) == 0 {
		fmt.Println("Файл не найден")
		return
	}
	image := resp.Images[0]
	os.WriteFile("downloaded_"+image.Filename, image.Data, 0644)
	fmt.Println("Скачано:", image.Filename)
}

func main() {
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("1.Загрузить изображения")
		fmt.Println("2.Посмотреть список изображений")
		fmt.Println("3.Скачать изображения")
		fmt.Println("4.Выйти")
		fmt.Print("Ваш выбор: ")

		reader := bufio.NewReader(os.Stdin)
		ch, _ := reader.ReadString('\n')
		ch = strings.TrimSpace(ch)

		switch ch {
		case "1":
			uploadImage(client)
		case "2":
			listImages(client)
		case "3":
			downloadImage(client)
		case "4":
			fmt.Println("Завершение работы")
			return
		default:
			fmt.Println("Попробуйте снова")
		}
	}
}
