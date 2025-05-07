package utils

import (
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func GenerateJWT(userID string, isTime int16) (string, error) {
	if isTime <= 0 {
		isTime = 1
	}
	claims := jwt.MapClaims{
		"sub": userID,
		"exp": time.Now().Add(time.Hour * time.Duration(isTime)).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(os.Getenv("JWT_SECRET")))
}

func ValidateJWT(tokenStr string) (string, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return "", err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims["sub"].(string), nil
}

func RefreshToken(tokenStr string) (string, string, error) {

	payload, err := ValidateJWT(tokenStr)

	if err != nil {
		return "", "", err
	}

	newAccessToken, err := GenerateJWT(payload, 3)
	refreshToken, err := GenerateJWT(payload, 168)

	if err != nil {
		return "", "", err
	}

	return newAccessToken, refreshToken, nil
}
