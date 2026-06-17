package v1_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"technical_test/internal/entity"
)

func TestHTTP_ListCompanies(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		companyUC := &mockCompanyUseCase{
			getListFn: func(ctx context.Context) ([]entity.Company, error) {
				return []entity.Company{
					{ID: "c0a80121-7ac0-11d1-898c-00c04fd8d5c1", Name: "Test Company A"},
				}, nil
			},
		}

		router := setupRouter(&mockItemUseCase{}, companyUC)

		req, _ := http.NewRequest("GET", "/api/v1/companies", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusOK, resp.Code, resp.Body.String())
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		companyUC := &mockCompanyUseCase{
			getListFn: func(ctx context.Context) ([]entity.Company, error) {
				return nil, errors.New("db error")
			},
		}

		router := setupRouter(&mockItemUseCase{}, companyUC)

		req, _ := http.NewRequest("GET", "/api/v1/companies", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusInternalServerError, resp.Code, resp.Body.String())
		}
	})
}
