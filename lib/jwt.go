package lib

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.RegisteredClaims
}

type Tokens struct {
	Token           string `json:"token"`
	TokenExp        int64  `json:"tokenExp"`
	RefreshToken    string `json:"refreshToken"`
	RefreshTokenExp int64  `json:"refreshTokenExp"`
}

func NewClaims(id, email string, expiry time.Time) *Claims {
	return &Claims{
		Id:    id,
		Email: email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiry),
			ID:        id,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
}

func GenerateTokens(id, email string) (Tokens, error) {
	tokenExpTime := time.Now().Add(24 * time.Hour)
	refTokenExpTime := time.Now().Add(24 * 7 * time.Hour)

	tokenClaims := NewClaims(id, email, tokenExpTime)
	refTokenClaims := NewClaims(id, email, refTokenExpTime)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refTokenClaims)

	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_TOKEN_SECRET")))

	if err != nil {
		return Tokens{}, err
	}

	refTokenStr, err1 := refToken.SignedString([]byte(os.Getenv("JWT_REFRESH_TOKEN_SECRET")))

	if err1 != nil {
		return Tokens{}, err1
	}

	fmt.Println("Access token", tokenString)
	fmt.Println("Refresh token", refTokenStr)

	return Tokens{Token: tokenString, TokenExp: tokenExpTime.Unix(), RefreshToken: refTokenStr, RefreshTokenExp: refTokenExpTime.Unix()}, nil
}

func ValidateToken(tk string, secret string) (Claims, bool) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tk, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil || !token.Valid {
		return Claims{}, false
	}

	return Claims{Id: claims.Id, Email: claims.Email, RegisteredClaims: claims.RegisteredClaims}, true
}
