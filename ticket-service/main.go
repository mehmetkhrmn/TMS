package main

import (
	"TMS/ticket-service/internal/database"
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/routers"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) //bu consolda çıktıları JSON formatında almak için
	slog.SetDefault(logger)                                 //varsayılanı değiştirdik
	gin.SetMode(gin.DebugMode)
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	required := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"JWT_SECRET",
		"NOTIFICATION_SERVICE_URL",
		"RABBITMQ_URL",
		"RABBITMQ_QUEUE",
		"RABBITMQ_ROUTING_KEY",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Required environment variable %s is not set", key)
		}
	}
	db, err := database.Connect()
	if err != nil {
		slog.Error("Sunucu hatası", "error", err.Error())
		return
	}

	defer func() {
		err := db.Close()
		if err != nil {
			log.Printf("Error closing DB connection")
		}
	}()
	rabbit, err := messaging.NewRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err := rabbit.Conn.Close()
		if err != nil {
			log.Printf("Error closing RabbitMQ connection")
		}
	}()

	defer func() {
		err := rabbit.Channel.Close()
		if err != nil {
			log.Printf("Error closing RabbitMQ channel")
		}
	}()
	if err := rabbit.QueueDeclare(); err != nil {
		log.Fatal(err)
	}

	r := routers.SetupRouter(db, rabbit)
	port := os.Getenv("PORT") //portumuzu alıyoruz

	if port == "" {
		log.Fatal("PORT is not set")
	}
	slog.Info("Sunucu başlatılıyor", "port", port)
	if err := r.Run(":" + port); err != nil { //burada route başlatıyoruz hata varsa err assinglıyor
		slog.Error("Sunucu hatası", "error", err.Error())
		os.Exit(1)
	}

}
