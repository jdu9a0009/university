package entity

import(
	"github.com/uptrace/bun"
)

type Department struct {
	bun.BaseModel `bun:"table:departments"`
	
	BasicEntity
	Name map[string]string `json:"name"     bun:"name"`
}