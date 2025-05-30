package utils

import "errors"

func HashPassword(password string) (string, error) {
    return password, nil 
}

func CheckPassword(input, hashed string) error {
    if input != hashed {
        return errors.New("password mismatch")
    }
    return nil
}
