package utlis

import (
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

// claims
type Claims struct {
	UserId  uint   `json:"userid"`
	Email   string `json:"email"`
	Role    string `json:"role"`
	TokenID string `josn:"tokenid"`
	jwt.RegisteredClaims
}

// generate jwt token
func Generate(userid uint, email, role string) (string, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading env file")
	}
	Jwtkey := []byte(os.Getenv("jwtkey"))

	expir := time.Now().Add(60 * time.Minute)
	claims := &Claims{
		UserId: userid,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expir),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(Jwtkey)
}

// refresh token
func Refresh(userid uint, email, role string) (string, error) {
	refreshkey := []byte(os.Getenv("refreshkey"))

	exp := time.Now().Add(7 * 24 * time.Hour)
	refreshid := uuid.New().String()

	claims := &Claims{
		UserId:  userid,
		Email:   email,
		Role:    role,
		TokenID: refreshid,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(exp),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "app",
			ID:        refreshid,
		},
	}
	refresh := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := refresh.SignedString(refreshkey)
	if err != nil {
		return "invalid refresh token", err
	}

	// storing token in redis
	err = StoreRefresh(refreshid, userid, time.Until(exp))
	if err != nil {
		return "failed to store token in redis", nil
	}

	return token, nil
}
