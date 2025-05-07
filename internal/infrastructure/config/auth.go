package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)


type JwtConfig struct {
	Secret string
	SigningMethod jwt.SigningMethod
	RefreshTokenTime time.Duration
	AccessTokenTime time.Duration
}


func LoadJwtConfig() (*JwtConfig, error){
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("")
	}

	refreshTokenTime, err := strconv.Atoi(os.Getenv("JWT_REFRESH_TOKEN_TIME"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT_REFRESH_TOKEN_TIME: %w", err)
	}
	accessTokenTime, err := strconv.Atoi(os.Getenv("JWT_ACCESS_TOKEN_TIME"))
	if err != nil {
		return nil, fmt.Errorf("failed to parse JWT_ACCESS_TOKEN_TIME: %w", err)
	}

	return &JwtConfig{
		Secret: os.Getenv("JWT_SECRET"),
		SigningMethod: jwt.GetSigningMethod(os.Getenv("JWT_SIGNING_METHOD")),
		RefreshTokenTime: time.Duration(refreshTokenTime) * time.Hour * 24 * 7,
		AccessTokenTime: time.Duration(accessTokenTime) * time.Minute,
	}, nil
}