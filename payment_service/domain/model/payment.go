package model

import "time"

type Payment struct {
	ID        string
	UserID    string
	Amount    float64
	Service   string
	Category  string
	Type      string    // "payment" — фиксировано для этого сервиса
	Status    string
	CreatedAt time.Time
}
