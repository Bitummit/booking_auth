package my_jwt

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/Bitummit/booking_auth/internal/models"
	"github.com/golang-jwt/jwt/v4"
)

type UserClaims struct {
	Id int64
	Username string
	Role string
	ExpiresAt int64
}

// var ErrorInvalidToken = errors.New("invalid token")
var ErrorTokenDuration = errors.New("invalid token duration")
var ErrorSigningToken = errors.New("token signing error")
var ErrorTokenExpired = errors.New("token expired")

func (u UserClaims) Valid() error {
	if u.ExpiresAt < time.Now().Unix() {
		return fmt.Errorf("invalid token: %w", ErrorTokenExpired)
	}
	return nil
}

func NewToken(user models.User) (string, error) {
	ttl, err := time.ParseDuration(os.Getenv("TOKEN_TTL"))
	if err != nil {
		return "", fmt.Errorf("creating token: %w", ErrorTokenDuration)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaims{
		Id: user.Id,
		Username: user.Username,
		Role: user.Role,
		ExpiresAt: time.Now().Add(ttl).Unix(),
	})

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", fmt.Errorf("creating token: %w", ErrorSigningToken)
	}

	return tokenString, nil
}

// func ParseToken(tokenString string) (models.User, error) {
// 	var userClaims UserClaims
// 	_, err := jwt.ParseWithClaims(tokenString, &userClaims, func(token *jwt.Token) (interface{}, error) {
//     	return []byte(os.Getenv("SECRET_KEY")), nil
// 	})
// 	if err != nil {
// 		return models.User{}, fmt.Errorf("parsing token: %w", ErrorInvalidToken)
// 	}
	
// 	user := models.User{
// 		Id: userClaims.Id,
// 		Username: userClaims.Username,
// 	}
// 	return user, nil
// }
