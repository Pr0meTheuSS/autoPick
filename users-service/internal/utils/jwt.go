package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	accessSecret  = []byte("your_access_secret_key")
	refreshSecret = []byte("your_refresh_secret_key")
)

// GenerateTokens создает access и refresh токены
func GenerateTokens(userID string) (string, string, error) {
	accessToken, err := generateJWT(userID, accessSecret, 15*time.Minute)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := generateJWT(userID, refreshSecret, 7*24*time.Hour)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// ValidateToken проверяет JWT и возвращает userID
func ValidateToken(tokenString string, isRefresh bool) (string, error) {
	var secret []byte
	if isRefresh {
		secret = refreshSecret
	} else {
		secret = accessSecret
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("неподдерживаемый метод подписи")
		}
		return secret, nil
	})

	if err != nil || !token.Valid {
		return "", errors.New("невалидный токен")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("ошибка получения claims")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", errors.New("ошибка получения user_id")
	}

	return userID, nil
}

// generateJWT создает токен
func generateJWT(userID string, secret []byte, duration time.Duration) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(secret)
}
