package lib

import (
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type Claims struct {
	Id    string `json:"id"`
	Email string `json:"email"`
	jwt.StandardClaims
}

type Tokens struct {
	Token           string
	TokenExp        int64
	RefreshToken    string
	RefreshTokenExp int64
}

func NewClaims(id, email string, expiry int64) *Claims {
	return &Claims{
		Id:    id,
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: int64(expiry),
		},
	}
}

func GenerateTokens(id, email string) (Tokens, error) {
	tokenExpTime := time.Now().Add(24 * time.Hour)
	refTokenExpTime := time.Now().Add(24 * 7 * time.Hour)

	tokenClaims := NewClaims(id, email, tokenExpTime.Unix())
	refTokenClaims := NewClaims(id, email, refTokenExpTime.Unix())

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims)
	refToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refTokenClaims)

	tokenString, err := token.SignedString(os.Getenv("JWT_TOKEN_SECRET"))

	if err != nil {
		return Tokens{}, err
	}

	refTokenStr, err1 := refToken.SignedString(os.Getenv("JWT_REFRESH_TOKEN_SECRET"))

	if err1 != nil {
		return Tokens{}, err1
	}

	return Tokens{Token: tokenString, TokenExp: tokenExpTime.Unix(), RefreshToken: refTokenStr, RefreshTokenExp: refTokenExpTime.Unix()}, nil
}
