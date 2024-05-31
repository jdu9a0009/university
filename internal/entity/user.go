package entity

import (
	"time"

	"github.com/uptrace/bun"
)

type User struct {
	bun.BaseModel `bun:"table:users"`

	BasicEntity
	Avatar        *string    `json:"avatar"     bun:"avatar"`
	Username      *string    `json:"username"   bun:"username"`
	Password      *string    `json:"password"   bun:"password"`
	FullName      *string    `json:"full_name"  bun:"full_name"`
	Phone         *string    `json:"phone"      bun:"phone"`
	Role          *string    `json:"role"       bun:"role"`
	BirthDistrict *int       `json:"birth_district" bun:"birth_district"`
	BirthDate     *time.Time `json:"birth_date" bun:"birth_date"`
}
