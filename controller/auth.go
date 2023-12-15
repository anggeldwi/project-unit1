package controller

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"project-unit-1/entities"
	"strings"
)

// Fitur login
func Login(db *sql.DB, phoneNumber, password string) (*entities.User, error) {
	var user entities.User

	// Jalankan query untuk mencari user berdasarkan nomor telepon dan password
	row := db.QueryRow("SELECT id, name, phone_number, alamat, email, password, created_at FROM users WHERE phone_number = ? AND password = ?", phoneNumber, password)

	// Scan hasil query ke dalam variabel user
	err := row.Scan(&user.ID, &user.Name, &user.Phone_number, &user.Alamat, &user.Email, &user.Password, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("login failed: user not found")
		}
		log.Fatal("error during login query:", err)
		return nil, fmt.Errorf("login failed: unexpected error")
	}

	return &user, nil
}

// Fitur read account
func ReadAccount(db *sql.DB, phoneNumber string) {
	var usersWithBalance []struct {
		User    entities.User
		Balance entities.Balance
	}

	// menjalankan perintah query select dengan operasi JOIN
	rows, errSelect := db.Query("SELECT u.id AS user_id, u.name, u.phone_number, u.alamat, u.email, u.created_at, b.amount, b.balance_at FROM users u LEFT JOIN balances b ON u.id = b.user_id WHERE u.phone_number = ?", phoneNumber)
	if errSelect != nil {
		log.Fatal("error run query select ", errSelect.Error())
	}

	for rows.Next() {
		var userWithBalance struct {
			User    entities.User
			Balance entities.Balance
		}
		errScan := rows.Scan(&userWithBalance.User.ID, &userWithBalance.User.Name, &userWithBalance.User.Phone_number, &userWithBalance.User.Alamat, &userWithBalance.User.Email, &userWithBalance.User.CreatedAt, &userWithBalance.Balance.Amount, &userWithBalance.Balance.UpdatedAt)
		if errScan != nil {
			log.Fatal("error scan select", errScan.Error())
		}
		usersWithBalance = append(usersWithBalance, userWithBalance)
	}

	for _, u := range usersWithBalance {

		// Format waktu pada Created At
		createdAt := u.User.CreatedAt.Format("2006-01-02 15:04:05")
		// Format waktu pada Updated At
		updatedAt := u.Balance.UpdatedAt.Format("2006-01-02 15:04:05")

		fmt.Printf("Nama: %v\nEmail: %v\nAlamat: %v\nNo.Phone: %v\nJumlah Saldo:%v\nCreated At:%v\nUpdated At: %v\n", u.User.Name, u.User.Email, u.User.Alamat, u.User.Phone_number, u.Balance.Amount, createdAt, updatedAt)
	}
}

// Fitur update account
func UpdateAccount(db *sql.DB, user *entities.User) error {
	// Menampilkan informasi akun sebelum diperbarui
	fmt.Printf("Informasi Akun Sebelum Diperbarui:\n")
	fmt.Printf("Nama: %s, Nomor Telepon: %s, Alamat: %s, Email: %s\n", user.Name, user.Phone_number, user.Alamat, user.Email)

	// Menerima input dari pengguna untuk pembaruan
	var newName, newPhoneNumber, newAlamat, newEmail string

	// Membersihkan newline yang mungkin masih ada di dalam buffer input
	fmt.Scanln()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Masukkan nama baru: ")
	newName, _ = reader.ReadString('\n')
	newName = strings.TrimSpace(newName)

	fmt.Print("Masukkan nomor telepon baru: ")
	newPhoneNumber, _ = reader.ReadString('\n')
	newPhoneNumber = strings.TrimSpace(newPhoneNumber)

	fmt.Print("Masukkan alamat baru: ")
	newAlamat, _ = reader.ReadString('\n')
	newAlamat = strings.TrimSpace(newAlamat)

	fmt.Print("Masukkan email baru: ")
	newEmail, _ = reader.ReadString('\n')
	newEmail = strings.TrimSpace(newEmail)

	// Update informasi akun di database
	_, err := db.Exec("UPDATE users SET name=?, phone_number=?, alamat=?, email=? WHERE id=?", newName, newPhoneNumber, newAlamat, newEmail, user.ID)
	if err != nil {
		log.Fatal("error during account update:", err)
		return fmt.Errorf("failed to update account information")
	}

	// Menampilkan informasi akun setelah diperbarui
	fmt.Printf("\nInformasi Akun Setelah Diperbarui:\n")
	fmt.Printf("Nama:           %s\n", newName)
	fmt.Printf("Nomor Telepon:  %s\n", newPhoneNumber)
	fmt.Printf("Alamat:         %s\n", newAlamat)
	fmt.Printf("Email:          %s\n", newEmail)

	return nil
}

// Fitur top up
func TopUp(db *sql.DB, user *entities.User) error {

	// Menerima input dari pengguna untuk top up
	var amount float64
	fmt.Print("Masukkan jumlah saldo yang ingin ditambahkan: ")
	_, err := fmt.Scan(&amount)
	if err != nil {
		return fmt.Errorf("error reading top-up amount: %w", err)
	}

	// Membersihkan newline yang mungkin masih ada di dalam buffer input
	fmt.Scanln()

	// Memastikan jumlah top up positif
	if amount <= 0 {
		fmt.Println("Jumlah top up harus lebih dari 0.")
		return nil
	}

	// Update saldo di database
	_, err = db.Exec("UPDATE balances SET amount = amount + ?, balance_at = NOW() WHERE user_id = ?", amount, user.ID)
	if err != nil {
		return fmt.Errorf("error during balance update: %w", err)
	}

	// Menyimpan riwayat top-up ke dalam tabel top_ups_history
	_, err = db.Exec("INSERT INTO top_ups_history (user_id, amount, top_up_at) VALUES (?, ?, NOW())", user.ID, amount)
	if err != nil {
		return fmt.Errorf("error saving top-up history: %w", err)
	}

	// Menampilkan informasi saldo setelah top up
	newBalance, err := getBalance(db, user.ID)
	if err != nil {
		return fmt.Errorf("error getting new balance: %w", err)
	}
	user.Balance.Amount = newBalance
	fmt.Printf("\nSaldo Setelah Top Up: %.2f\n", user.Balance.Amount)

	return nil
}

// getBalance mengambil balance dari database berdasarkan user ID
func getBalance(db *sql.DB, userID uint) (float64, error) {
	var balance float64
	err := db.QueryRow("SELECT amount FROM balances WHERE user_id = ?", userID).Scan(&balance)
	if err != nil {
		return 0, err
	}
	return balance, nil
}
