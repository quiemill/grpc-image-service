package main

import (
	"bufio"
	"fmt"
	pb "grpc-image-service/api/gen/image_service"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
)

func uploadImage(client pb.ImageServiceClient) {
	fmt.Println("Функция uploadImage")
}

func listImages(client pb.ImageServiceClient) {
	fmt.Println("Функция listImages")
}

func downloadImage(client pb.ImageServiceClient) {
	fmt.Println("Функция downloadImage")
}

func main() {
	// Подключаемся к серверу
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()

	client := pb.NewImageServiceClient(conn)
	for {
		fmt.Println("\nВыберите действие:")
		fmt.Println("Загрузить изображения")
		fmt.Println("Посмотреть список изображений")
		fmt.Println("Скачать изображения")
		fmt.Println("Выйти")

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
			fmt.Println("🚪 Выход...")
			return
		default:
			fmt.Println("Попробуйте снова.")
		}
	}
}
