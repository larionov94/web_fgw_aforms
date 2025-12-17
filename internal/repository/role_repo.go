package repository

import (
	"context"
	"database/sql"
	"errors"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"fmt"
)

type RoleRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewRoleRepo(mssql *sql.DB, logger *common.Logger) *RoleRepo {
	return &RoleRepo{mssql: mssql, logg: logger}
}

type RoleRepository interface {
	FindById(ctx context.Context, id int) (*model.Role, error)
}

func (r *RoleRepo) FindById(ctx context.Context, id int) (*model.Role, error) {
	var role model.Role

	if err := r.mssql.QueryRowContext(ctx, FGWsvRoleFindByIdQuery, id).Scan(
		&role.Id,
		&role.Name,
		&role.Desc,
		&role.AuditRec.CreatedAt,
		&role.AuditRec.CreatedBy,
		&role.AuditRec.UpdatedAt,
		&role.AuditRec.UpdatedBy,
	); err != nil {
		r.logg.LogE(msg.E3204, err)

		if errors.Is(err, sql.ErrNoRows) {
			r.logg.LogE(msg.E3206, err)

			return nil, err
		}
		return nil, fmt.Errorf("%s: %v", msg.E3202, err)
	}

	return &role, nil
}
