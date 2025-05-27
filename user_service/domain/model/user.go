package model

import "time"

// User — чистая бизнес-сущность без зависимостей от Mongo, JSON и т.п.
type User struct {
    ID         string    // UUID или ObjectID (но без привязки к Mongo!)
    Name       string
    Email      string
    Password   string    // Хешированный пароль
    IsVerified bool
    CreatedAt  time.Time
    Balance    float64
}
