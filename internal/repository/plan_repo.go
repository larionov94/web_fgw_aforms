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
	All(ctx context.Context, sortField, sortOrder, startDate, endDate string, idProduction, idSector *int) ([]*model.Plan, error)
	Add(ctx context.Context, p *model.Plan) error
}

func (p *PlanRepo) All(ctx context.Context, sortField, sortOrder, startDate, endDate string, idProduction, idSector *int) ([]*model.Plan, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvAFormsPlansWithSortingAndFilteringQuery, sortField, sortOrder, startDate, endDate, idProduction, idSector)
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
			&plan.PlanShift,
			&plan.ExtProduction,
			&plan.ProductionModel.PrName,
			&plan.ProductionModel.PrShortName,
			&plan.ProductionModel.PrType,
			&plan.ProductionModel.PrArticle,
			&plan.ProductionModel.PrColor,
			&plan.ProductionModel.PrCount,
			&plan.ProductionModel.PrRows,
			&plan.ProductionModel.PrHWD,
			&plan.ProductionModel.PrWeight,
			&plan.ExtSector,
			&plan.SectorModel.SectorName,
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

func (p *PlanRepo) Add(ctx context.Context, plan *model.Plan) error {
	if _, err := p.mssql.ExecContext(ctx, FGWsvAFormsPlanAddQuery,
		&plan.PlanShift,
		&plan.PlanDate,
		&plan.ExtSector,
		&plan.PlanCount,
		&plan.ExtProduction,
		&plan.PlanInfo,
		&plan.AuditRec.CreatedBy,
		&plan.AuditRec.UpdatedBy,
	); err != nil {
		p.logg.LogE(msg.E3204, err)

		return err
	}

	return nil
}
