package auth

import "github.com/golang-jwt/jwt/v5"

type AccessClaims struct {
    Role string `json:"role"`
    Type string `json:"typ"`
    jwt.RegisteredClaims
}

type RefreshClaims struct {
    SessionID string `json:"sid"`
    Type      string `json:"typ"`
    jwt.RegisteredClaims
}