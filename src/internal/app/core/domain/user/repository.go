package user

import (
	"context"

	"github.com/google/uuid"
)

type IRepository interface {
	Create(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
	Update(ctx context.Context, user *User) error
	GetUser(ctx context.Context, id uuid.UUID) (*User, error)
	GetUsers(ctx context.Context, limit int, page int) ([]*User, int, error)
}
