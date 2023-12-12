package entities

import "time"

type User struct {
	ID           uint
	Name         string
	Phone_number string
	Alamat       string
	Email        string
	Password     string
	CreatedAt    time.Time
}

type Balance struct {
	ID        int
	UserID    int
	Amount    float64
	UpdatedAt string
}
