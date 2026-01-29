package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type PlanService struct {
	planRepo repository.PlanRepository
	logg     *common.Logger
}

func NewPlanService(planRepo repository.PlanRepository, logger *common.Logger) *PlanService {
	return &PlanService{planRepo: planRepo, logg: logger}
}

type PlanUseCase interface {
	AllPlans(ctx context.Context, sortField, sortOrder string) ([]*model.Plan, error)
}

func (p *PlanService) AllPlans(ctx context.Context, sortField, sortOrder string) ([]*model.Plan, error) {
	plans, err := p.planRepo.All(ctx, sortField, sortOrder)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return plans, nil
}
