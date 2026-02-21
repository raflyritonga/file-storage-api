package service

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/raflyritonga/file-storage-api/internal/dto"
	"github.com/raflyritonga/file-storage-api/internal/model"
	"github.com/raflyritonga/file-storage-api/internal/repository"
	"github.com/raflyritonga/file-storage-api/pkg/helper/auth"
)

type AuthService struct {
    userRepo    repository.UserRepository
    sessionRepo repository.SessionRepository
    tokenMgr    *auth.TokenManager
}

func NewAuthService(userRepo repository.UserRepository, sessionRepo repository.SessionRepository, tokenMgr *auth.TokenManager) *AuthService {
    return &AuthService{userRepo: userRepo, sessionRepo: sessionRepo, tokenMgr: tokenMgr}
}

func (s *AuthService) Register(ctx context.Context, req dto.RegisterRequest) error {
    _, err := s.userRepo.FindByEmail(ctx, req.Email)
    if err == nil {
        return errors.New("email already registered")
    }

    if !errors.Is(err, sql.ErrNoRows) {
        return err
    }

    hashedPass, err := auth.HashPassword(req.Password)
    if err != nil {
        return err
    }

    user := &model.User{
        ID:           uuid.NewString(),
        Email:        req.Email,
		Username: 		req.Username,
        Password: 		hashedPass,
        Role:        req.Role,
    }

    return s.userRepo.Create(ctx, user)
}

func (s *AuthService) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
    user, err := s.userRepo.FindByUsername(ctx, req.Username)
    if err != nil {
        return nil, errors.New("invalid credentials")
    }

    if err := auth.ComparePassword(user.Password, req.Password); err != nil {
        return nil, errors.New("invalid credentials")
    }

    sessionID := uuid.NewString()

    accessToken, _, accessExp, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Role)
    if err != nil {
        return nil, err
    }

    refreshToken, refreshJTI, refreshExp, err := s.tokenMgr.GenerateRefreshToken(user.ID, sessionID)
    if err != nil {
        return nil, err
    }

    refreshTTL := int64(time.Until(refreshExp).Seconds())
    if err := s.sessionRepo.SaveRefreshSession(ctx, user.ID, sessionID, refreshJTI, user.Role, refreshTTL); err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int(time.Until(accessExp).Seconds()),
    }, nil
}

func (s *AuthService) Refresh(ctx context.Context, rawRefreshToken string) (*dto.AuthResponse, error) {
    claims, err := s.tokenMgr.ParseRefreshToken(rawRefreshToken)
    if err != nil {
        return nil, errors.New("invalid refresh token")
    }

    isValid, err := s.sessionRepo.ValidateRefreshSession(ctx, claims.Subject, claims.SessionID, claims.ID)
    if err != nil || !isValid {
        return nil, errors.New("refresh session invalid")
    }

    user, err := s.userRepo.FindByID(ctx, claims.Subject)
    if err != nil {
        return nil, err
    }

    _ = s.sessionRepo.DeleteRefreshSession(ctx, claims.Subject, claims.SessionID)

    newSessionID := uuid.NewString()
    accessToken, _, accessExp, err := s.tokenMgr.GenerateAccessToken(user.ID, user.Role)
    if err != nil {
        return nil, err
    }

    refreshToken, refreshJTI, refreshExp, err := s.tokenMgr.GenerateRefreshToken(user.ID, newSessionID)
    if err != nil {
        return nil, err
    }

    refreshTTL := int64(time.Until(refreshExp).Seconds())
    if err := s.sessionRepo.SaveRefreshSession(ctx, user.ID, newSessionID, refreshJTI, user.Role, refreshTTL); err != nil {
        return nil, err
    }

    return &dto.AuthResponse{
        AccessToken:  accessToken,
        RefreshToken: refreshToken,
        TokenType:    "Bearer",
        ExpiresIn:    int(time.Until(accessExp).Seconds()),
    }, nil
}