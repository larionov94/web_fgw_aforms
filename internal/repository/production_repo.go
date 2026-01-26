package repository

import (
	"context"
	"database/sql"
	"errors"
	"fgw_web_aforms/internal/config/db"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"fmt"
)

type ProductionRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewProductionRepo(mssql *sql.DB, logg *common.Logger) *ProductionRepo {
	return &ProductionRepo{mssql: mssql, logg: logg}
}

type ProductionRepository interface {
	All(ctx context.Context, sortField, sortOrder string) ([]*model.Production, error)
	Filter(ctx context.Context, articlePattern, namePattern, idPattern string) ([]*model.Production, error)
	Add(ctx context.Context, production *model.Production) error
	Upd(ctx context.Context, idProduction int, production *model.Production) error
	FindById(ctx context.Context, idProduction int) (*model.Production, error)
}

func (p *ProductionRepo) All(ctx context.Context, sortField, sortOrder string) ([]*model.Production, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvAFormsProductionAllQuery, sortField, sortOrder)
	if err != nil {
		p.logg.LogE(msg.E3203, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var productions []*model.Production
	for rows.Next() {
		var production model.Production
		if err = rows.Scan(
			&production.IdProduction,
			&production.PrShortName,
			&production.PrPackName,
			&production.PrArticle,
			&production.PrColor,
			&production.PrCount,
			&production.PrRows,
			&production.PrHWD,
			&production.PrEditDate,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}
		productions = append(productions, &production)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return productions, nil
}

func (p *ProductionRepo) Filter(ctx context.Context, articlePattern, namePattern, idPattern string) ([]*model.Production, error) {
	rows, err := p.mssql.QueryContext(ctx, FGWsvAFormsProductionFilterQuery, articlePattern, namePattern, idPattern)
	if err != nil {
		p.logg.LogE(msg.E3206, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var productions []*model.Production
	for rows.Next() {
		var production model.Production

		if err = rows.Scan(
			&production.IdProduction,
			&production.PrShortName,
			&production.PrPackName,
			&production.PrArticle,
			&production.PrColor,
			&production.PrCount,
			&production.PrRows,
			&production.PrHWD,
			&production.PrEditDate,
		); err != nil {
			p.logg.LogE(msg.E3204, err)

			return nil, err
		}

		productions = append(productions, &production)
	}

	if err = rows.Err(); err != nil {
		p.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return productions, nil
}

func (p *ProductionRepo) Add(ctx context.Context, production *model.Production) error {
	if _, err := p.mssql.ExecContext(ctx, FGWsvAFormsProductionAddQuery,
		&production.PrName,
		&production.PrShortName,
		&production.PrPackName,
		&production.PrDecl,
		&production.PrSun,
		&production.PrProdType,
		&production.PrParty,
		&production.PrUmbrella,
		&production.PrPerfumery,
		&production.PrColor,
		&production.PrGL,
		&production.PrArticle,
		&production.PrSAP,
		&production.PrCount,
		&production.PrRows,
		&production.PrWeight,
		&production.PrHWD,
		&production.PrInfo,
		&production.PrPart,
		&production.PrPartLastDate,
		&production.PrPartAutoInc,
		&production.PrPerGodn,
		&production.AuditRec.CreatedBy,
		&production.AuditRec.UpdatedBy,
	); err != nil {
		p.logg.LogE(msg.E3204, err)

		return err
	}
	return nil
}

func (p *ProductionRepo) Upd(ctx context.Context, idProduction int, production *model.Production) error {
	_, err := p.mssql.ExecContext(ctx, FGWsvAFormsProductionUpdQuery, idProduction,
		production.PrName,
		production.PrShortName,
		production.PrPackName,
		production.PrDecl,
		production.PrSun,
		production.PrProdType,
		production.PrParty,
		production.PrUmbrella,
		production.PrPerfumery,
		production.PrColor,
		production.PrGL,
		production.PrArticle,
		production.PrSAP,
		production.PrCount,
		production.PrRows,
		production.PrWeight,
		production.PrHWD,
		production.PrInfo,
		production.PrPart,
		production.PrPartLastDate,
		production.PrPartAutoInc,
		production.PrPerGodn,
		production.AuditRec.CreatedBy,
		production.AuditRec.UpdatedBy,
	)
	if err != nil {
		p.logg.LogE(msg.E3216, err)

		return err
	}

	return nil
}

func (p *ProductionRepo) FindById(ctx context.Context, idProduction int) (*model.Production, error) {
	var production model.Production

	if err := p.mssql.QueryRowContext(ctx, FGWsvAFormsProductionFindByIdQuery, idProduction).Scan(
		&production.IdProduction,
		&production.PrName,
		&production.PrShortName,
		&production.PrPackName,
		&production.PrType,
		&production.PrArticle,
		&production.PrColor,
		&production.PrBarCode,
		&production.PrCount,
		&production.PrRows,
		&production.PrWeight,
		&production.PrHWD,
		&production.PrInfo,
		&production.PrStatus,
		&production.PrEditDate,
		&production.PrEditUser,
		&production.PrPart,
		&production.PrPartLastDate,
		&production.PrPartAutoInc,
		&production.PrPartRealDate,
		&production.PrArchive,
		&production.PrPerGodn,
		&production.PrSAP,
		&production.PrProdType,
		&production.PrUmbrella,
		&production.PrSun,
		&production.PrDecl,
		&production.PrParty,
		&production.PrGL,
		&production.PrVP,
		&production.PrML,
		&production.AuditRec.CreatedAt,
		&production.AuditRec.CreatedBy,
		&production.AuditRec.UpdatedAt,
		&production.AuditRec.UpdatedBy,
		&production.PrCreationDate,
		&production.PrPerfumery,
	); err != nil {
		p.logg.LogE(msg.E3204, err)

		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%s: %v", msg.E3206, err)
		}

		return nil, err
	}

	return &production, nil
}
