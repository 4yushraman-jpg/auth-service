package service

import (
	"context"

	"github.com/google/uuid"

	appErrors "github.com/4yushraman-jpg/auth-service/internal/errors"
	"github.com/4yushraman-jpg/auth-service/internal/model"
	"github.com/4yushraman-jpg/auth-service/internal/repository"
)

type AuthService struct {
	userRepo        *repository.UserRepository
	passwordService *PasswordService
}

func NewAuthService(
	userRepo *repository.UserRepository,
	passwordService *PasswordService,
) *AuthService {
	return &AuthService{
		userRepo:        userRepo,
		passwordService: passwordService,
	}
}

func (s *AuthService) Register(
	ctx context.Context,
	email string,
	password string,
) error {
	existingUser, _ := s.userRepo.GetByEmail(ctx, email)

	if existingUser != nil {
		return appErrors.ErrUserAlreadyExists
	}

	hashedPassword, err := s.passwordService.Hash(password)
	if err != nil {
		return err
	}

	user := &model.User{
		ID:           uuid.New(),
		Email:        email,
		PasswordHash: hashedPassword,
		Role:         "user",
	}

	return s.userRepo.Create(ctx, user)
}
