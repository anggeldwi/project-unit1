package entities

import "time"

func NewUser() *User {
	return &User{
		Balance: Balance{},
	}
}

type User struct {
	ID           uint
	Name         string
	Phone_number string
	Alamat       string
	Email        string
	Password     string
	CreatedAt    time.Time
	Balance      Balance
}

type Balance struct {
	ID        int
	UserID    int
	Amount    float64
	UpdatedAt time.Time
}

type TopUpHistory struct {
	ID      int
	UserID  int
	Amount  float64
	TopUpAt time.Time
}
