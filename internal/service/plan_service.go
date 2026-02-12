package service

import (
	"context"
	"errors"
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
	AllPlans(ctx context.Context, sortField, sortOrder, startDate, endDate string, idProduction, idSector *int) ([]*model.Plan, error)
	AddPlan(ctx context.Context, p *model.Plan) error
}

func (p *PlanService) AllPlans(ctx context.Context, sortField, sortOrder, startDate, endDate string, idProduction, idSector *int) ([]*model.Plan, error) {
	plans, err := p.planRepo.All(ctx, sortField, sortOrder, startDate, endDate, idProduction, idSector)
	if err != nil {
		p.logg.LogE(msg.E3209, err)

		return nil, err
	}

	return plans, nil
}

func (p *PlanService) AddPlan(ctx context.Context, plan *model.Plan) error {
	if plan == nil {
		return errors.New("сменно-суточное задание не должно быть nil")
	}

	return p.planRepo.Add(ctx, plan)
}
