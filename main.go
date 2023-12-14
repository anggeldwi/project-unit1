package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"project-unit-1/controller"
	"project-unit-1/entities"

	_ "github.com/go-sql-driver/mysql"
)

// type AppConfig struct {
// 	DB_USERNAME string
// 	DB_PASSWORD string
// 	DB_HOST     string
// 	DB_PORT     string
// 	DB_NAME     string
// }

func InitDB() (*sql.DB, error) {
	// DSN (Data Source Name):
	// <username>:<password>@tcp(<hostname>:<port-db>)/<db-name>
	// var cfg = AppConfig{
	// 	DB_USERNAME: os.Getenv("DB_USERNAME"),
	// 	DB_PASSWORD: os.Getenv("DB_PASSWORD"),
	// 	DB_HOST:     os.Getenv("DB_HOST"),
	// 	DB_PORT:     os.Getenv("DB_PORT"),
	// 	DB_NAME:     os.Getenv("DB_NAME"),
	// }
	// var connectionString = fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", cfg.DB_USERNAME, cfg.DB_PASSWORD, cfg.DB_HOST, cfg.DB_PORT, cfg.DB_NAME)
	var connectionString = os.Getenv("CONNECTION_DB")
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
	// close connection
	defer db.Close()

	var user *entities.User = entities.NewUser() // variabel untuk menyimpan informasi user yang berhasil login

	for {
		// membuat menu
		fmt.Println("Pilih menu:")
		fmt.Println("[1]: Login")
		fmt.Println("[2]: Lihat profil")
		fmt.Println("[3]: Update Profil")
		fmt.Println("[4]: Top Up")
		fmt.Println("[5]: Riwayat Top Up")
		fmt.Println("[0]: Keluar aplikasi")
		var pilihan int
		fmt.Println("Masukkan angka sesuai pilihan menu:")
		_, err := fmt.Scan(&pilihan)
		if err != nil {
			fmt.Println("Error reading input:", err)
			continue
		}
		// membuat switch case
		switch pilihan {
		case 1:
			fmt.Println("Login:")

			// Membersihkan newline yang mungkin masih ada di dalam buffer input
			fmt.Scanln()

			// Tambahkan dua variabel untuk menyimpan nomor telepon dan kata sandi
			var nomorTelepon, password string

			fmt.Print("Masukkan nomor telepon: ")
			fmt.Scanln(&nomorTelepon)

			fmt.Print("Masukkan kata sandi: ")
			fmt.Scanln(&password)

			// memanggil fungsi login
			loggedInUser, err := controller.Login(db, nomorTelepon, password)
			if err != nil {
				fmt.Println("Login failed:", err)
				continue
			}
			fmt.Printf("Login successful! Welcome, %s\n", loggedInUser.Name)
			// menyimpan informasi user yang berhasil login
			user = loggedInUser
		case 2:
			// Memastikan bahwa user sudah login sebelum mengakses menu ini
			if user == nil || user.Name == "" {
				fmt.Println("Silakan login terlebih dahulu.")
				continue
			}
			fmt.Println("Lihat Profil:")
			// memanggil fungsi ReadAccount dengan user yang sudah login
			controller.ReadAccount(db, user.Phone_number, user.Password)

		case 3:
			// Memastikan bahwa user sudah login sebelum mengakses menu ini
			if user == nil || user.Name == "" {
				fmt.Println("Silakan login terlebih dahulu.")
				continue
			}
			fmt.Println("Update Data Akun:")
			// memanggil fungsi UpdateAccount dengan user yang sudah login
			err := controller.UpdateAccount(db, user)
			if err != nil {
				fmt.Println("Gagal memperbarui informasi akun:", err)
			}

		case 4:
			// Memastikan bahwa user sudah login sebelum mengakses menu ini
			if user == nil || user.Name == "" {
				fmt.Println("Silakan login terlebih dahulu.")
				continue
			}
			fmt.Println("Top Up Saldo:")
			// memanggil fungsi TopUp dengan user yang sudah login
			err := controller.TopUp(db, user)
			if err != nil {
				fmt.Println("Gagal melakukan top up saldo:", err)
			}

		case 5:
			// Memastikan bahwa user sudah login sebelum mengakses menu ini
			if user == nil || user.Name == "" {
				fmt.Println("Silakan login terlebih dahulu.")
				continue
			}
			fmt.Println("Lihat Riwayat Top Up:")
			// memanggil fungsi ViewTopUpHistory dengan ID pengguna yang sudah login
			controller.ViewTopUpHistory(db, user.ID)

		case 0:
			fmt.Println("Sukses keluar dari aplikasi.")
			return
		}
	}
}
