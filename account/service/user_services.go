package service

import (
	"context"
	"github.com/google/uuid"
	"github.com/weslleyrsr/auth-engine/account/model"
	"github.com/weslleyrsr/auth-engine/account/model/apperrors"
	"log"
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

// Signup reaches out to a UserRepository to verify that the email address is available
// and sign up the user if this is the case
func (s *UserService) Signup(ctx context.Context, u *model.User) error {
	pw, err := hashPassword(u.Password)
	if err != nil {
		log.Printf("Unable to signup user fr email: %v\n", u.Email)
		return apperrors.NewInternal()
	}

	u.Password = pw

	if err := s.UserRepository.Create(ctx, u); err != nil {
		return err
	}

	//if we get around to adding events, we'd publish it here
	//err:= s.EventsBroker.PublishUserUpdated(u, true)

	//if err := nil {
	//	return nil, apperrors.NewInternal()
	//}

	return nil
}
