package main

import (
	"TMS/ticket-service/internal/database"
	"TMS/ticket-service/internal/routers"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) //bu consolda çıktıları JSON formatında almak için
	slog.SetDefault(logger)                                 //varsayılanı değiştirdik
	gin.SetMode(gin.DebugMode)

	db, err := database.Connect()
	if err != nil {
		slog.Error("Sunucu hatası", "error", err.Error())
	}
	r := routers.SetupRouter(db)
	defer db.Close()
	port := os.Getenv("PORT") //portumuzu alıyoruz

	if port == "" {
		port = "8080"
	}
	slog.Info("Sunucu başlatılıyor", "port", port)
	if err := r.Run(":" + port); err != nil { //burada route başlatıyoruz hata varsa err assinglıyor
		slog.Error("Sunucu hatası", "error", err.Error())
		os.Exit(1)
	}

}
