package main

import (
	"bufio"
	"context"
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func uploadImage(client pb.ImageServiceClient) {
	fmt.Println("Введите пути к изображениям через запятую:")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Ошибка чтения ввода:", err)
		return
	}
	input = strings.TrimSpace(input)
	paths := strings.Split(input, ",")
	if len(paths) == 0 {
		fmt.Println("Не указаны файлы для загрузки.")
		return
	}
	var images []*pb.ImageData
	for _, path := range paths {
		path = strings.TrimSpace(path)
		data, err := os.ReadFile(path)
		if err != nil {
			fmt.Println("Ошибка чтения файла:", path, err)
			continue
		}
		image := &pb.ImageData{
			Filename: path,
			Data:     data,
		}
		images = append(images, image)
	}
	resp, err := client.UploadImage(context.Background(), &pb.ImageBatch{Images: images})
	if err != nil {
		fmt.Println("Ошибка при загрузке файлов:", err)
		return
	}
	fmt.Println("Файлы загружены. Ответ сервера:", resp.Info)
}

func listImages(client pb.ImageServiceClient) {
	fmt.Print("Функция listImages")
}

func downloadImage(client pb.ImageServiceClient) {
	fmt.Print("Функция downloadImage ")
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
