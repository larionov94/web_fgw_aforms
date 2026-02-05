package repository

import (
	"context"
	"database/sql"
	"fgw_web_aforms/internal/config/db"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type SectorRepo struct {
	mssql *sql.DB
	logg  *common.Logger
}

func NewSectorRepo(mssql *sql.DB, logger *common.Logger) *SectorRepo {
	return &SectorRepo{mssql: mssql, logg: logger}
}

type SectorRepository interface {
	All(ctx context.Context) ([]*model.Sector, error)
}

func (s *SectorRepo) All(ctx context.Context) ([]*model.Sector, error) {
	rows, err := s.mssql.QueryContext(ctx, FGWsvAFormsSectorAllQuery)
	if err != nil {
		s.logg.LogE(msg.E3203, err)

		return nil, err
	}
	defer db.RowsClose(rows)

	var sectors []*model.Sector
	for rows.Next() {
		var sector model.Sector
		if err = rows.Scan(
			&sector.IdSector,
			&sector.SectorName,
			&sector.SectorEditDate,
			&sector.SectorEditUser,
			&sector.SecVMPL,
			&sector.PerformerId,
			&sector.Dtact,
			&sector.TicketSize,
		); err != nil {
			s.logg.LogE(msg.E3204, err)

			return nil, err
		}

		sectors = append(sectors, &sector)
	}
	if err = rows.Err(); err != nil {
		s.logg.LogE(msg.E3205, err)

		return nil, err
	}

	return sectors, nil
}
