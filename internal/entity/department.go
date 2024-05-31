package entity

import(
	"github.com/uptrace/bun"
)

type Department struct {
	bun.BaseModel `bun:"table:department"`
	
	BasicEntity
	Name map[string]string `json:"name"     bun:"name"`
}