package main

import (
	"TMS/notification-service/internal/database"
	"TMS/notification-service/internal/messaging"
	"TMS/notification-service/internal/repository"
	"TMS/notification-service/internal/routers"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)
	gin.SetMode(gin.DebugMode)

	db, err := database.Connect()
	if err != nil {
		slog.Error("Sunucu hatası", "error", err.Error())
	}
	repo := repository.NewRepository(db)
	r := routers.SetupRouter(repo)

	defer func() {
		_ = db.Close()
	}()
	rabbit, err := messaging.NewRabbitMQ(os.Getenv("RABBITMQ_URL"))
	if err != nil {
		slog.Error("RabbitMQ connection failed", "error", err)
		os.Exit(1)
	}
	defer func() {
		_ = rabbit.Conn.Close()
	}()

	go func() {
		if err := rabbit.Consume(repo); err != nil {
			slog.Error("RabbitMQ consumer stopped", "error", err)
		}
	}()

	port := os.Getenv("PORT") //portumuzu alıyoruz

	if port == "" {
		port = "8081"
	}
	slog.Info("Sunucu başlatılıyor", "port", port)
	if err := r.Run(":" + port); err != nil { //burada route başlatıyoruz hata varsa err assinglıyor
		slog.Error("Sunucu hatası", "error", err.Error())
		os.Exit(1)
	}

}
