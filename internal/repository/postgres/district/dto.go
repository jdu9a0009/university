package district

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit    *int
	Offset   *int
	Page     *int
	Search   *string
	RegionID *int
}

type GetListResponse struct {
	ID            int      `json:"id"`
	Name          *string  `json:"name"`
	RegionID      *int     `json:"region_id" bun:"region_id"`
	RegionName    *string  `json:"region_name" bun:"region_name"`
	NameLanguages []string `json:"name_languages"`
}

type GetListByRegionIDResponse struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

type GetDetailByIdResponse struct {
	ID       int               `json:"id"`
	Name     map[string]string `json:"name"`
	RegionID *int              `json:"region_id"`
	Region   map[string]string `json:"region"`
}

type CreateRequest struct {
	Name     map[string]string `json:"name" form:"name"`
	RegionID *int              `json:"region_id" form:"region_id"`
}

type CreateResponse struct {
	bun.BaseModel `bun:"table:district"`

	ID int `json:"id" bun:"-"`

	Name      map[string]string `json:"name"       bun:"name"`
	RegionID  *int              `json:"region_id" bun:"region_id"`
	CreatedAt time.Time         `json:"-"          bun:"created_at"`
	CreatedBy int               `json:"-"          bun:"created_by"`
}

type UpdateRequest struct {
	ID       int               `json:"id" form:"id"`
	Name     map[string]string `json:"name" form:"name"`
	RegionID *int              `json:"region_id" form:"region_id"`
}
