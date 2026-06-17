package usecase_test

import (
	"context"
	"testing"

	"technical_test/internal/entity"
	"technical_test/internal/usecase"
)

type mockCompanyRepo struct {
	getListFn func(ctx context.Context) ([]entity.Company, error)
}

func (m *mockCompanyRepo) GetList(ctx context.Context) ([]entity.Company, error) {
	return m.getListFn(ctx)
}

func TestCompanyUseCase_GetList(t *testing.T) {
	expectedCompanies := []entity.Company{
		{ID: "comp-1", Name: "Company 1"},
		{ID: "comp-2", Name: "Company 2"},
	}

	repo := &mockCompanyRepo{
		getListFn: func(ctx context.Context) ([]entity.Company, error) {
			return expectedCompanies, nil
		},
	}
	uc := usecase.NewCompanyUseCase(repo)

	res, err := uc.GetList(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(res) != 2 {
		t.Fatalf("expected 2 companies, got %d", len(res))
	}

	if res[0].ID != "comp-1" || res[0].Name != "Company 1" {
		t.Errorf("unexpected company data at index 0: %v", res[0])
	}
}
