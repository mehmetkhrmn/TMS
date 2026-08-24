package database

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

func Connect() (*sql.DB, error) {

	connStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, err := sql.Open("postgres", connStr) //iki değişken döndürüyor ve sıralamsı db,err şeklinde
	if err != nil {
		return nil, err
	}
	for attempt := 1; attempt <= 5; attempt++ {
		err = db.Ping()
		if err == nil { //hata yoksa db dondurur
			return db, nil

		}
		if attempt < 5 {
			delay := time.Duration(attempt*3) * time.Second //linear bekleme suresi artisi
			log.Println("Database unavailable, retrying...")
			time.Sleep(delay)
		}
	}
	return nil, err
}
