package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type SectorService struct {
	sectorRepo repository.SectorRepository
	logg       *common.Logger
}

func NewSectorService(repo repository.SectorRepository, logger *common.Logger) *SectorService {
	return &SectorService{repo, logger}
}

type SectorUseCase interface {
	AllSector(ctx context.Context) ([]*model.Sector, error)
}

func (s *SectorService) AllSector(ctx context.Context) ([]*model.Sector, error) {
	sectors, err := s.sectorRepo.All(ctx)
	if err != nil {
		s.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return sectors, nil
}
