package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	ID    uint   `json:"user_id"`
	Email string `json:"email"`
	Role  string `json:"role"`
	jwt.RegisteredClaims
}

// Valid implements jwt.Claims.
func (j *JwtCustomClaims) Valid() error {
	panic("unimplemented")
}

var secretKey = []byte(os.Getenv("JWT_SECRET_KEY"))

func GenerateJWT(id uint, email string, role string) (string, error) {
	// Set custom claims
	customClaims := &JwtCustomClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "coffee_shop_app",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(15 * time.Minute)),
		},
	}

	// Create token with claims
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Generate encoded token and send it as response.
	token, err := jwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GenerateRefresherJWT(id uint, email string, role string) (string, error) {
	// Set custom claims
	customClaims := &JwtCustomClaims{
		ID:    id,
		Email: email,
		Role:  role,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)), // 7 days
			Issuer:    "coffee_shop_app",
		},
	}

	// Create token with claims
	refresherJwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)

	// Generate encoded token and send it as response.
	token, err := refresherJwtToken.SignedString(secretKey)
	if err != nil {
		return "", err
	}
	return token, nil
}

func GetSecretKey() []byte {
	return secretKey
}

// func ValidateToken(signedToken string) (*JwtCustomClaims, error) {
// 	token, err := jwt.ParseWithClaims(
// 		signedToken,
// 		&JwtCustomClaims{},
// 		func(token *jwt.Token) (interface{}, error) {
// 			return secretKey, nil
// 		},
// 	)
// 	if err != nil {
// 		return nil, err
// 	}
// 	claims, ok := token.Claims.(*JwtCustomClaims)
// 	if !ok {
// 		return nil, err
// 	}
// 	return claims, nil
// }
