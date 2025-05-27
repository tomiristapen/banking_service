package utils

import (
    "os"
    "testing"
)

func TestGenerateAndParseJWT(t *testing.T) {
    // Set test secret manually
    os.Setenv("JWT_SECRET", "test_secret_key")

    userID := "123456"
    email := "test@example.com"

    // Generate token
    token, err := GenerateJWT(userID, email)
    if err != nil {
        t.Fatalf("failed to generate JWT: %v", err)
    }

    if token == "" {
        t.Fatal("expected token, got empty string")
    }

    // Parse token
    parsedID, err := ParseToken(token)
    if err != nil {
        t.Fatalf("failed to parse JWT: %v", err)
    }

    if parsedID != userID {
        t.Errorf("expected userID %s, got %s", userID, parsedID)
    }
}
