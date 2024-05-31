package region

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"project/foundation/web"
	"project/internal/entity"
	"project/internal/pkg/repository/postgresql"
	"project/internal/repository/postgres"

	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Repository struct {
	*postgresql.Database
}

func NewRepository(database *postgresql.Database) *Repository {
	return &Repository{Database: database}
}

func (r Repository) GetById(ctx context.Context, id int) (entity.Region, error) {
	var detail entity.Region

	err := r.NewSelect().Model(&detail).Where("id = ?", id).Scan(ctx)

	return detail, err
}

func (r Repository) GetList(ctx context.Context, filter Filter) ([]GetListResponse, int, error) {
	_, err := r.CheckClaims(ctx)
	if err != nil {
		return nil, 0, err
	}
	lang := r.GetLang(ctx)

	whereQuery := fmt.Sprintf(`
			WHERE 
				rg.deleted_at IS NULL
			`)

	if filter.Search != nil {
		whereQuery += fmt.Sprintf(` AND 
			(
				rg.name->>'%s' ilike '%s'
			)`,
			lang, "%"+*filter.Search+"%",
		)
	}

	if filter.RepublicID != nil {
		whereQuery += fmt.Sprintf(` AND rg.republic_id = %d`, *filter.RepublicID)
	}

	orderQuery := "ORDER BY rg.created_at desc"

	var limitQuery, offsetQuery string

	if filter.Page != nil && filter.Limit != nil {
		offset := (*filter.Page - 1) * (*filter.Limit)
		filter.Offset = &offset
	}

	if filter.Limit != nil {
		limitQuery += fmt.Sprintf(" LIMIT %d", *filter.Limit)
	}

	if filter.Offset != nil {
		offsetQuery += fmt.Sprintf(" OFFSET %d", *filter.Offset)
	}

	query := fmt.Sprintf(`
		SELECT
			rg.id,
			rg.name,
			rg.republic_id,
			rp.name
		FROM
		    region rg
		LEFT JOIN republic rp ON rg.republic_id = rp.id
		%s %s %s %s
	`, whereQuery, orderQuery, limitQuery, offsetQuery)

	rows, err := r.QueryContext(ctx, query)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting regions"), http.StatusBadRequest)
	}

	var list []GetListResponse

	for rows.Next() {
		var detail GetListResponse
		var nameByte []byte
		var republicNameByte []byte
		if err = rows.Scan(&detail.ID, &nameByte, &detail.RepublicID, &republicNameByte); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning region"), http.StatusBadRequest)
		}

		if nameByte != nil {
			name := make(map[string]string)
			if err = json.Unmarshal(nameByte, &name); err != nil {
				return nil, 0, web.NewRequestError(errors.Wrap(err, "region name unmarshal"), http.StatusBadRequest)
			}
			for k, v := range name {
				if k == lang {
					n := v
					detail.Name = &n
				}
				if v != "" {
					detail.NameLanguages = append(detail.NameLanguages, k)
				}
			}
		}

		if republicNameByte != nil {
			republicName := make(map[string]string)
			if err = json.Unmarshal(republicNameByte, &republicName); err != nil {
				return nil, 0, web.NewRequestError(errors.Wrap(err, "region republicName unmarshal"), http.StatusBadRequest)
			}
			for k, v := range republicName {
				if k == lang {
					n := v
					detail.RepublicName = &n
				}
				if v != "" {
					detail.RepublicLanguages = append(detail.RepublicLanguages, k)
				}
			}
		}

		list = append(list, detail)
	}

	countQuery := fmt.Sprintf(`
		SELECT
			count(rg.id)
		FROM
		    region rg
		LEFT JOIN republic rp ON rg.republic_id = rp.id
		%s
	`, whereQuery)

	countRows, err := r.QueryContext(ctx, countQuery)
	if err == sql.ErrNoRows {
		return nil, 0, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}
	if err != nil {
		return nil, 0, web.NewRequestError(errors.Wrap(err, "selecting regions"), http.StatusBadRequest)
	}

	count := 0

	for countRows.Next() {
		if err = countRows.Scan(&count); err != nil {
			return nil, 0, web.NewRequestError(errors.Wrap(err, "scanning region count"), http.StatusBadRequest)
		}
	}

	return list, count, nil
}

func (r Repository) GetDetailById(ctx context.Context, id int) (GetDetailByIdResponse, error) {
	_, err := r.CheckClaims(ctx)
	if err != nil {
		return GetDetailByIdResponse{}, err
	}

	query := fmt.Sprintf(`
		SELECT
			rg.id,
			rg.name,
			rg.republic_id,
			rp.name
		FROM
				region rg
		LEFT JOIN republic as rp ON rg.republic_id = rp.id
		WHERE rg.deleted_at IS NULL AND rg.id = %d
	`, id)

	var detail GetDetailByIdResponse
	var nameByte []byte
	var nameRepublicByte []byte

	err = r.QueryRowContext(ctx, query).Scan(
		&detail.ID,
		&nameByte,
		&detail.RepublicID,
		&nameRepublicByte,
	)

	if err == sql.ErrNoRows {
		return GetDetailByIdResponse{}, web.NewRequestError(postgres.ErrNotFound, http.StatusNotFound)
	}

	if nameByte != nil {
		if err = json.Unmarshal(nameByte, &detail.Name); err != nil {
			return GetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "region name unmarshal"), http.StatusBadRequest)
		}
	}
	if nameRepublicByte != nil {
		if err = json.Unmarshal(nameRepublicByte, &detail.Republic); err != nil {
			return GetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "republic name unmarshal"), http.StatusBadRequest)
		}
	}
	if err != nil {
		return GetDetailByIdResponse{}, web.NewRequestError(errors.Wrap(err, "selecting region detail"), http.StatusBadRequest)
	}

	return detail, nil
}

func (r Repository) Create(ctx context.Context, request CreateRequest) (CreateResponse, error) {
	claims, err := r.CheckClaims(ctx)
	if err != nil {
		return CreateResponse{}, err
	}

	if err := r.ValidateStruct(&request, "Name", "RepublicID"); err != nil {
		return CreateResponse{}, err
	}

	var response CreateResponse

	response.Name = request.Name
	response.RepublicID = request.RepublicID
	response.CreatedAt = time.Now()
	response.CreatedBy = claims.UserId

	_, err = r.NewInsert().Model(&response).Returning("id").Exec(ctx, &response.ID)
	if err != nil {
		return CreateResponse{}, web.NewRequestError(errors.Wrap(err, "creating region"), http.StatusBadRequest)
	}

	return response, nil
}

func (r Repository) UpdateAll(ctx context.Context, request UpdateRequest) error {
	if err := r.ValidateStruct(&request, "ID", "Name", "RepublicID"); err != nil {
		return err
	}

	claims, err := r.CheckClaims(ctx)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("region").Where("deleted_at IS NULL AND id = ?", request.ID)

	q.Set("name = ?", request.Name)
	q.Set("republic_id = ?", request.RepublicID)
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating region"), http.StatusBadRequest)
	}

	return nil
}

func (r Repository) UpdateColumns(ctx context.Context, request UpdateRequest) error {
	if err := r.ValidateStruct(&request, "ID"); err != nil {
		return err
	}

	claims, err := r.CheckClaims(ctx)
	if err != nil {
		return err
	}

	q := r.NewUpdate().Table("region").Where("deleted_at IS NULL AND id = ?", request.ID)

	if request.Name != nil {
		q.Set("name = ?", request.Name)
	}

	if request.RepublicID != nil {
		q.Set("republic_id = ?", request.RepublicID)
	}
	q.Set("updated_at = ?", time.Now())
	q.Set("updated_by = ?", claims.UserId)

	_, err = q.Exec(ctx)
	if err != nil {
		return web.NewRequestError(errors.Wrap(err, "updating region"), http.StatusBadRequest)
	}

	return nil
}


func (r Repository) Delete(ctx context.Context, id int) error {
	return r.DeleteRow(ctx, "region", id)
}
