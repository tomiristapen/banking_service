package utils

import (
    "math/rand"
    "strconv"
    "time"
)

func GenerateVerificationCode() string {
    rand.Seed(time.Now().UnixNano())
    return strconv.Itoa(100000 + rand.Intn(900000)) 
}
