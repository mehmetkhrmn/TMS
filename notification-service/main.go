package main

import (
	"TMS/notification-service/internal/database"
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
	r := routers.SetupRouter(db)
	defer db.Close()
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
