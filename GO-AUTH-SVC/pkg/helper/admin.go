package helper

import (
	"time"

	"github.com/akhi9550/pkg/config"
	"github.com/akhi9550/pkg/utils/models"
	"github.com/golang-jwt/jwt"
)

type AuthAdminClaims struct {
	Id      uint   `json:"id"`
	Email   string `json:"email"`
	Isadmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GenerateAccessTokenAdmin(user models.AdminResponse) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := GenerateAdminToken(user.Id, user.Email, user.Isadmin, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshTokenAdmin(user models.AdminResponse) (string, error) {
	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokenString, err := GenerateAdminToken(user.Id, user.Email, user.Isadmin, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateAdminToken(userID uint, Email string, isadmin bool, expirationTime time.Time) (string, error) {
	cfg, _ := config.LoadConfig()
	claims := &AuthAdminClaims{
		Id:      userID,
		Email:   Email,
		Isadmin: isadmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.KEY_ADMIN))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
