package jwt

import (
	"testing"
	"time"
)

func TestPasswordHashAndCompare(t *testing.T) {
	hash, err := HashPassword("super-secret-123")
	if err != nil {
		t.Fatalf("HashPassword: %v", err)
	}
	if err := ComparePassword("super-secret-123", hash); err != nil {
		t.Fatalf("ComparePassword (correct): %v", err)
	}
	if err := ComparePassword("wrong", hash); err == nil {
		t.Fatalf("ComparePassword (wrong) should have errored")
	}
}

func TestGenerateAndValidateToken(t *testing.T) {
	secret := "test-secret"
	tok, err := GenerateToken(secret, 1, 42, "admin", "admin")
	if err != nil {
		t.Fatalf("GenerateToken: %v", err)
	}
	if tok == "" {
		t.Fatal("token is empty")
	}

	claims, err := ValidateToken(secret, tok)
	if err != nil {
		t.Fatalf("ValidateToken: %v", err)
	}
	if claims.UserID != 42 || claims.Login != "admin" || claims.Role != "admin" {
		t.Fatalf("claims wrong: %+v", claims)
	}
}

func TestValidateToken_BadSecret(t *testing.T) {
	tok, _ := GenerateToken("right", 1, 1, "u", "admin")
	if _, err := ValidateToken("wrong", tok); err == nil {
		t.Fatal("ValidateToken with wrong secret should have errored")
	}
}

func TestValidateToken_EmptyInputs(t *testing.T) {
	if _, err := ValidateToken("", "x"); err == nil {
		t.Fatal("empty secret should error")
	}
	if _, err := ValidateToken("k", ""); err == nil {
		t.Fatal("empty token should error")
	}
}

func TestGenerateToken_BadInputs(t *testing.T) {
	if _, err := GenerateToken("", 1, 1, "u", "admin"); err == nil {
		t.Fatal("empty secret should error")
	}
	if _, err := GenerateToken("k", 0, 1, "u", "admin"); err == nil {
		t.Fatal("non-positive expiry should error")
	}
}

func TestTokenExpires(t *testing.T) {
	tok, _ := GenerateToken("k", 1, 1, "u", "admin")

	claims, err := ValidateToken("k", tok)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if claims.ExpiresAt == nil {
		t.Fatal("no expiry set")
	}
	if time.Until(claims.ExpiresAt.Time) <= 0 {
		t.Fatal("token already expired")
	}
}
