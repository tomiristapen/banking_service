package utils

import (
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtKey = []byte(os.Getenv("JWT_SECRET"))

func GenerateJWT(userID string, email string) (string, error) {
    claims := jwt.MapClaims{
        "user_id": userID,
        "email":   email,
        "exp":     time.Now().Add(24 * time.Hour).Unix(),
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtKey)
}

func ParseToken(tokenStr string) (string, error) {
    token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })

    if err != nil || !token.Valid {
        return "", errors.New("invalid token")
    }

    claims, ok := token.Claims.(jwt.MapClaims)
    if !ok {
        return "", errors.New("invalid claims")
    }

    userID, ok := claims["user_id"].(string)
    if !ok {
        return "", errors.New("user_id not found in token")
    }

    return userID, nil
}
