package republic

import (
	"context"
	"project/internal/repository/postgres/republic"
)

type Republic interface {
	GetList(ctx context.Context, filter republic.Filter) ([]republic.GetListResponse, int, error)
	GetDetailById(ctx context.Context, id int) (republic.GetDetailByIdResponse, error)
	Create(ctx context.Context, request republic.CreateRequest) (republic.CreateResponse, error)
	UpdateAll(ctx context.Context, request republic.UpdateRequest) error
	UpdateColumns(ctx context.Context, request republic.UpdateRequest) error
	Delete(ctx context.Context, id int) error
}
