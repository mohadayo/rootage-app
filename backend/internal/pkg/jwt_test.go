package pkg

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const testSecret = "test-secret-key"

func TestGenerateAndValidateToken(t *testing.T) {
	token, err := GenerateToken("user-123", "test@example.com", "user", testSecret)
	if err != nil {
		t.Fatalf("GenerateToken failed: %v", err)
	}
	if token == "" {
		t.Fatal("token should not be empty")
	}

	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.UserID != "user-123" {
		t.Errorf("UserID = %q, want %q", claims.UserID, "user-123")
	}
	if claims.Email != "test@example.com" {
		t.Errorf("Email = %q, want %q", claims.Email, "test@example.com")
	}
	if claims.Role != "user" {
		t.Errorf("Role = %q, want %q", claims.Role, "user")
	}
}

func TestValidateToken_WrongSecret(t *testing.T) {
	token, _ := GenerateToken("user-123", "test@example.com", "user", testSecret)

	_, err := ValidateToken(token, "wrong-secret")
	if err == nil {
		t.Error("ValidateToken should fail with wrong secret")
	}
}

func TestValidateToken_InvalidToken(t *testing.T) {
	_, err := ValidateToken("not.a.valid.token", testSecret)
	if err == nil {
		t.Error("ValidateToken should fail with invalid token string")
	}
}

func TestValidateToken_EmptyToken(t *testing.T) {
	_, err := ValidateToken("", testSecret)
	if err == nil {
		t.Error("ValidateToken should fail with empty token")
	}
}

func TestValidateToken_ExpiredToken(t *testing.T) {
	claims := Claims{
		UserID: "user-123",
		Email:  "test@example.com",
		Role:   "user",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, _ := token.SignedString([]byte(testSecret))

	_, err := ValidateToken(tokenString, testSecret)
	if err == nil {
		t.Error("ValidateToken should fail with expired token")
	}
}

func TestGenerateToken_AdminRole(t *testing.T) {
	token, _ := GenerateToken("admin-1", "admin@example.com", "admin", testSecret)
	claims, err := ValidateToken(token, testSecret)
	if err != nil {
		t.Fatalf("ValidateToken failed: %v", err)
	}
	if claims.Role != "admin" {
		t.Errorf("Role = %q, want %q", claims.Role, "admin")
	}
}
