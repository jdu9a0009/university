package entity

import (
	"github.com/uptrace/bun"
)

type Position struct {
	bun.BaseModel `bun:"table:positions"`

	BasicEntity
	Name map[string]string `json:"name"     bun:"name"`
}
