package auth

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)



const (
	MinCost     int = 4  // the minimum allowable cost as passed in to GenerateFromPassword
	MaxCost     int = 31 // the maximum allowable cost as passed in to GenerateFromPassword
	DefaultCost int = 10 // the cost that will actually be set if a cost below MinCost is passed into GenerateFromPassword
	// expiresIn = 5
)


func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil 

}
func CheckPassworhash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	
}

func MakeJWT(userID uuid.UUID, tokenSecret string, expiresIn time.Duration) (string, error) {
	
	signingKey := []byte(tokenSecret)

	token:= jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer: "chirpy",
		IssuedAt: jwt.NewNumericDate(time.Now().UTC()),
		ExpiresAt: jwt.NewNumericDate(time.Now().UTC().Add(expiresIn)),
		Subject: userID.String(),

	})
	signedToken, err := token.SignedString(signingKey)
	if err != nil {
		return "", fmt.Errorf("error signing token: %w", err)
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecret string) (uuid.UUID, error) {
	claims := &jwt.RegisteredClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(tokenSecret), nil
	})

	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing token: %w", err)
	}

	if !token.Valid {
		return uuid.Nil, fmt.Errorf("invalid token")
	}

	if claims.Issuer != "chirpy" {
    return uuid.Nil, fmt.Errorf("invalid token issuer")
	}

	if claims.ExpiresAt == nil || time.Now().After(claims.ExpiresAt.Time) {
    return uuid.Nil, fmt.Errorf("token has expired")
	}


	userID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, fmt.Errorf("error parsing user ID:%w", err)
	}
	return userID, nil

}

func GetBearerToken(headers http.Header) (string, error) {
	authHeader:= headers.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("no header")
	}
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 {
		return "", errors.New("malformed authorization header")

	}

	if parts[0] != "Bearer" {
		return "", errors.New("authorization header must start with Bearer")
	}
	return parts[1], nil
}	


func MakeRefreshToken() (string, error) {
	data := make([]byte, 32)

	_, err := rand.Read(data)
	if err != nil {
		fmt.Println("error generating random data:", err)
		return "", err
	}

	encodedData := hex.EncodeToString(data)
	return encodedData, nil
}
	