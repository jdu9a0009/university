package region

import (
	"time"

	"github.com/uptrace/bun"
)

type Filter struct {
	Limit      *int
	Offset     *int
	Page       *int
	Search     *string
	RepublicID *int
}

type GetListResponse struct {
	ID                int      `json:"id"`
	Name              *string  `json:"name"`
	RepublicID        *int     `json:"republic_id" bun:"republic_id"`
	RepublicName      *string  `json:"republic_name" bun:"republic_name"`
	NameLanguages     []string `json:"name_languages"`
	RepublicLanguages []string `json:"republic_languages"`
}

type GetListByRepublicIDResponse struct {
	ID   int     `json:"id"`
	Name *string `json:"name"`
}

type GetDetailByIdResponse struct {
	ID         int               `json:"id"`
	Name       map[string]string `json:"name"`
	RepublicID *int              `json:"republic_id"`
	Republic   map[string]string `json:"republic"`
}

type CreateRequest struct {
	Name       map[string]string `json:"name" form:"name"`
	RepublicID *int              `json:"republic_id" form:"republic_id"`
}

type CreateResponse struct {
	bun.BaseModel `bun:"table:region"`

	ID int `json:"id" bun:"-"`

	Name       map[string]string `json:"name"       bun:"name"`
	RepublicID *int              `json:"republic_id" bun:"republic_id"`
	CreatedAt  time.Time         `json:"-"          bun:"created_at"`
	CreatedBy  int               `json:"-"          bun:"created_by"`
}

type UpdateRequest struct {
	ID         int               `json:"id" form:"id"`
	Name       map[string]string `json:"name" form:"name"`
	RepublicID *int              `json:"republic_id" form:"republic_id"`
}
