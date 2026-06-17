package usecase

import (
	"context"

	"technical_test/internal/entity"
)

type ItemUseCase struct {
	repo ItemRepo
}

func NewItemUseCase(r ItemRepo) *ItemUseCase {
	return &ItemUseCase{repo: r}
}

func (uc *ItemUseCase) Create(ctx context.Context, item entity.Item) (entity.Item, error) {
	item.Status = entity.StatusActive
	return uc.repo.Insert(ctx, item)
}

func (uc *ItemUseCase) GetList(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
	return uc.repo.GetByFilter(ctx, companyID, filter)
}

func (uc *ItemUseCase) GetDetail(ctx context.Context, companyID, itemID string) (entity.Item, error) {
	return uc.repo.GetByID(ctx, companyID, itemID)
}

func (uc *ItemUseCase) Update(ctx context.Context, item entity.Item) (entity.Item, error) {
	existing, err := uc.repo.GetByID(ctx, item.CompanyID, item.ID)
	if err != nil {
		return entity.Item{}, err
	}
	if existing.Status == entity.StatusArchived {
		return entity.Item{}, entity.ErrAlreadyArchived
	}

	return uc.repo.Update(ctx, item)
}

func (uc *ItemUseCase) Archive(ctx context.Context, companyID, itemID string) error {
	existing, err := uc.repo.GetByID(ctx, companyID, itemID)
	if err != nil {
		return err
	}
	if existing.Status == entity.StatusArchived {
		return entity.ErrAlreadyArchived
	}

	return uc.repo.Archive(ctx, companyID, itemID)
}
