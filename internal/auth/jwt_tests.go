package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeAndValidateJWT(t *testing.T) {
	// make token
	expiresIn, err := time.ParseDuration("1m")
	if err != nil {
		t.Error("Failed to create time.Duration")
	}

	userID := uuid.New()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	tokenSecret := "testy"
	signedToken, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Errorf("Error signing token: %v", err)
	}

	// validate token made above
	claims := &jwt.RegisteredClaims{}

	jwtTokenToCheck, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", err)
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		t.Errorf("Error parsing claims: %v", err)
	}

	if !jwtTokenToCheck.Valid {
		t.Errorf("Error validating token: %v", err)
	}
}

func TestExpiredJWTTokens(t *testing.T) {
	// make token
	expiresIn, err := time.ParseDuration("1ms")
	if err != nil {
		t.Error("Failed to create time.Duration")
	}

	userID := uuid.New()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	tokenSecret := "testy"
	signedToken, err := jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Errorf("Error signing token: %v", err)
	}

	// intentionally wait 50ms so that token expires
	time.Sleep(time.Millisecond * 50)

	// validate token made above
	claims := &jwt.RegisteredClaims{}

	jwtTokenToCheck, err := jwt.ParseWithClaims(signedToken, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", err)
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		t.Errorf("Error parsing claims: %v", err)
	}

	if jwtTokenToCheck.Valid {
		t.Error("Token is valid when it should be invalid due to expiration")
	}
}

func TestWrongSecretJWTTokens(t *testing.T) {
	// make token
	expiresIn, err := time.ParseDuration("5s")
	if err != nil {
		t.Error("Failed to create time.Duration")
	}

	userID := uuid.New()

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject:   userID.String(),
	})

	// true secret
	tokenSecret := "testy"
	_, err = jwtToken.SignedString([]byte(tokenSecret))
	if err != nil {
		t.Errorf("Error signing token: %v", err)
	}

	// validate token made above
	claims := &jwt.RegisteredClaims{}

	// intentionally put in wrong secret to parse with
	jwtTokenToCheck, err := jwt.ParseWithClaims("wrong secret", claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			t.Errorf("Unexpected signing method: %v", err)
		}
		return []byte(tokenSecret), nil
	})
	if err != nil {
		t.Errorf("Error parsing claims: %v", err)
	}

	if jwtTokenToCheck.Valid {
		t.Error("Token is valid when it should be invalid from wrong secret")
	}
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name          string
		authHeader    string
		expectedToken string
		expectError   bool
	}{
		{
			name:          "Valid Bearer token",
			authHeader:    "Bearer abc.def.ghi",
			expectedToken: "abc.def.ghi",
			expectError:   false,
		},
		{
			name:        "Missing Authorization header",
			authHeader:  "",
			expectError: true,
		},
		{
			name:        "Authorization header with only Bearer",
			authHeader:  "Bearer",
			expectError: true,
		},
		{
			name:        "Malformed Authorization header",
			authHeader:  "SomethingElse abc.def.ghi",
			expectError: true,
		},
		{
			name:          "Authorization header with extra spaces",
			authHeader:    "   Bearer    abc.def.ghi   ",
			expectedToken: "abc.def.ghi",
			expectError:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			headers := http.Header{}
			if tt.authHeader != "" {
				headers.Set("Authorization", tt.authHeader)
			}

			token, err := GetBearerToken(headers)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Did not expect error, got: %v", err)
				}
				if token != tt.expectedToken {
					t.Errorf("Expected token %q, got %q", tt.expectedToken, token)
				}
			}
		})
	}
}
