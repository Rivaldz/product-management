package v1

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "technical_test/docs"
	"technical_test/internal/usecase"
	"technical_test/pkg/logger"
)

func NewRouter(handler *gin.Engine, itemUC usecase.Item, companyUC usecase.Company, l logger.Interface) {
	// Options
	handler.Use(gin.Logger())
	handler.Use(gin.Recovery())

	// Swagger
	swaggerHandler := ginSwagger.DisablingWrapHandler(swaggerFiles.Handler, "DISABLE_SWAGGER_HTTP_HANDLER")
	handler.GET("/swagger/*any", swaggerHandler)

	// API routes
	h := handler.Group("/")
	{
		newItemRoutes(h, itemUC, l)
		newCompanyRoutes(h, companyUC, l)
	}
}
