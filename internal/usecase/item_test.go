package usecase_test

import (
	"context"
	"errors"
	"testing"

	"technical_test/internal/entity"
	"technical_test/internal/usecase"
)

type mockItemRepo struct {
	insertFn      func(ctx context.Context, item entity.Item) (entity.Item, error)
	getByFilterFn func(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error)
	getByIDFn     func(ctx context.Context, companyID, itemID string) (entity.Item, error)
	updateFn      func(ctx context.Context, item entity.Item) (entity.Item, error)
	archiveFn     func(ctx context.Context, companyID, itemID string) error
}

func (m *mockItemRepo) Insert(ctx context.Context, item entity.Item) (entity.Item, error) {
	return m.insertFn(ctx, item)
}

func (m *mockItemRepo) GetByFilter(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
	return m.getByFilterFn(ctx, companyID, filter)
}

func (m *mockItemRepo) GetByID(ctx context.Context, companyID, itemID string) (entity.Item, error) {
	return m.getByIDFn(ctx, companyID, itemID)
}

func (m *mockItemRepo) Update(ctx context.Context, item entity.Item) (entity.Item, error) {
	return m.updateFn(ctx, item)
}

func (m *mockItemRepo) Archive(ctx context.Context, companyID, itemID string) error {
	return m.archiveFn(ctx, companyID, itemID)
}

func TestItemUseCase_Create(t *testing.T) {
	repo := &mockItemRepo{
		insertFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
			item.ID = "generated-id"
			return item, nil
		},
	}
	uc := usecase.NewItemUseCase(repo)

	item := entity.Item{
		CompanyID: "comp-1",
		Code:      "CODE1",
		Name:      "Product 1",
		Type:      entity.TypeProduct,
		Price:     100.0,
	}

	res, err := uc.Create(context.Background(), item)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != "generated-id" {
		t.Errorf("expected ID 'generated-id', got %s", res.ID)
	}
	if res.Status != entity.StatusActive {
		t.Errorf("expected Status '%s', got %s", entity.StatusActive, res.Status)
	}
}

func TestItemUseCase_GetList(t *testing.T) {
	expectedItems := []entity.Item{
		{ID: "item-1", CompanyID: "comp-1", Code: "CODE1", Name: "Product 1"},
	}
	repo := &mockItemRepo{
		getByFilterFn: func(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
			return expectedItems, nil
		},
	}
	uc := usecase.NewItemUseCase(repo)

	res, err := uc.GetList(context.Background(), "comp-1", entity.ItemFilter{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 1 {
		t.Errorf("expected 1 item, got %d", len(res))
	}
}

func TestItemUseCase_GetDetail(t *testing.T) {
	expectedItem := entity.Item{ID: "item-1", CompanyID: "comp-1", Code: "CODE1"}
	repo := &mockItemRepo{
		getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
			return expectedItem, nil
		},
	}
	uc := usecase.NewItemUseCase(repo)

	res, err := uc.GetDetail(context.Background(), "comp-1", "item-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if res.ID != "item-1" {
		t.Errorf("expected ID 'item-1', got %s", res.ID)
	}
}

func TestItemUseCase_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{ID: itemID, CompanyID: companyID, Status: entity.StatusActive}, nil
			},
			updateFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
				return item, nil
			},
		}
		uc := usecase.NewItemUseCase(repo)

		item := entity.Item{ID: "item-1", CompanyID: "comp-1", Code: "CODE1", Name: "Updated Name", Status: entity.StatusActive}
		res, err := uc.Update(context.Background(), item)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if res.Name != "Updated Name" {
			t.Errorf("expected Name 'Updated Name', got %s", res.Name)
		}
	})

	t.Run("error get by id failed", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{}, dbErr
			},
		}
		uc := usecase.NewItemUseCase(repo)

		item := entity.Item{ID: "item-1", CompanyID: "comp-1"}
		_, err := uc.Update(context.Background(), item)
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected db error, got %v", err)
		}
	})

	t.Run("error already archived", func(t *testing.T) {
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{ID: itemID, CompanyID: companyID, Status: entity.StatusArchived}, nil
			},
		}
		uc := usecase.NewItemUseCase(repo)

		item := entity.Item{ID: "item-1", CompanyID: "comp-1"}
		_, err := uc.Update(context.Background(), item)
		if !errors.Is(err, entity.ErrAlreadyArchived) {
			t.Fatalf("expected ErrAlreadyArchived, got %v", err)
		}
	})
}

func TestItemUseCase_Archive(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{ID: itemID, CompanyID: companyID, Status: entity.StatusActive}, nil
			},
			archiveFn: func(ctx context.Context, companyID, itemID string) error {
				return nil
			},
		}
		uc := usecase.NewItemUseCase(repo)

		err := uc.Archive(context.Background(), "comp-1", "item-1")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
	})

	t.Run("error get by id failed", func(t *testing.T) {
		dbErr := errors.New("db error")
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{}, dbErr
			},
		}
		uc := usecase.NewItemUseCase(repo)

		err := uc.Archive(context.Background(), "comp-1", "item-1")
		if !errors.Is(err, dbErr) {
			t.Fatalf("expected db error, got %v", err)
		}
	})

	t.Run("error already archived", func(t *testing.T) {
		repo := &mockItemRepo{
			getByIDFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{ID: itemID, CompanyID: companyID, Status: entity.StatusArchived}, nil
			},
		}
		uc := usecase.NewItemUseCase(repo)

		err := uc.Archive(context.Background(), "comp-1", "item-1")
		if !errors.Is(err, entity.ErrAlreadyArchived) {
			t.Fatalf("expected ErrAlreadyArchived, got %v", err)
		}
	})
}
