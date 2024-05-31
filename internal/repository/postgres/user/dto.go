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
	ID            int        `json:"id"`
	Avatar        *string    `json:"avatar"`
	FullName      *string    `json:"full_name"`
	Username      *string    `json:"username"`
	Phone         *string    `json:"phone"`
	Role          *string    `json:"role"`
	BirthDistrict *string    `json:"birth_district"`
	BirthDate     *time.Time `json:"birth_date"`
}

type GetDetailByIdResponse struct {
	ID            int        `json:"id"`
	Avatar        *string    `json:"avatar"`
	Username      *string    `json:"username"`
	FullName      *string    `json:"full_name"`
	Phone         *string    `json:"phone"`
	Role          *string    `json:"role"`
	BirthDistrict *string    `json:"birth_district"`
	BirthDate     *time.Time `json:"birth_date"`
}

type CreateRequest struct {
	Username      *string               `json:"username" form:"username"`
	FullName      *string               `json:"full_name" form:"full_name"`
	Password      *string               `json:"password" form:"password"`
	Phone         *string               `json:"phone" form:"phone"`
	Avatar        *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink    *string               `json:"-" form:"-"`
	Role          *string               `json:"role" form:"role"`
	BirthDistrict *string               `json:"birth_district" form:"birth_district"`
	BirthDate     *time.Time            `json:"birth_date" form:"birth_date"`
}

type CreateResponse struct {
	bun.BaseModel `bun:"table:users"`

	ID            int        `json:"id" bun:"-"`
	Avatar        *string    `json:"avatar"     bun:"avatar"`
	Username      *string    `json:"username"   bun:"username"`
	Password      *string    `json:"-"   bun:"password"`
	FullName      *string    `json:"full_name" bun:"full_name"`
	Phone         *string    `json:"phone"      bun:"phone"`
	Role          *string    `json:"role" bun:"role"`
	BirthDistrict *string    `json:"birth_district" bun:"birth_district"`
	BirthDate     *time.Time `json:"birth_date" bun:"birth_date"`
	CreatedAt     time.Time  `json:"-"          bun:"created_at"`
	CreatedBy     int        `json:"-"          bun:"created_by"`
}

type UpdateRequest struct {
	ID            int                   `json:"id" form:"id"`
	Username      *string               `json:"username" form:"username"`
	FullName      *string               `json:"full_name" form:"full_name"`
	Password      *string               `json:"password" form:"password"`
	Phone         *string               `json:"phone" form:"phone"`
	Avatar        *multipart.FileHeader `json:"-" form:"avatar"`
	AvatarLink    *string               `json:"-" form:"-"`
	Role          *string               `json:"role" form:"role"`
	BirthDistrict *string               `json:"birth_district" form:"birth_district"`
	BirthDate     *time.Time            `json:"birth_date" form:"birth_date"`
}
