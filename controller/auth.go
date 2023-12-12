package controller

import (
	"database/sql"
	"fmt"
	"log"
	"project-unit-1/entities"
)

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

func ReadAccount(db *sql.DB, phoneNumber, password string) {
	var usersWithBalance []struct {
		User    entities.User
		Balance entities.Balance
	}

	// menjalankan perintah query select dengan operasi JOIN
	rows, errSelect := db.Query("SELECT u.id AS user_id, u.name, u.phone_number, u.alamat, u.email, b.amount, b.balance_at FROM users u LEFT JOIN balances b ON u.id = b.user_id WHERE u.phone_number = ? AND u.password = ?", phoneNumber, password)
	if errSelect != nil {
		log.Fatal("error run query select ", errSelect.Error())
	}

	for rows.Next() {
		var userWithBalance struct {
			User    entities.User
			Balance entities.Balance
		}
		errScan := rows.Scan(&userWithBalance.User.ID, &userWithBalance.User.Name, &userWithBalance.User.Phone_number, &userWithBalance.User.Alamat, &userWithBalance.User.Email, &userWithBalance.Balance.Amount, &userWithBalance.Balance.UpdatedAt)
		if errScan != nil {
			log.Fatal("error scan select", errScan.Error())
		}
		usersWithBalance = append(usersWithBalance, userWithBalance)
	}

	for _, u := range usersWithBalance {
		fmt.Printf("Nama: %v, Email: %v, Alamat: %v, Jumlah Saldo: %v, Updated At: %v\n", u.User.Name, u.User.Email, u.User.Alamat, u.Balance.Amount, u.Balance.UpdatedAt)
	}
}
