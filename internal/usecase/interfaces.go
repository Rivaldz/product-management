package usecase

import (
	"context"

	"technical_test/internal/entity"
)

type (
	Item interface {
		Create(ctx context.Context, item entity.Item) (entity.Item, error)
		GetList(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error)
		GetDetail(ctx context.Context, companyID, itemID string) (entity.Item, error)
		Update(ctx context.Context, item entity.Item) (entity.Item, error)
		Archive(ctx context.Context, companyID, itemID string) error
	}

	ItemRepo interface {
		Insert(ctx context.Context, item entity.Item) (entity.Item, error)
		GetByFilter(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error)
		GetByID(ctx context.Context, companyID, itemID string) (entity.Item, error)
		Update(ctx context.Context, item entity.Item) (entity.Item, error)
		Archive(ctx context.Context, companyID, itemID string) error
	}

	Company interface {
		GetList(ctx context.Context) ([]entity.Company, error)
	}

	CompanyRepo interface {
		GetList(ctx context.Context) ([]entity.Company, error)
	}
)
