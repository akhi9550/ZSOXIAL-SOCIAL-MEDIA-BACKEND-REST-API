package helper

import (
	"fmt"

	"github.com/akhi9550/api-gateway/pkg/config"
	"github.com/golang-jwt/jwt"
)

type AuthUserClaims struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Isadmin  bool   `json:"is_admin"`
	jwt.StandardClaims
}

type AuthAdminClaims struct {
	Id      uint   `json:"id"`
	Email   string `json:"email"`
	Isadmin bool   `json:"is_admin"`
	jwt.StandardClaims
}

func GetTokenFromHeader(header string) string {
	if len(header) > 7 && header[:7] == "Bearer " {
		return header[7:]
	}
	return header
}
func ExtractUserIDFromToken(tokenString string) (int, string, error) {
	cfg, _ := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenString, &AuthUserClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(cfg.KEY), nil
	})

	if err != nil {
		fmt.Println("errors:-", err)
		return 0, "", err
	}

	claims, ok := token.Claims.(*AuthUserClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid token claims")
	}

	return int(claims.Id), claims.Username, nil

}

func ExtractAdminFromToken(tokenString string) (*AuthAdminClaims, error) {
	cfg, _ := config.LoadConfig()
	token, err := jwt.ParseWithClaims(tokenString, &AuthAdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("invalid signing method")
		}
		return []byte(cfg.KEY_ADMIN), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*AuthAdminClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	return claims, nil

}
