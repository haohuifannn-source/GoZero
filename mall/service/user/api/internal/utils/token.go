package utils

import "github.com/golang-jwt/jwt/v4"

// 一般来说，还需要有个Role的KV
func GenerateAccessToken(secret string, iat, seconds, userId int64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"iat":    iat,
		"exp":    iat + seconds,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}

func GenerateRefreshToken(secret string, iat, seconds, userId int64) (string, error) {
	claims := jwt.MapClaims{
		"userId": userId,
		"iat":    iat,
		"exp":    iat + seconds,
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).
		SignedString([]byte(secret))
}
