package service

import (
	"context"

	"telecomx-provisioning-service/internal/domain/model"
	"telecomx-provisioning-service/internal/infrastructure/adapter/repository"
)

type ProvisioningService struct {
	repo *repository.MongoRepository
}

func NewProvisioningService(repo *repository.MongoRepository) *ProvisioningService {
	return &ProvisioningService{repo: repo}
}

func (s *ProvisioningService) Create(ctx context.Context, p *model.Provisioning) error {
	return s.repo.Create(ctx, p)
}

func (s *ProvisioningService) UpdateStatus(ctx context.Context, userID, status string) error {
	return s.repo.UpdateStatus(ctx, userID, status)
}

func (s *ProvisioningService) Delete(ctx context.Context, userID string) error {
	return s.repo.DeleteByUserID(ctx, userID)
}

func (s *ProvisioningService) GetAll(ctx context.Context) ([]model.Provisioning, error) {
	return s.repo.GetAll(ctx)
}
