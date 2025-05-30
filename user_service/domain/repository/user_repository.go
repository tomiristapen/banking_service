package repository

import (
	"context"

	"github.com/tomiristapen/banking_service/user_service/domain/model"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user *model.User) (*model.User, error)
	FindByEmail(ctx context.Context, email string) (*model.User, error)
	FindByID(ctx context.Context, id string) (*model.User, error)
	UpdateVerificationStatus(ctx context.Context, userID string, verified bool) error


	UpdateUser(ctx context.Context, user *model.User) error
}
