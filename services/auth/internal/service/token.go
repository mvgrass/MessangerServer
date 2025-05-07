package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserId string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateToken(userId, accessJwtSecret, refreshJwtSecret string) (accessToken, refreshToken string, err error) {

	accessTokenExp := time.Now().Add(15 * time.Minute)
	accessClaims := UserClaims{
		UserId: userId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessTokenExp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "myapp",
		},
	}

	accessToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims).SignedString([](byte)(accessJwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to create access token: %w", err)
	}

	refreshTokenExp := time.Now().Add(7 * 24 * time.Hour)
	refreshClaims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(refreshTokenExp),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "myapp",
		Subject:   userId,
	}

	refreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([](byte)(refreshJwtSecret))
	if err != nil {
		return "", "", fmt.Errorf("failed to create refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func ParseRefreshToken(tokenString, secretString string) (*jwt.RegisteredClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwt.RegisteredClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretString), nil
		},
	)

	if claims, ok := token.Claims.(*jwt.RegisteredClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

func ParseAccessTokenWithoutExparation(tokenString, secretString string) (*UserClaims, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		&UserClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secretString), nil
		},
	)

	if claims, ok := token.Claims.(*UserClaims); ok {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return claims, nil
		}

		if err != nil {
			return nil, fmt.Errorf("error during parsing access token:%w", err)
		}
		return claims, nil
	}

	return nil, err
}
