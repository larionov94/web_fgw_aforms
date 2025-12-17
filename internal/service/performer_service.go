package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type PerformerService struct {
	performerRepo repository.PerformerRepository
	logg          *common.Logger
}

func NewPerformerService(performerRepo repository.PerformerRepository, logger *common.Logger) *PerformerService {
	return &PerformerService{performerRepo: performerRepo, logg: logger}
}

type PerformerUseCase interface {
	AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error)
	FindByIdPerformer(ctx context.Context, id int) (*model.Performer, error)
}

func (p *PerformerService) AuthPerformer(ctx context.Context, id int, password string) (*model.AuthPerformer, error) {
	if id <= 0 || password == "" {
		p.logg.LogE(msg.E3211, nil)

		return &model.AuthPerformer{Success: false, Message: msg.E3211}, nil
	}

	authOK, err := p.performerRepo.AuthByIdAndPass(ctx, id, password)
	if err != nil {
		p.logg.LogE(msg.E3210, err)

		return &model.AuthPerformer{Success: false, Message: msg.E3210}, err
	}

	if !authOK {
		p.logg.LogE(msg.E3210, err)

		return &model.AuthPerformer{Success: false, Message: msg.E3210}, err
	}

	performer, err := p.performerRepo.FindById(ctx, id)
	if err != nil {
		return &model.AuthPerformer{Success: false, Message: msg.E3212}, err
	}

	return &model.AuthPerformer{
		Success:   true,
		Performer: *performer,
		Message:   "Успешный вход",
	}, nil
}

func (p *PerformerService) FindByIdPerformer(ctx context.Context, id int) (*model.Performer, error) {
	performer, err := p.performerRepo.FindById(ctx, id)
	if err != nil {
		p.logg.LogE(msg.E3212, err)

		return nil, err
	}

	return performer, nil
}
