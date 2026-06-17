package v1

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"

	"technical_test/internal/entity"
	"technical_test/internal/usecase"
	"technical_test/pkg/logger"
)

type itemRoutes struct {
	i usecase.Item
	l logger.Interface
}

func newItemRoutes(handler *gin.RouterGroup, i usecase.Item, l logger.Interface) {
	r := &itemRoutes{i, l}

	h := handler.Group("/companies/:company_id/items")
	h.Use(UUIDParamValidator("company_id"))
	{
		h.POST("", r.create)
		h.GET("", r.list)
		h.GET("/:item_id", UUIDParamValidator("item_id"), r.detail)
		h.PATCH("/:item_id", UUIDParamValidator("item_id"), r.update)
		h.PATCH("/:item_id/archive", UUIDParamValidator("item_id"), r.archive)
	}
}

type CreateItemRequest struct {
	Code         string   `json:"code" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Type         string   `json:"type" binding:"required,oneof=PRODUCT SERVICE"`
	Price        *float64 `json:"price" binding:"required,min=0"`
	CategoryName string   `json:"category_name"`
}

type UpdateItemRequest struct {
	Code         string   `json:"code" binding:"required"`
	Name         string   `json:"name" binding:"required"`
	Type         string   `json:"type" binding:"required,oneof=PRODUCT SERVICE"`
	Price        *float64 `json:"price" binding:"required,min=0"`
	CategoryName string   `json:"category_name"`
	Status       string   `json:"status" binding:"required,oneof=ACTIVE INACTIVE"`
}

type Response struct {
	Success bool        `json:"success"`
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

func errorResponse(c *gin.Context, httpCode int, code, msg string, errs interface{}) {
	c.JSON(httpCode, Response{
		Success: false,
		Code:    code,
		Message: msg,
		Errors:  errs,
	})
}

func successResponse(c *gin.Context, httpCode int, code, msg string, data interface{}) {
	c.JSON(httpCode, Response{
		Success: true,
		Code:    code,
		Message: msg,
		Data:    data,
	})
}

// @Summary      Create a new item
// @Description  Create a new product or service for a company
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        company_id  path      string              true  "Company ID"
// @Param        request     body      CreateItemRequest   true  "Item creation payload"
// @Success      201         {object}  Response{data=entity.Item}
// @Failure      400         {object}  Response
// @Failure      409         {object}  Response
// @Failure      500         {object}  Response
// @Router       /companies/{company_id}/items [post]
func (r *itemRoutes) create(c *gin.Context) {
	companyID := c.Param("company_id")

	var req CreateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	item := entity.Item{
		CompanyID:    companyID,
		Code:         req.Code,
		Name:         req.Name,
		Type:         req.Type,
		Price:        *req.Price,
		CategoryName: req.CategoryName,
	}

	created, err := r.i.Create(c.Request.Context(), item)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateCode) {
			errorResponse(c, http.StatusConflict, "ITEM_CODE_ALREADY_EXISTS", err.Error(), nil)
			return
		}
		r.l.Error(err, "HTTP - v1 - create")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusCreated, "ITEM_CREATED", "Item created successfully", created)
}

// @Summary      List items
// @Description  Get a list of items belonging to a company with filters
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        company_id  path      string  true   "Company ID"
// @Param        type        query     string  false  "Item Type (PRODUCT, SERVICE)"
// @Param        status      query     string  false  "Item Status (ACTIVE, INACTIVE, ARCHIVED)"
// @Param        keyword     query     string  false  "Search by code or name"
// @Success      200         {object}  Response{data=[]entity.Item}
// @Failure      500         {object}  Response
// @Router       /companies/{company_id}/items [get]
func (r *itemRoutes) list(c *gin.Context) {
	companyID := c.Param("company_id")
	filter := entity.ItemFilter{
		Type:    c.Query("type"),
		Status:  c.Query("status"),
		Keyword: c.Query("keyword"),
	}

	items, err := r.i.GetList(c.Request.Context(), companyID, filter)
	if err != nil {
		r.l.Error(err, "HTTP - v1 - list")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusOK, "ITEM_LIST_RETRIEVED", "Item list retrieved successfully", items)
}

// @Summary      Get item detail
// @Description  Get details of a specific item for a company
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        company_id  path      string  true  "Company ID"
// @Param        item_id     path      string  true  "Item ID"
// @Success      200         {object}  Response{data=entity.Item}
// @Failure      404         {object}  Response
// @Failure      500         {object}  Response
// @Router       /companies/{company_id}/items/{item_id} [get]
func (r *itemRoutes) detail(c *gin.Context) {
	companyID := c.Param("company_id")
	itemID := c.Param("item_id")

	item, err := r.i.GetDetail(c.Request.Context(), companyID, itemID)
	if err != nil {
		if errors.Is(err, entity.ErrItemNotFound) {
			errorResponse(c, http.StatusNotFound, "ITEM_NOT_FOUND", err.Error(), nil)
			return
		}
		r.l.Error(err, "HTTP - v1 - detail")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusOK, "ITEM_DETAIL_RETRIEVED", "Item detail retrieved successfully", item)
}

// @Summary      Update an item
// @Description  Update a specific item's details. Archived items cannot be updated.
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        company_id  path      string              true  "Company ID"
// @Param        item_id     path      string              true  "Item ID"
// @Param        request     body      UpdateItemRequest   true  "Item update payload"
// @Success      200         {object}  Response{data=entity.Item}
// @Failure      400         {object}  Response
// @Failure      404         {object}  Response
// @Failure      409         {object}  Response
// @Failure      500         {object}  Response
// @Router       /companies/{company_id}/items/{item_id} [patch]
func (r *itemRoutes) update(c *gin.Context) {
	companyID := c.Param("company_id")
	itemID := c.Param("item_id")

	var req UpdateItemRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		errorResponse(c, http.StatusBadRequest, "VALIDATION_ERROR", "Validation failed", err.Error())
		return
	}

	item := entity.Item{
		ID:           itemID,
		CompanyID:    companyID,
		Code:         req.Code,
		Name:         req.Name,
		Type:         req.Type,
		Price:        *req.Price,
		CategoryName: req.CategoryName,
		Status:       req.Status,
	}

	updated, err := r.i.Update(c.Request.Context(), item)
	if err != nil {
		if errors.Is(err, entity.ErrDuplicateCode) {
			errorResponse(c, http.StatusConflict, "ITEM_CODE_ALREADY_EXISTS", err.Error(), nil)
			return
		}
		if errors.Is(err, entity.ErrAlreadyArchived) {
			errorResponse(c, http.StatusConflict, "ITEM_ALREADY_ARCHIVED", err.Error(), nil)
			return
		}
		if errors.Is(err, entity.ErrItemNotFound) {
			errorResponse(c, http.StatusNotFound, "ITEM_NOT_FOUND", err.Error(), nil)
			return
		}
		r.l.Error(err, "HTTP - v1 - update")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusOK, "ITEM_UPDATED", "Item updated successfully", updated)
}

// @Summary      Archive an item
// @Description  Archive a specific item. Archived items cannot be archived again.
// @Tags         items
// @Accept       json
// @Produce      json
// @Param        company_id  path      string  true  "Company ID"
// @Param        item_id     path      string  true  "Item ID"
// @Success      200         {object}  Response
// @Failure      404         {object}  Response
// @Failure      409         {object}  Response
// @Failure      500         {object}  Response
// @Router       /companies/{company_id}/items/{item_id}/archive [patch]
func (r *itemRoutes) archive(c *gin.Context) {
	companyID := c.Param("company_id")
	itemID := c.Param("item_id")

	err := r.i.Archive(c.Request.Context(), companyID, itemID)
	if err != nil {
		if errors.Is(err, entity.ErrAlreadyArchived) {
			errorResponse(c, http.StatusConflict, "ITEM_ALREADY_ARCHIVED", err.Error(), nil)
			return
		}
		if errors.Is(err, entity.ErrItemNotFound) {
			errorResponse(c, http.StatusNotFound, "ITEM_NOT_FOUND", err.Error(), nil)
			return
		}
		r.l.Error(err, "HTTP - v1 - archive")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusOK, "ITEM_ARCHIVED", "Item archived successfully", gin.H{
		"item_id":    itemID,
		"company_id": companyID,
		"status":     entity.StatusArchived,
	})
}
