package model

import "time"

type User struct {
    ID         string    
    Name       string
    Email      string
    Password   string    
    IsVerified bool
    CreatedAt  time.Time
    Balance    float64
}
