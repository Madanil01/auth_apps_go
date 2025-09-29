package utils

import (
	"time"
	"github.com/golang-jwt/jwt/v5"
)

// func GenerateToken(userID uint, email string, duration time.Duration) (string, error) {
// 	claims := jwt.MapClaims{
// 		"user_id": userID,
// 		"email":   email,
// 		"exp":     time.Now().Add(duration).Unix(),
// 	}

// 	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
// 	return token.SignedString(SecretKey)
// }

func GenerateTokens(userID uint, email string) (accessToken string, refreshToken string, err error) {
	// Access Token (misalnya 15 menit)
	accessClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(120 * time.Minute).Unix(),
	}
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessToken, err = at.SignedString(SecretKey)
	if err != nil {
		return
	}

	// Refresh Token (misalnya 7 hari)
	refreshClaims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
	}
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshToken, err = rt.SignedString(SecretKey)
	return
}