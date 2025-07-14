// package auth
// import (
// 	"time"
// 	"github.com/google/uuid"
// 	"testing"
// )

// func TestMakeAndValidateJWT(t *testing.T) {
// 	userID := uuid.New()
// 	secret := "test-secret"
// 	validDuration := time.Hour

// 	token, err := MakeJWT(userID, secret, validDuration)
// 	if err != nil {
// 		t.Fatalf("MakeJWT failed: %v", err)
// 	}

// 	// Test valid token
// 	parsedID, err := ValidateJWT(token, secret)
// 	if err != nil {
// 		t.Fatalf("ValidateJWT failed: %v", err)
// 	}
// 	if parsedID != userID {
// 		t.Fatalf("Expected user ID %v, got %v", userID, parsedID)
// 	}

// 	// Test expired token
// 	expiredToken, _ := MakeJWT(userID, secret, -time.Hour)
// 	_, err = ValidateJWT(expiredToken, secret)
// 	if err == nil {
// 		t.Fatal("Expected error for expired token")
// 	}

// 	// Test wrong secret
// 	_, err = ValidateJWT(token, "wrong-secret")
// 	if err == nil {
// 		t.Fatal("Expected error for wrong secret")
// 	}
// }




// auth_test.go
package auth

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	testCases := []struct {
		name      string
		userID    uuid.UUID
		secret    string
		expiresIn time.Duration
		wantErr   bool
	}{
		{
			name:      "Valid token",
			userID:    uuid.New(),
			secret:    "test-secret",
			expiresIn: time.Hour,
			wantErr:   false,
		},
		{
			name:      "Expired token",
			userID:    uuid.New(),
			secret:    "test-secret",
			expiresIn: -time.Hour,
			wantErr:   true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			token, err := MakeJWT(tc.userID, tc.secret, tc.expiresIn)
			if err != nil {
				t.Fatalf("MakeJWT failed: %v", err)
			}

			_, err = ValidateJWT(token, tc.secret)
			if (err != nil) != tc.wantErr {
				t.Errorf("ValidateJWT() error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}