package service

import "testing"

func TestValidatePassword(t *testing.T) {
	cases := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{"too short", "1234567", true},
		{"exactly min", "abcd1234", false},
		{"common weak password", "password", true},
		{"common weak numeric", "12345678", true},
		{"strong", "Str0ng-Passphrase", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := validatePassword(tc.password)
			if (err != nil) != tc.wantErr {
				t.Errorf("validatePassword(%q) err = %v, wantErr = %v", tc.password, err, tc.wantErr)
			}
		})
	}
}

func TestNormalizeEmail(t *testing.T) {
	cases := map[string]string{
		"  Taro.Yamada@Example.com ": "taro.yamada@example.com",
		"USER@TEST.EXAMPLE.COM":       "user@test.example.com",
		"already@lower.com":           "already@lower.com",
	}
	for in, want := range cases {
		if got := normalizeEmail(in); got != want {
			t.Errorf("normalizeEmail(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestIsValidDifficulty(t *testing.T) {
	valid := []string{"beginner", "intermediate", "advanced"}
	for _, d := range valid {
		if !isValidDifficulty(d) {
			t.Errorf("isValidDifficulty(%q) = false, want true", d)
		}
	}
	for _, d := range []string{"中級", "easy", "", "BEGINNER"} {
		if isValidDifficulty(d) {
			t.Errorf("isValidDifficulty(%q) = true, want false", d)
		}
	}
}
