package repository

import (
	"context"
	"database/sql"
	"fgw_web_aforms/internal/config/db"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type CatalogRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewCatalogRepo(mssql *sql.DB, logger *common.Logger) *CatalogRepo {
	return &CatalogRepo{mssql: mssql, logg: logger}
}

type CatalogRepository interface {
	DesignNameAll(ctx context.Context) ([]*model.Catalog, error)
	ColorAll(ctx context.Context) ([]*model.Catalog, error)
}

func (c *CatalogRepo) DesignNameAll(ctx context.Context) ([]*model.Catalog, error) {
	rows, err := c.mssql.QueryContext(ctx, FGWsvAFormsDesignNameAllQuery)
	if err != nil {
		c.logg.LogE(msg.E3203, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var catalogs []*model.Catalog
	for rows.Next() {
		var catalog model.Catalog
		if err = rows.Scan(
			&catalog.Id,
			&catalog.Name,
		); err != nil {
			c.logg.LogE(msg.E3204, err)

			return nil, err
		}

		catalogs = append(catalogs, &catalog)
	}

	if err = rows.Err(); err != nil {
		c.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return catalogs, nil
}

func (c *CatalogRepo) ColorAll(ctx context.Context) ([]*model.Catalog, error) {
	rows, err := c.mssql.QueryContext(ctx, FGWsvAFormsColorAllQuery)
	if err != nil {
		c.logg.LogE(msg.E3203, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var colors []*model.Catalog
	for rows.Next() {
		var color model.Catalog
		if err = rows.Scan(
			&color.Name,
			&color.DopInt1,
		); err != nil {
			c.logg.LogE(msg.E3204, err)

			return nil, err
		}

		colors = append(colors, &color)
	}

	if err = rows.Err(); err != nil {
		c.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return colors, nil
}
