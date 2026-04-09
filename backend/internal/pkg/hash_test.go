package pkg

import "testing"

func TestHashPassword(t *testing.T) {
	hash, err := HashPassword("test1234")
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}
	if hash == "" {
		t.Fatal("hash should not be empty")
	}
	if hash == "test1234" {
		t.Fatal("hash should not equal plaintext")
	}
}

func TestCheckPassword(t *testing.T) {
	password := "mysecretpass"
	hash, _ := HashPassword(password)

	if !CheckPassword(password, hash) {
		t.Error("CheckPassword should return true for correct password")
	}
	if CheckPassword("wrongpass", hash) {
		t.Error("CheckPassword should return false for wrong password")
	}
	if CheckPassword("", hash) {
		t.Error("CheckPassword should return false for empty password")
	}
}

func TestHashPasswordDifferentHashes(t *testing.T) {
	hash1, _ := HashPassword("same")
	hash2, _ := HashPassword("same")
	if hash1 == hash2 {
		t.Error("same password should produce different hashes (bcrypt salt)")
	}
}
