package v1_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"

	v1 "technical_test/internal/controller/http/v1"
	"technical_test/internal/entity"
	"technical_test/pkg/logger"
)

type mockItemUseCase struct {
	createFn    func(ctx context.Context, item entity.Item) (entity.Item, error)
	getListFn   func(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error)
	getDetailFn func(ctx context.Context, companyID, itemID string) (entity.Item, error)
	updateFn    func(ctx context.Context, item entity.Item) (entity.Item, error)
	archiveFn   func(ctx context.Context, companyID, itemID string) error
}

func (m *mockItemUseCase) Create(ctx context.Context, item entity.Item) (entity.Item, error) {
	return m.createFn(ctx, item)
}

func (m *mockItemUseCase) GetList(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
	return m.getListFn(ctx, companyID, filter)
}

func (m *mockItemUseCase) GetDetail(ctx context.Context, companyID, itemID string) (entity.Item, error) {
	return m.getDetailFn(ctx, companyID, itemID)
}

func (m *mockItemUseCase) Update(ctx context.Context, item entity.Item) (entity.Item, error) {
	return m.updateFn(ctx, item)
}

func (m *mockItemUseCase) Archive(ctx context.Context, companyID, itemID string) error {
	return m.archiveFn(ctx, companyID, itemID)
}

type mockCompanyUseCase struct {
	getListFn func(ctx context.Context) ([]entity.Company, error)
}

func (m *mockCompanyUseCase) GetList(ctx context.Context) ([]entity.Company, error) {
	return m.getListFn(ctx)
}

func setupRouter(itemUC *mockItemUseCase, companyUC *mockCompanyUseCase) *gin.Engine {
	gin.SetMode(gin.TestMode)
	handler := gin.New()
	l := logger.New("error")
	v1.NewRouter(handler, itemUC, companyUC, l)
	return handler
}

func TestHTTP_CreateItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			createFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
				item.ID = "generated-uuid"
				item.Status = entity.StatusActive
				return item, nil
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		reqBody, _ := json.Marshal(v1.CreateItemRequest{
			Code:         "CODE1",
			Name:         "Product 1",
			Type:         "PRODUCT",
			Price:        125000,
			CategoryName: "Electronics",
		})

		req, _ := http.NewRequest("POST", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d. Body: %s", http.StatusCreated, resp.Code, resp.Body.String())
		}
	})

	t.Run("validation error", func(t *testing.T) {
		router := setupRouter(&mockItemUseCase{}, &mockCompanyUseCase{})

		reqBody, _ := json.Marshal(v1.CreateItemRequest{
			Code:  "", // invalid
			Name:  "Product 1",
			Type:  "INVALID_TYPE", // invalid
			Price: -100,           // invalid
		})

		req, _ := http.NewRequest("POST", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusBadRequest {
			t.Errorf("expected status %d, got %d", http.StatusBadRequest, resp.Code)
		}
	})

	t.Run("conflict error", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			createFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
				return entity.Item{}, entity.ErrDuplicateCode
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		reqBody, _ := json.Marshal(v1.CreateItemRequest{
			Code:  "CODE1",
			Name:  "Product 1",
			Type:  "PRODUCT",
			Price: 100,
		})

		req, _ := http.NewRequest("POST", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, resp.Code)
		}
	})
}

func TestHTTP_ListItems(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			getListFn: func(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
				return []entity.Item{
					{ID: "id-1", Code: "C1", Name: "N1", Type: "PRODUCT"},
				}, nil
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		req, _ := http.NewRequest("GET", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items?type=PRODUCT", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			getListFn: func(ctx context.Context, companyID string, filter entity.ItemFilter) ([]entity.Item, error) {
				return nil, errors.New("db connection failed")
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		req, _ := http.NewRequest("GET", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, resp.Code)
		}
	})
}

func TestHTTP_GetItemDetail(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			getDetailFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{ID: itemID, CompanyID: companyID, Code: "C1"}, nil
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		req, _ := http.NewRequest("GET", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/item-id-1", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("not found", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			getDetailFn: func(ctx context.Context, companyID, itemID string) (entity.Item, error) {
				return entity.Item{}, entity.ErrItemNotFound
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		req, _ := http.NewRequest("GET", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/invalid-id", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, resp.Code)
		}
	})
}

func TestHTTP_UpdateItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			updateFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
				return item, nil
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		reqBody, _ := json.Marshal(v1.UpdateItemRequest{
			Code:   "C2",
			Name:   "N2",
			Type:   "SERVICE",
			Price:  1000,
			Status: "ACTIVE",
		})

		req, _ := http.NewRequest("PATCH", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/item-1", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})

	t.Run("conflict - item archived", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			updateFn: func(ctx context.Context, item entity.Item) (entity.Item, error) {
				return entity.Item{}, entity.ErrAlreadyArchived
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		reqBody, _ := json.Marshal(v1.UpdateItemRequest{
			Code:   "C2",
			Name:   "N2",
			Type:   "SERVICE",
			Price:  1000,
			Status: "ACTIVE",
		})

		req, _ := http.NewRequest("PATCH", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/item-1", bytes.NewBuffer(reqBody))
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusConflict {
			t.Errorf("expected status %d, got %d", http.StatusConflict, resp.Code)
		}
	})
}

func TestHTTP_ArchiveItem(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		itemUC := &mockItemUseCase{
			archiveFn: func(ctx context.Context, companyID, itemID string) error {
				return nil
			},
		}
		router := setupRouter(itemUC, &mockCompanyUseCase{})

		req, _ := http.NewRequest("PATCH", "/api/v1/companies/c0a80121-7ac0-11d1-898c-00c04fd8d5c1/items/item-1/archive", nil)
		resp := httptest.NewRecorder()

		router.ServeHTTP(resp, req)

		if resp.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, resp.Code)
		}
	})
}
