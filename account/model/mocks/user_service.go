package mocks

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/weslleyrsr/auth-engine/model"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Get(ctx context.Context, uid uuid.UUID) (*model.User, error) {
	res := m.Called(ctx, uid)

	var r0 *model.User
	if res.Get(0) != nil {
		r0 = res.Get(0).(*model.User)
	}

	var r1 error
	if res.Get(1) != nil {
		r1 = res.Get(1).(error)
	}

	return r0, r1
}
