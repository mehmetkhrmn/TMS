package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
		os.Exit(1)
	}
	connStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, dberr := sql.Open("postgres", connStr) //iki değişken döndürüyor ve sıralamsı db,err şeklinde
	//bu main bitince çalışıyor

	if dberr != nil {
		log.Fatal(dberr)
	}
	if dberr := db.Ping(); dberr != nil { //buradada db ye pingliyoruz hata döndürüyorsa yine logluyoruz
		log.Fatal(dberr)
	}
	return db, nil
}
