package auth

import (
	"context"
	"project/internal/entity"
)

type User interface {
	GetByUsername(ctx context.Context, username string) (entity.User, error)
}
