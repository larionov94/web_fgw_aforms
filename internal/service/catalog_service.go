package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type CatalogService struct {
	catalogRepo repository.CatalogRepository
	logg        *common.Logger
}

func NewCatalogService(catalogRepo repository.CatalogRepository, logger *common.Logger) *CatalogService {
	return &CatalogService{catalogRepo, logger}
}

type CatalogUseCase interface {
	DesignNameAll(ctx context.Context) ([]*model.Catalog, error)
	ColorAll(ctx context.Context) ([]*model.Catalog, error)
}

func (c *CatalogService) DesignNameAll(ctx context.Context) ([]*model.Catalog, error) {
	designNameList, err := c.catalogRepo.DesignNameAll(ctx)
	if err != nil {
		c.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return designNameList, nil
}

func (c *CatalogService) ColorAll(ctx context.Context) ([]*model.Catalog, error) {
	colorList, err := c.catalogRepo.ColorAll(ctx)
	if err != nil {
		c.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return colorList, nil
}
