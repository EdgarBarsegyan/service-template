package user

import (
	"context"
	"service-template/internal/domain/user"

	"github.com/google/uuid"
)

type Repository interface {
	Create(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *user.User) error
	GetUser(ctx context.Context, id uuid.UUID) (*user.User, error)
	GetUsers(ctx context.Context, limit int, page int) ([]*user.User, int, error)
}
