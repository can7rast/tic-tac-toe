package application

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"school21/internal/domain"
	"school21/internal/infrastructure/datasource"
	"school21/internal/web/dto"
	"school21/pkg"
	"time"
)

type AuthService interface {
	SignUp(ctx context.Context, req dto.SignUpRequest) error
	Login(ctx context.Context, req dto.SignUpRequest) (*domain.User, error)
}

type authService struct {
	authRepo datasource.UserRepository
}

func (a *authService) SignUp(ctx context.Context, req dto.SignUpRequest) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	exitstUser, err := a.authRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return fmt.Errorf("data base error %w", err)
	}
	if exitstUser != nil {
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
	err = a.authRepo.Create(ctx, u)
	if err != nil {
		return err
	}
	return nil
}

func (a *authService) Login(ctx context.Context, req dto.SignUpRequest) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	u, err := a.authRepo.FindByLogin(ctx, req.Login)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}
	if err = pkg.CheckPassword(u.PasswordHash, req.Password); err != nil {
		return nil, fmt.Errorf("invalid password")
	}
	return u, nil
}

func NewAuthService(authRepo datasource.UserRepository) AuthService {
	return &authService{authRepo: authRepo}
}
