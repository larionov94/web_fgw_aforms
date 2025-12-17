package service

import (
	"context"
	"fgw_web_aforms/internal/model"
	"fgw_web_aforms/internal/repository"
	"fgw_web_aforms/pkg/common"
	"fgw_web_aforms/pkg/common/msg"
)

type RoleService struct {
	roleRepo repository.RoleRepository
	logg     *common.Logger
}

func NewRoleService(roleRepo repository.RoleRepository, logger *common.Logger) *RoleService {
	return &RoleService{roleRepo: roleRepo, logg: logger}
}

type RoleUseCase interface {
	FindRoleById(ctx context.Context, id int) (*model.Role, error)
}

func (r *RoleService) FindRoleById(ctx context.Context, id int) (*model.Role, error) {
	role, err := r.roleRepo.FindById(ctx, id)
	if err != nil {
		r.logg.LogE(msg.E3212, err)

		return nil, err
	}

	return role, nil
}
