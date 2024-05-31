package entity

import (
	"github.com/uptrace/bun"
)

type District struct {
	bun.BaseModel `bun:"table:district"`

	BasicEntity
	Name     map[string]string `json:"name"     bun:"name"`
	RegionID *int              `json:"region_id" bun:"region_id"`
}
