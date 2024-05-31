package entity

import (
	"github.com/uptrace/bun"
)

type Region struct {
	bun.BaseModel `bun:"table:region"`

	BasicEntity
	Name       map[string]string `json:"name"     bun:"name"`
	RepublicID *int              `json:"republic_id" bun:"republic_id"`
}