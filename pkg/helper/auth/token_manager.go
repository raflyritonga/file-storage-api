package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type TokenManager struct {
    accessSecret   []byte
    refreshSecret  []byte
    issuer         string
    accessTTLMin   int
    refreshTTLHour int
}

func NewTokenManager(accessSecret, refreshSecret, issuer string, accessTTLMin, refreshTTLHour int) *TokenManager {
    return &TokenManager{
        accessSecret:   []byte(accessSecret),
        refreshSecret:  []byte(refreshSecret),
        issuer:         issuer,
        accessTTLMin:   accessTTLMin,
        refreshTTLHour: refreshTTLHour,
    }
}

func (m *TokenManager) GenerateAccessToken(userID, role string) (token string, jti string, exp time.Time, err error) {
    jti = uuid.NewString()
    exp = time.Now().Add(time.Duration(m.accessTTLMin) * time.Minute)

    claims := AccessClaims{
        Role: role,
        Type: "access",
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   userID,
            Issuer:    m.issuer,
            ID:        jti,
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(exp),
        },
    }

    signed := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err = signed.SignedString(m.accessSecret)
    return
}

func (m *TokenManager) GenerateRefreshToken(userID, sessionID string) (token string, jti string, exp time.Time, err error) {
    jti = uuid.NewString()
    exp = time.Now().Add(time.Duration(m.refreshTTLHour) * time.Hour)

    claims := RefreshClaims{
        SessionID: sessionID,
        Type:      "refresh",
        RegisteredClaims: jwt.RegisteredClaims{
            Subject:   userID,
            Issuer:    m.issuer,
            ID:        jti,
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            ExpiresAt: jwt.NewNumericDate(exp),
        },
    }

    signed := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    token, err = signed.SignedString(m.refreshSecret)
    return
}

func (m *TokenManager) ParseAccessToken(raw string) (*AccessClaims, error) {
    token, err := jwt.ParseWithClaims(raw, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
        return m.accessSecret, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*AccessClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid access token")
    }

    if claims.Type != "access" || claims.Issuer != m.issuer {
        return nil, errors.New("invalid access claims")
    }

    return claims, nil
}

func (m *TokenManager) ParseRefreshToken(raw string) (*RefreshClaims, error) {
    token, err := jwt.ParseWithClaims(raw, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
        return m.refreshSecret, nil
    })
    if err != nil {
        return nil, err
    }

    claims, ok := token.Claims.(*RefreshClaims)
    if !ok || !token.Valid {
        return nil, errors.New("invalid refresh token")
    }

    if claims.Type != "refresh" || claims.Issuer != m.issuer {
        return nil, errors.New("invalid refresh claims")
    }

    return claims, nil
}