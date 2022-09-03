package auth

import (
	"os"
	"time"

	"avila-common/model"

	"github.com/dgrijalva/jwt-go"
)

var secretString = os.Getenv("JWT_SECRET")
var jwtKey = []byte(secretString)
var defaultExpirationTime = 1 // TODO GET FROM ENV

type JWTClaim struct {
	UserId     int64       `json:"user_id"`
	Email      string      `json:"email"`
	CustomerId int64       `json:"customer_id"`
	Roles      model.Roles `json:"roles"`
	jwt.StandardClaims
}

func GenerateJwt(userId int64, email string, customerId int64, roles model.Roles) (tokenString string, err error) {

	expirationTime := time.Now().Add(time.Duration(defaultExpirationTime) * time.Hour)
	claims := &JWTClaim{
		UserId:     userId,
		Email:      email,
		CustomerId: customerId,
		Roles:      roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func GenerateJwtAndRefresh(userId int64, email string, customerId int64, roles model.Roles) (tokenString string, tokenRefresh string, err error) {

	expirationTime := time.Now().Add(time.Duration(defaultExpirationTime) * time.Hour)
	claims := &JWTClaim{
		UserId:     userId,
		Email:      email,
		CustomerId: customerId,
		Roles:      roles,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response.
	// The signing string should be secret (a generated UUID works too)
	tokenString, err = token.SignedString(jwtKey)
	if err != nil {

		return "", "", err
	}

	refreshToken := jwt.New(jwt.SigningMethodHS256)
	refreshClaims := refreshToken.Claims.(jwt.MapClaims)
	refreshClaims["exp"] = time.Now().Add(time.Hour * 24 * 7).Unix()

	refreshString, err := refreshToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	return tokenString, refreshString, err
}

func ValidateJWT(tokenString string) (claims *JWTClaim, err error) {

	token, err := jwt.ParseWithClaims(tokenString, &JWTClaim{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JWTClaim); ok && token.Valid {
		return claims, nil
	}
	return nil, err
}
