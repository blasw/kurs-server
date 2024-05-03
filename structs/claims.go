package structs

import "github.com/golang-jwt/jwt"

type JwtUserClaims struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	jwt.StandardClaims
}
