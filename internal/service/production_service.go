package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type ProductionService struct {
	productionRepo repository.ProductionRepository
	logg           *common.Logger
}

func NewProductionService(production repository.ProductionRepository, logg *common.Logger) *ProductionService {
	return &ProductionService{production, logg}
}

type ProductionUserCase interface {
	AllProductions(ctx context.Context) ([]*model.Production, error)
}

func (p *ProductionService) AllProductions(ctx context.Context) ([]*model.Production, error) {
	productions, err := p.productionRepo.All(ctx)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return productions, nil
}
