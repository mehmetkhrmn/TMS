package database

import (
	"database/sql"
	"log"
	"os"
	"time"
)

func Connect() (*sql.DB, error) {

	connStr := "host=" + os.Getenv("DB_HOST") +
		" port=" + os.Getenv("DB_PORT") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	db, dberr := sql.Open("postgres", connStr) //iki değişken döndürüyor ve sıralamsı db,err şeklinde

	if dberr != nil {
		return nil, dberr

	}
	for { //bağlantıdan olumlu dönüş alana kadar tekrarlayan mekanizma
		if err := db.Ping(); err == nil {
			break
		}

		log.Println("Database unavailable, retrying...")
		time.Sleep(5 * time.Second)
	}
	return db, nil
}
