package repository

import (
	"context"
	"database/sql"
	"fgw_web_aforms/internal/config/db"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type PlanRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewPlanRepo(mssql *sql.DB, logger *common.Logger) *PlanRepo {
	return &PlanRepo{mssql: mssql, logg: logger}
}

type PlanRepository interface {
	All(ctx context.Context, sortField, sortOrder string) ([]*model.Plan, error)
}

func (p *PlanRepo) All(ctx context.Context, sortField, sortOrder string) ([]*model.Plan, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvAFormsPlanAllQuery, sortField, sortOrder)
	if err != nil {
		p.logg.LogE(msg.E3203, err)

		return nil, err
	}

	defer db.RowsClose(rows)

	var plans []*model.Plan
	for rows.Next() {
		var plan model.Plan
		if err = rows.Scan(
			&plan.IdPlan,
			&plan.IdPlan,
			&plan.PlanShift,
			&plan.ExtProduction,
			&plan.ExtSector,
			&plan.PlanCount,
			&plan.PlanDate,
			&plan.PlanInfo,
			&plan.PlEditDate,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}

		plans = append(plans, &plan)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return plans, nil
}
