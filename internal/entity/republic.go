package entity

import (
	"github.com/uptrace/bun"
)

type Republic struct {
	bun.BaseModel `bun:"table:republic"`

	BasicEntity
	Name map[string]string `json:"name"     bun:"name"`
}
