package service

import (
	"context"
	"errors"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
	"strings"
)

type ProductionService struct {
	productionRepo repository.ProductionRepository
	logg           *common.Logger
}

func NewProductionService(production repository.ProductionRepository, logg *common.Logger) *ProductionService {
	return &ProductionService{production, logg}
}

type ProductionUseCase interface {
	AllProductions(ctx context.Context, sortField, sortOrder string) ([]*model.Production, error)
	SearchProductions(ctx context.Context, articlePattern, namePattern, idPattern string) ([]*model.Production, error)
	AddProduction(ctx context.Context, production *model.Production) error
}

func (p *ProductionService) AllProductions(ctx context.Context, sortField, sortOrder string) ([]*model.Production, error) {
	productions, err := p.productionRepo.All(ctx, sortField, sortOrder)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return productions, nil
}

func (p *ProductionService) SearchProductions(ctx context.Context, articlePattern, namePattern, idPattern string) ([]*model.Production, error) {
	articlePattern = strings.TrimSpace(articlePattern)
	namePattern = strings.TrimSpace(namePattern)
	idPattern = strings.TrimSpace(idPattern)

	return p.productionRepo.Filter(ctx, articlePattern, namePattern, idPattern)
}

func (p *ProductionService) AddProduction(ctx context.Context, production *model.Production) error {
	if production == nil {
		return errors.New("продукция не должна быть nil")
	}

	return p.productionRepo.Add(ctx, production)
}
