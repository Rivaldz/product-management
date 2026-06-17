package usecase

import (
	"context"

	"technical_test/internal/entity"
)

type CompanyUseCase struct {
	repo CompanyRepo
}

func NewCompanyUseCase(r CompanyRepo) *CompanyUseCase {
	return &CompanyUseCase{repo: r}
}

func (uc *CompanyUseCase) GetList(ctx context.Context) ([]entity.Company, error) {
	return uc.repo.GetList(ctx)
}
