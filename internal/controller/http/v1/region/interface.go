package region

import (
	"context"
	"project/internal/repository/postgres/region"
)

type Region interface {
	GetList(ctx context.Context, filter region.Filter) ([]region.GetListResponse, int, error)
	GetDetailById(ctx context.Context, id int) (region.GetDetailByIdResponse, error)
	Create(ctx context.Context, request region.CreateRequest) (region.CreateResponse, error)
	UpdateAll(ctx context.Context, request region.UpdateRequest) error
	UpdateColumns(ctx context.Context, request region.UpdateRequest) error
	Delete(ctx context.Context, id int) error
	GetRegionByRepublicIDList(ctx context.Context, republicID int, filter region.Filter) ([]region.GetListByRepublicIDResponse, int, error)
}
