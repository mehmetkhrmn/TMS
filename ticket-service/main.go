package main

import (
	"TMS/ticket-service/internal/database"
	"TMS/ticket-service/internal/messaging"
	"TMS/ticket-service/internal/routers"
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) //bu consolda çıktıları JSON formatında almak için
	slog.SetDefault(logger)                                 //varsayılanı değiştirdik
	err := godotenv.Load(".env")
	if err != nil {
		slog.Info(".env file not found,:environment variables will be used")
	}
	required := []string{
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASSWORD",
		"DB_NAME",
		"DB_SSLMODE",
		"JWT_SECRET",
		"RABBITMQ_URL",
		"RABBITMQ_QUEUE",
		"RABBITMQ_ROUTING_KEY",
	}

	for _, key := range required {
		if os.Getenv(key) == "" {
			log.Fatalf("Required environment variable %s is not set", key)
		}
	}
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}
	db, err := database.Connect()
	if err != nil {
		slog.Error("Server error", "error", err.Error())
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

	r := routers.SetupRouter(db, rabbit)
	port := os.Getenv("PORT") //portumuzu alıyoruz

	if port == "" {
		log.Fatal("PORT is not set")
	}
	slog.Info("Server starting", "port", port)

	server := &http.Server{ //serveri kendimiz olusturduk
		Addr:    ":" + port,
		Handler: r,
	}

	go func() { //serveri ayri go routine de baslattik
		if err := server.ListenAndServe(); err != nil && //istekleri dinliyoruz
			err != http.ErrServerClosed {
			slog.Error("Server error", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(
		quit,
		syscall.SIGINT,  //ctrl+c
		syscall.SIGTERM, //quit
	)

	<-quit //sinyal gelene kadar main burada bekliyor ve kapanmiyor

	ctx, cancel := context.WithTimeout( //zaman siniri koyuyor
		context.Background(), //bos bir context
		10*time.Second,       //10 saniye yasiyor
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil { //buradada ctx in olmesini bekliyor olunce shutdown atacak
		slog.Error("Server shutdown failed", "error", err)
	}
}
