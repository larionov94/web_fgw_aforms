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

type PerformerRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewPerformerRepo(mssql *sql.DB, logger *common.Logger) *PerformerRepo {
	return &PerformerRepo{mssql: mssql, logg: logger}
}

type PerformerRepository interface {
	AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error)
	FindById(ctx context.Context, id int) (*model.Performer, error)
}

// AuthByIdAndPass проверка существования в БД сотрудника.
func (p *PerformerRepo) AuthByIdAndPass(ctx context.Context, id int, password string) (bool, error) {
	var authSuccess bool

	err := p.mssql.QueryRowContext(ctx, FGWsvPerformerAuthQuery, id, password).Scan(&authSuccess)
	if err != nil {
		p.logg.LogE(msg.E3202, err)

		return false, err
	}

	return authSuccess, nil
}

// FindById ищет сотрудника по ИД.
func (p *PerformerRepo) FindById(ctx context.Context, id int) (*model.Performer, error) {
	var performer model.Performer

	if err := p.mssql.QueryRowContext(ctx, FGWsvPerformerFindByIdQuery, id).Scan(
		&performer.Id,
		&performer.FIO,
		&performer.BC,
		&performer.Pass,
		&performer.Archive,
		&performer.IdRoleAForms,
		&performer.IdRoleAFGW,
		&performer.AuditRec.CreatedAt,
		&performer.AuditRec.CreatedBy,
		&performer.AuditRec.UpdatedAt,
		&performer.AuditRec.UpdatedBy,
	); err != nil {
		p.logg.LogE(msg.E3204, err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %v", msg.E3206, err)
		}
		return nil, err
	}

	return &performer, nil
}
