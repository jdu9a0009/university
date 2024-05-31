package district

import (
	"context"
	"project/internal/repository/postgres/district"
)

type District interface {
	GetList(ctx context.Context, filter district.Filter) ([]district.GetListResponse, int, error)
	GetDetailById(ctx context.Context, id int) (district.GetDetailByIdResponse, error)
	Create(ctx context.Context, request district.CreateRequest) (district.CreateResponse, error)
	UpdateAll(ctx context.Context, request district.UpdateRequest) error
	UpdateColumns(ctx context.Context, request district.UpdateRequest) error
	Delete(ctx context.Context, id int) error
	GetDistrictsListByRegionID(ctx context.Context, regionID int, filter district.Filter) ([]district.GetListByRegionIDResponse, int, error)
}
