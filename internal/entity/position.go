package entity

import (
	"github.com/uptrace/bun"
)

type Position struct {
	bun.BaseModel `bun:"table:position"`

	BasicEntity
	Name map[string]string `json:"name"     bun:"name"`
}
