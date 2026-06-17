package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"technical_test/internal/entity"
	"technical_test/internal/usecase"
	"technical_test/pkg/logger"
)

type companyRoutes struct {
	c usecase.Company
	l logger.Interface
}

func newCompanyRoutes(handler *gin.RouterGroup, c usecase.Company, l logger.Interface) {
	r := &companyRoutes{c, l}

	h := handler.Group("/companies")
	{
		h.GET("", r.list)
	}
}

// @Summary      List companies
// @Description  Get a list of all active companies containing uuid and name
// @Tags         companies
// @Accept       json
// @Produce      json
// @Success      200         {object}  Response{data=[]entity.Company}
// @Failure      500         {object}  Response
// @Router       /companies [get]
func (r *companyRoutes) list(c *gin.Context) {
	var companies []entity.Company
	var err error

	companies, err = r.c.GetList(c.Request.Context())
	if err != nil {
		r.l.Error(err, "HTTP - v1 - company - list")
		errorResponse(c, http.StatusInternalServerError, "INTERNAL_ERROR", "Internal server error", nil)
		return
	}

	successResponse(c, http.StatusOK, "COMPANY_LIST_RETRIEVED", "Company list retrieved successfully", companies)
}
