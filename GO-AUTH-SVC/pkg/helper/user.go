package helper

import (
	"errors"
	"time"

	"github.com/akhi9550/pkg/config"
	"github.com/akhi9550/pkg/utils/models"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type AuthUserClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Isadmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GenerateAccessTokenUser(user models.UserResponse) (string, error) {
	expirationTime := time.Now().Add(15 * time.Minute)
	tokenString, err := GenerateToken(user.Id, user.Username, user.Isadmin, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateRefreshTokenUser(user models.UserResponse) (string, error) {
	expirationTime := time.Now().Add(24 * 90 * time.Hour)
	tokenString, err := GenerateToken(user.Id, user.Username, user.Isadmin, expirationTime)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func GenerateToken(userID uint, userName string, isadmin bool, expirationTime time.Time) (string, error) {
	cfg, _ := config.LoadConfig()
	claims := &AuthUserClaims{
		Id:       userID,
		Username: userName,
		Isadmin:  isadmin,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(cfg.KEY))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func PasswordHash(password string) (string, error) {
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashPassword)
	return hash, nil
}

func PasswordHashing(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return "", errors.New("internal server error")
	}
	hash := string(hashedPassword)
	return hash, nil
}

func CompareHashAndPassword(a string, b string) error {
	err := bcrypt.CompareHashAndPassword([]byte(a), []byte(b))
	if err != nil {
		return err
	}
	return nil
}