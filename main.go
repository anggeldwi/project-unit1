package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type AppConfig struct {
	DB_USERNAME string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func InitDB() (*sql.DB, error) {
	// DSN (Data Source Name):
	// <username>:<password>@tcp(<hostname>:<port-db>)/<db-name>
	var cfg = AppConfig{
		DB_USERNAME: os.Getenv("DB_USERNAME"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}
	var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	var db *sql.DB
	var err error
	// Get a database handle.
	db, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatal("error open connection to db: ", err)
	}

	// db.SetConnMaxLifetime(time.Minute * 3)
	// db.SetMaxOpenConns(10)
	// db.SetMaxIdleConns(10)

	// cek apakah sudah bisa connect ke db
	pingErr := db.Ping()
	if pingErr != nil {
		log.Println("err ping conenction: ", pingErr)
		return nil, pingErr
	}
	fmt.Println("success connect to db!")
	return db, nil
}
func main() {
	db, errInitDB := InitDB()
	if errInitDB != nil {
		log.Fatal("error connect to db", errInitDB)
	}
	//close connection
	defer db.Close()

	// membuat menu
	fmt.Println("Pilih menu:")
	fmt.Println("[1]: Read data")
	fmt.Println("[2]: Add data")
	fmt.Println("[3]: Delete data")
	fmt.Println("[4]: Update data")
	fmt.Println("[0]: keluar aplikasi")
	var pilihan int
	fmt.Println("Masukkan angka sesuai pilihan menu:")
	fmt.Scanln(&pilihan)

	//membuat switch case
	switch pilihan {
	case 1:
		fmt.Println("menu1.")

	case 2:
		fmt.Println("menu2.")
	}
}
