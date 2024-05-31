package user

import (
	"mime/multipart"
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit  *int
	Offset *int
	Page   *int
	Search *string
	Role   *string
}

type SignInRequest struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

type AuthClaims struct {
	ID   int
	Role string
}

type GetListResponse struct {
	ID            int     `json:"id"`
	Avatar        *string `json:"avatar"`
	FullName      *string `json:"full_name"`
	Username      *string `json:"username"`
	Role          *string `json:"role"`
	BirthDistrict *string `json:"birth_district_id"`
	BirthDate     *string `json:"birth_date"`
}

type GetDetailByIdResponse struct {
	ID                int     `json:"id"`
	Avatar            *string `json:"avatar"`
	Username          *string `json:"username"`
	FullName          *string `json:"full_name"`
	Role              *string `json:"role"`
	BirthDistrict     *int    `json:"birth_district_id"`
	BirthDistrictName *string `json:"birth_district_name"`
	BirthDate         *string `json:"birth_date"`
}

type CreateRequest struct {
	Username      *string               `json:"username" form:"username"`
	FullName      *string               `json:"full_name" form:"full_name"`
	Password      *string               `json:"password" form:"password"`
	Avatar        *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink    *string               `json:"-" form:"-"`
	Role          *string               `json:"role" form:"role"`
	BirthDistrict *string               `json:"birth_district_id" form:"birth_district_id"`
	BirthDate     *string               `json:"birth_date" form:"birth_date"`
}

type CreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID            int        `json:"id" bun:"-"`
	Avatar        *string    `json:"avatar"     bun:"avatar"`
	Username      *string    `json:"username"   bun:"username"`
	Password      *string    `json:"-"   bun:"password"`
	FullName      *string    `json:"full_name" bun:"full_name"`
	Role          *string    `json:"role" bun:"role"`
	BirthDistrict *string    `json:"birth_district_id" bun:"birth_district_id"`
	BirthDate     *time.Time `json:"birth_date" bun:"birth_date"`
	CreatedAt     time.Time  `json:"-"          bun:"created_at"`
	CreatedBy     int        `json:"-"          bun:"created_by"`
}

type UpdateRequest struct {
	ID            int                   `json:"id" form:"id"`
	Username      *string               `json:"username" form:"username"`
	FullName      *string               `json:"full_name" form:"full_name"`
	Password      *string               `json:"password" form:"password"`
	Avatar        *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink    *string               `json:"-" form:"-"`
	Role          *string               `json:"role" form:"role"`
	BirthDistrict *string               `json:"birth_district_id" form:"birth_district_id"`
	BirthDate     *string               `json:"birth_date" form:"birth_date"`
}
