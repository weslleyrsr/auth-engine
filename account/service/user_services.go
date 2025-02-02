package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/weslleyrsr/auth-engine/account/model"
)

type UserService struct {
	UserRepository model.UserRepository
}

type USConfig struct {
	UserRepository model.UserRepository
}

// NewUserService is a factory function for initializing a UserService with its repository layer dependencies
func NewUserService(c *USConfig) model.UserService {
	return &UserService{
		UserRepository: c.UserRepository,
	}
}

// Get retrieves a user based on their uuid
func (s *UserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	u, err := s.UserRepository.FindByID(ctx, uid)

	return u, err
}

// Signup just panics
func (s *UserService) Signup(ctx context.Context, u *model.User) error {
	panic("Method not implemented yet.")
}
