package service

import "testing"

const TOKEN = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEifQ.1LUx7XnbjaKzc9eaKZc9PsZWyzm_VuMYbaOIU2J7uQE"

func TestValidateTokenSuccess(t *testing.T) {
	t.Setenv("JWT_KEY", "testing")
	err := ValidateToken(TOKEN)
	if err != nil {
		t.Fatalf("Expected no error")
	}
}

func TestValidateTokenError(t *testing.T) {
	err := ValidateToken(TOKEN)
	if err == nil {
		t.Fatalf("Expected error")
	}
}
