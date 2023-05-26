package models

import "github.com/golang-jwt/jwt"

type AddClaims struct {
	UserId string `json:"userId"`
	jwt.StandardClaims
}
