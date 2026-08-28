package main

import (
	"TMS/notification-service/internal/database"
	"TMS/notification-service/internal/messaging"
	"TMS/notification-service/internal/repository"
	"TMS/notification-service/internal/routers"
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {

	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load(".env")
	if err != nil {
		slog.Info(".env file not found,:environment variables will be used")
	}
	if os.Getenv("GIN_MODE") != "" {
		gin.SetMode(os.Getenv("GIN_MODE"))
	}

	db, err := database.Connect()

	if err != nil {
		slog.Error("Server error", "error", err.Error())
		return
	}

	rabbit, err := messaging.NewRabbitMQ(os.Getenv("RABBITMQ_URL"))

	if err != nil {
		slog.Error("RabbitMQ connection failed", "error", err)
		return
	}
	repo := repository.NewRepository(db)
	r := routers.SetupRouter(repo, rabbit)

	defer func() {
		_ = db.Close()
	}()
	//----------------------------------------
	// burada mq kapanirsa giderse yeniden acilacak
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done() //bu goroutine kapandigidna wg.done diyecek ve tetiklenecek
		for {

			if err := rabbit.Consume(ctx, repo); err != nil { //consume bittiginde direkt  conn.close yapamz // contexi iceri gonderdik daha sonra consumeyi kapatmak icin ctx den bildirecegiz
				slog.Error("RabbitMQ consumer stopped", "error", err) //burada baglanti gitti
			}
			if ctx.Err() != nil {
				break //ctx ten cancel ile shutdown tetiklendiyse yeni mq acma
			}
			_ = rabbit.Conn.Close()
			//sonra baglanti gittigi icin yeni bir mq aciyoruz
			select {
			case <-time.After(5 * time.Second): //burada shutdown komutu  gelirse beklemeyelim diye done ile beraeber koyduk
				//ilk veri kimden gelirse onu sececek
				// yeni mq olusturuyoruz
				newRabbit, err := messaging.NewRabbitMQ(os.Getenv("RABBITMQ_URL"))
				if err != nil {
					slog.Error("RabbitMQ reconnect failed", "error", err)
					continue
				}
				rabbit = newRabbit //yeni mq yu eskisi ile degistiryoruz
			case <-ctx.Done():
				return
			}

		}
	}()
	//----------------------------------------
	port := os.Getenv("PORT") //portumuzu alıyoruz

	if port == "" {
		port = "8081"
	}
	slog.Info("Server starting", "port", port)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: r,
	}
	go func() {
		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			slog.Error("Server error", "error", err)
		}
	}()
	quit := make(chan os.Signal, 1) //ayni channel mekanizmasi burada

	signal.Notify(
		quit,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	<-quit
	cancel()
	shutdownCtx, shutdownCancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer shutdownCancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		slog.Error("Server shutdown failed", "error", err)
	}
	wg.Wait() //consumerin kapanmasini bekliyoruz defer ile koymustuk

	if rabbit != nil {
		_ = rabbit.Conn.Close()

	}
}
