package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web/dto"
	"school21/internal/web/jwt"
	"school21/pkg"
	"time"
)

type AuthService interface {
	SignUp(ctx context.Context, req dto.JwtRequest) error
	Login(ctx context.Context, req dto.JwtRequest) (*dto.JwtResponse, error)
	RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error)
	RefreshRefreshToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error)
	GetUserInfo(ctx context.Context, token string) (*dto.UserInfoResponse, error)
}

type authService struct {
	authRepo datasource.UserRepository
	jvt      *jwt.Provider
}

func (a *authService) GetUserInfo(ctx context.Context, token string) (*dto.UserInfoResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	claims, err := a.jvt.ValidateToken(token)
	if err != nil {
		return nil, fmt.Errorf("validate token not validate: %w", err)
	}
	login := claims.Login
	if login == "" {
		return nil, fmt.Errorf("login not valid")
	}
	id := claims.UserID
	if id == uuid.Nil {
		return nil, fmt.Errorf("user id is not valid")
	}
	return &dto.UserInfoResponse{
		ID:    id,
		Login: login,
	}, nil
}

func (a *authService) SignUp(ctx context.Context, req dto.JwtRequest) error {
	existsUser, err := a.authRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return fmt.Errorf("data base error %w", err)
	}
	if existsUser != nil {
		return fmt.Errorf("user with login %s already exists", req.Login)
	}

	PasswordHash, err := pkg.HashPassword(req.Password)
	if err != nil {
		return err
	}
	u := &domain.User{
		ID:           uuid.New(),
		Login:        req.Login,
		PasswordHash: PasswordHash,
	}
	return a.authRepo.Create(ctx, u)
}

func (a *authService) Login(ctx context.Context, req dto.JwtRequest) (*dto.JwtResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	u, err := a.authRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("data base error %w", err)
	}
	if u == nil {
		return nil, fmt.Errorf("user not found")
	}
	if err = pkg.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	accessToken, err := a.jvt.GenerateAccessToken(u.ID, u.Login)
	if err != nil {
		return nil, fmt.Errorf("can't generate access token in Login")
	}
	refreshToken, err := a.jvt.GenerateRefreshToken(u.ID)
	if err != nil {
		return nil, fmt.Errorf("can't generate refresh token in Login")
	}
	return &dto.JwtResponse{
		Token:        "Bearer",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (a *authService) RefreshAccessToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	claims, err := a.jvt.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}
	NewAccessToken, err := a.jvt.GenerateAccessToken(claims.UserID, claims.Login)
	if err != nil {
		return nil, fmt.Errorf("can't regenerate access token")
	}

	return &dto.JwtResponse{
		Token:        "Bearer",
		AccessToken:  NewAccessToken,
		RefreshToken: refreshToken,
	}, nil

}
func (a *authService) RefreshRefreshToken(ctx context.Context, refreshToken string) (*dto.JwtResponse, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	claims, err := a.jvt.ValidateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token")
	}
	NewRefreshToken, err := a.jvt.GenerateRefreshToken(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("can't regenerate refresh token")
	}
	NewAccessToken, err := a.jvt.GenerateAccessToken(claims.UserID, claims.Login)
	if err != nil {
		return nil, fmt.Errorf("can't regenerate access token")
	}

	return &dto.JwtResponse{
		Token:        "Bearer",
		AccessToken:  NewAccessToken,
		RefreshToken: NewRefreshToken,
	}, nil

}

func NewAuthService(authRepo datasource.UserRepository, provider *jwt.Provider) AuthService {
	return &authService{
		authRepo: authRepo,
		jvt:      provider,
	}
}
