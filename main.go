package main

import (
	"TMS/internal/routers"
	"database/sql"
	"log"
	"log/slog"
	"os"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil)) //bu consolda çıktıları JSON formatında almak için
	slog.SetDefault(logger)                                 //varsayılanı değiştirdik
	gin.SetMode(gin.DebugMode)
	connStr := "host=localhost port=5432 user=postgres password=g42v24bh dbname=postgres sslmode=disable"

	db, dberr := sql.Open("postgres", connStr) //iki değişken döndürüyor ve sıralamsı db,err şeklinde
	defer db.Close()                           //bu main bitince çalışıyor

	if dberr != nil {
		log.Fatal(dberr)
	}
	if dberr := db.Ping(); dberr != nil { //buradada db ye pingliyoruz hata döndürüyorsa yine logluyoruz
		log.Fatal(dberr)
	}

	r := routers.SetupRouter(db)

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
